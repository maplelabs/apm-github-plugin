/* Package Configurator is responsible for creating tasks from auditjobs mentioned in the config.yaml
 */

package configurator

import (
	"errors"
	"strings"
	"time"

	"github.com/maplelabs/github-audit/gitprovider"
	"github.com/maplelabs/github-audit/input"
	"github.com/maplelabs/github-audit/internal/task"
	"github.com/maplelabs/github-audit/logger"
	"github.com/maplelabs/github-audit/utils"
)

const (
	// Default Scheduling interval for each audit job.
	DEFAULTDURATION = time.Duration(5 * time.Minute)
)

var (
	log logger.Logger
	// ErrNoTaskConfigured if there is no audit job configured.
	ErrNoTaskConfigured = errors.New("no audit task configured")
)

func init() {
	log = logger.GetLogger()
}

// StartProcessing returns list of tasks with valid credentials which can be scheduled.
func StartProcessing(config input.Config) ([]*task.Task, error) {
	log.Info("starting creation of tasks for configured audit jobs")
	tasks := make([]*task.Task, 0)
	for _, aj := range config.AuditJobs {
		task, err := createTask(aj, config.Targets)
		if err != nil {
			log.Errorf("error[%v] in creating task for audit job %v", err, aj.Name)
			continue
		}
		log.Debugf("configured tasks for scheduling audit job %v with taskID %v", aj.Name, task.ID)
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		log.Errorf("error[%v] in creating audit tasks", ErrNoTaskConfigured)
		return nil, ErrNoTaskConfigured
	}
	return tasks, nil
}

// createTask creates a single task based on audit job config and targets.
func createTask(auditJob input.AuditJob, targets []input.Target) (*task.Task, error) {
	var decodedKey string
	var err error
	task := task.Newtask()
	if auditJob.AccessToken != "" {
		decodedKey, err = utils.DecodeAccessKey(auditJob.AccessToken)
		if err != nil {
			log.Errorf("error[%v] in decoding accessToken for audit job %v", err, auditJob.Name)
			return task, err
		}
	}
	gp := gitprovider.NewGitProvider(auditJob.RepositoryHost, auditJob.RepositoryOwner, auditJob.RepositoryName, auditJob.Username, decodedKey)
	err = gp.CheckCredentials()
	if err != nil {
		log.Errorf("error[%v] in authenticating gitprovider for audit job %v", err, auditJob.Name)
		return task, err
	}
	taskParams := createTaskParam(auditJob, targets, decodedKey)
	task.AddTaskParams(taskParams)
	interval := convertIntervalToDuration(auditJob.PollingInterval)
	task.AddInterval(interval)
	return task, nil
}

// createTaskParam creates params for the task that are needed for running.
func createTaskParam(auditJob input.AuditJob, targets []input.Target, decodedKey string) task.TaskParams {
	var taskParam task.TaskParams
	taskParam.Config = auditJob
	jobTargets := make([]input.Target, 0)
	for _, tar := range targets {
		for i := range auditJob.Output.TargetName {
			if auditJob.Output.TargetName[i] == tar.Name {
				jobTargets = append(jobTargets, tar)
			}
		}
	}
	taskParam.Targets = jobTargets
	taskParam.DecodeAccessKey = decodedKey
	return taskParam
}

// convertIntervalToDuration converts the scheduling interval as provided in config.yaml to golang's duration
func convertIntervalToDuration(interval string) time.Duration {
	interval = strings.ToLower(interval)
	duration, err := time.ParseDuration(interval)
	if err != nil {
		log.Errorf("error[%v] in converting scheduling duration to golang's duration , setting default duration of 5 minutes", err)
		return DEFAULTDURATION
	}
	return duration
}
