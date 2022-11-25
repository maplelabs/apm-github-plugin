/* Package Task contains Task related methods and information
 */
package task

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/maplelabs/github-audit/gitprovider"
	"github.com/maplelabs/github-audit/input"
	"github.com/maplelabs/github-audit/internal/dataprocessor"
	"github.com/maplelabs/github-audit/logger"
	"github.com/maplelabs/github-audit/publisher"
)

var (
	log logger.Logger
	// TaskStatsMap is global variable holding stats related to task with task ID as key.
	TaskStatsMap = make(map[string]TaskStats)
	// taskStatsFile is the file where tasks stats are saved regularly.
	taskStatsFile = "taskStats.json"
	// taskStatsFile keeps concurrency good between multiple goroutines handling task stats.
	taskStatsMutex = &sync.Mutex{}
	// saveTaskPeriodicInterval is interval for which checkpoint for task is done
	saveTaskPeriodicInterval = 30 * time.Second
)

// Task represents a single task where one auditjob = one task
type Task struct {
	// ID is unique for a task Format: auditJobName$repoOwner$repoName.
	ID string

	// Scheduling interval between two runs of the task.
	SchedulingInterval time.Duration

	// isRunning flag to know if the task is currently running. This
	// flag will ensure that only one instance of the task will run at any time.
	IsRunning int64

	// previous run time of the task.
	PreviousRunTime time.Time

	// next run time of the task.
	NextRunTime time.Time

	// task params needed for task execution.
	TaskParams
}

// TaskParams represents task params needed to run a task.
type TaskParams struct {
	// config contains auditjob config needed to run task.
	Config input.AuditJob

	// Targets contains all target where data needs to be published.
	Targets []input.Target

	// DecodeAccessKey represents decoded access key.
	DecodeAccessKey string
}

// TaskStats represents task running stats to handle githu-audit restarts and checkpointing.
type TaskStats struct {
	// TaskID is unique for a task Format: auditJobName$repoOwner$repoName.
	TaskID string

	// LastSuccessFullRunTime of the task.
	LastSuccessFullRunTime time.Time

	// LastPullRequestNo represents last fetched pull request number.
	LastPullRequestNo int

	// LastCommitTime represents the last commit time for the audit job.
	LastCommitTime map[string]time.Time

	// LastIssueTime represents the last issue time for the audit job.
	LastIssueTime time.Time
}

func init() {
	log = logger.GetLogger()

	// checking if task stats file is present or not
	_, err := os.Stat(taskStatsFile)
	// only if file exists , unmarshalling to TaskStats global variable
	if err == nil {
		fileByte, err := os.ReadFile(taskStatsFile)
		if err != nil {
			log.Error("error[%v] in reading task stats file", err)
		}
		err = json.Unmarshal(fileByte, &TaskStatsMap)
		if err != nil {
			log.Error("error[%v] in unmarshalling task stats file", err)
		}
	} else {
		log.Error("error[%v] in reading task stats file", err)
	}
}

// SaveTaskStatsPeriodic runs periodically to save task stats to file system (checkpointing mechanism).
func SaveTaskStatsPeriodic(ctx context.Context) {
	ticker := time.NewTicker(saveTaskPeriodicInterval)
	for {
		select {
		case <-ticker.C:
			taskStatsMutex.Lock()
			tBytes, err := json.Marshal(TaskStatsMap)
			if err != nil {
				log.Errorf("error[%v] in marshalling TaskStatsMap for checkpointing", err)
			}
			err = os.WriteFile(taskStatsFile, tBytes, 0644)
			if err != nil {
				log.Errorf("error[%v] in writing taskStatsFile file", err)
			}
			taskStatsMutex.Unlock()
		case <-ctx.Done():
			log.Info("stopping periodic save task stats routine")
			return
		}
	}
}

// getTaskStats returns task stats for particular task.
func getTaskStats(id string) (TaskStats, error) {
	taskStatsMutex.Lock()
	defer taskStatsMutex.Unlock()
	if v, ok := TaskStatsMap[id]; ok {
		return v, nil
	}
	return TaskStats{}, errors.New("task not found")
}

// saveTaskStats saves task stats for a particular task.
func saveTaskStats(id string, ts TaskStats) {
	taskStatsMutex.Lock()
	defer taskStatsMutex.Unlock()
	TaskStatsMap[id] = ts
}

// Newtask returns new task instance.
func Newtask() *Task {
	task := new(Task)
	// keeping previous time to -N minutes before scheduling interval.
	task.PreviousRunTime = time.Now().Add(-task.SchedulingInterval)
	// keeping nextruntime as current time.
	task.NextRunTime = time.Now()
	return task
}

// AddTaskParams methods adds tasks params.
func (t *Task) AddTaskParams(tp TaskParams) {
	t.TaskParams = tp
	t.ID = tp.Config.Name + "$" + tp.Config.RepositoryOwner + "$" + tp.Config.RepositoryName
}

// AddInterval methods adds tasks scheduling interval.
func (t *Task) AddInterval(interval time.Duration) {
	t.SchedulingInterval = interval
}

// scheduleNextRunOfTask schedules the next run of the task.
func (task *Task) ScheduleNextRunOfTask() {
	task.NextRunTime = time.Now().Add(task.SchedulingInterval)
}

// Start methods starts teh execution of a particular task.
func (t *Task) Start() error {
	var ts TaskStats
	ts, err := getTaskStats(t.ID)
	// if error reading previous stats , putting default values for task stats map
	if err != nil {
		ts.LastIssueTime = t.PreviousRunTime
		for _, br := range t.Config.Branches {
			ts.LastCommitTime[br] = t.PreviousRunTime
		}
		ts.LastPullRequestNo = 0
		ts.TaskID = t.ID
		// saving default task stats for a task.
		saveTaskStats(t.ID, ts)
	}
	gp := gitprovider.NewGitProvider(t.Config.RepositoryHost, t.Config.RepositoryOwner, t.Config.RepositoryName, t.Config.Username, t.DecodeAccessKey)

	// max concurrency guard to control goroutines
	maxConcurrencyGuard := make(chan struct{}, runtime.NumCPU()*2)
	wg := &sync.WaitGroup{}
	// executing each target in task in separate goroutines
	for _, tar := range t.Targets {
		maxConcurrencyGuard <- struct{}{}
		wg.Add(1)
		go func(tar input.Target, maxConcurrencyGuard chan struct{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() { <-maxConcurrencyGuard }()
			// getting new publisher
			pb, err := publisher.NewPublisher(tar.Type, tar.TargetConfig)
			if err != nil {
				log.Errorf("error[%v] in getting publisher for the task with ID %v", err, t.ID)
				return
			}
			// getting new dataprocessor
			dp := dataprocessor.NewDataProcessor(t.Config.RepositoryHost, t.Config.RepositoryName, t.Config.RepositoryURL)
			err = t.collectAndPublishCommits(gp, pb, dp, ts)
			if err != nil {
				log.Errorf("error[%v] in collecting commits for task with ID %v", err, t.ID)
				return
			}
			err = t.collectAndPublishPullRequests(gp, pb, dp, ts)
			if err != nil {
				log.Errorf("error[%v] in collecting pull requests for task with ID %v", err, t.ID)
				return
			}
			err = t.collectAndPublishIssues(gp, pb, dp, ts)
			if err != nil {
				log.Errorf("error[%v] in collecting issues for task with ID %v", err, t.ID)
				return
			}
		}(tar, maxConcurrencyGuard, wg)
	}
	// waiting for all concurrent goroutines to complete
	wg.Wait()
	ts, err = getTaskStats(t.ID)
	if err != nil {
		log.Errorf("error[%v] in getting task stats for task with ID %v", err, t.ID)
		return err
	}
	ts.LastSuccessFullRunTime = time.Now()
	// saving task stats.
	saveTaskStats(t.ID, ts)
	return nil
}

//TODO: add stop function in future if required
func (t *Task) Stop() error {
	return nil
}

// collectAndPublishCommits collects commits and publish them to targets.
func (t *Task) collectAndPublishCommits(gp gitprovider.GitProvider, pb publisher.Publisher, dp dataprocessor.DataProcessor, ts TaskStats) error {
	var err error
	errChan := make(chan error)
	// max concurrency guard to control goroutines
	maxConcurrencyGuard := make(chan struct{}, runtime.NumCPU()*2)
	for _, br := range t.Config.Branches {
		maxConcurrencyGuard <- struct{}{}
		// executing each branch in separate goroutines
		go func(br string, maxConcurrencyGuard chan struct{}, errChan chan error) {
			defer func() { <-maxConcurrencyGuard }()
			commitBytes, err := gp.GetCommits(ts.LastCommitTime[br], time.Now(), br)
			if err != nil {
				log.Errorf("error[%v] in getting commits from gitprovider for task with ID %v", err, t.ID)
				errChan <- err
				return
			}
			processed, err := dp.ProcessCommits(commitBytes, t.Config.Tags)
			if err != nil {
				log.Errorf("error[%v] in processing commits for task with ID %v", err, t.ID)
				errChan <- err
				return
			}
			err = pb.Publish(processed)
			if err != nil {
				log.Errorf("error[%v] in publishing commits for task with ID %v", err, t.ID)
				errChan <- err
				return
			}
			// saving stats after finished task
			for _, v := range processed {
				//taking latest commit time
				lastCommitTime := v.(dataprocessor.Commit).CreatedAt
				ts.LastCommitTime[br] = lastCommitTime
				saveTaskStats(t.ID, ts)
				break
			}
			errChan <- nil
		}(br, maxConcurrencyGuard, errChan)
	}
	//waiting for all goroutines to complete based on errChannel
	for errVal := range errChan {
		if errVal != nil {
			err = errVal
		}
	}
	return err
}

// collectAndPublishPullRequests collects pull requests and publish them to targets.
func (t *Task) collectAndPublishPullRequests(gp gitprovider.GitProvider, pb publisher.Publisher, dp dataprocessor.DataProcessor, ts TaskStats) error {
	prBytes, err := gp.GetPullRequests(ts.LastPullRequestNo)
	if err != nil {
		log.Errorf("error[%v] in getting commits from gitprovider for task with ID %v", err, t.ID)
		return err
	}
	processed, err := dp.ProcessCommits(prBytes, t.Config.Tags)
	if err != nil {
		log.Errorf("error[%v] in processing pull requests for task with ID %v", err, t.ID)
		return err
	}
	err = pb.Publish(processed)
	if err != nil {
		log.Errorf("error[%v] in publishing commits for task with ID %v", err, t.ID)
		return err
	}
	// saving stats after finished task
	for _, v := range processed {
		//taking latest pull request number
		lastPrNo := v.(dataprocessor.PullRequest).PullRequestNo
		ts.LastPullRequestNo, _ = strconv.Atoi(lastPrNo)
		saveTaskStats(t.ID, ts)
		break
	}
	return nil
}

// collectAndPublishPullRequests collects issues and publish them to targets.
func (t *Task) collectAndPublishIssues(gp gitprovider.GitProvider, pb publisher.Publisher, dp dataprocessor.DataProcessor, ts TaskStats) error {
	issuesBytes, err := gp.GetIssues(ts.LastIssueTime)
	if err != nil {
		log.Errorf("error[%v] in getting commits from gitprovider for task with ID %v", err, t.ID)
		return err
	}
	processed, err := dp.ProcessCommits(issuesBytes, t.Config.Tags)
	if err != nil {
		log.Errorf("error[%v] in processing issues for task with ID %v", err, t.ID)
		return err
	}
	err = pb.Publish(processed)
	if err != nil {
		log.Errorf("error[%v] in publishing commits for task with ID %v", err, t.ID)
		return err
	}
	// saving stats after finished task
	for _, v := range processed {
		//taking latest issue time
		lastIssueTime := v.(dataprocessor.Issue).CreatedAt
		ts.LastIssueTime = lastIssueTime
		saveTaskStats(t.ID, ts)
		break
	}
	return nil
}
