/* Package Taskmanager is schduling tasks at regular interval similar to cron jobs
 */
package taskmanager

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/maplelabs/github-audit/internal/task"
	"github.com/maplelabs/github-audit/logger"
)

var (
	log logger.Logger
)

func init() {
	log = logger.GetLogger()
}

// TaskManager holds information regarding the scheduler which runs tasks.
type TaskManager struct {
	// The tasks queue that the scheduler holds.
	taskQueue []*task.Task

	// The mutex used for add tasks.
	taskQueueMutex sync.Mutex

	// The max concurrency guard will make sure that only fixed number of
	// allowed tasks to run concurrently.
	maxConcurrency chan struct{}
}

// StartTasks starts the tasks passed as input.
func StartTasks(ctx context.Context, tasks []*task.Task) {
	maxConcurrency := runtime.NumCPU() * 2
	tm := NewTaskManager(int64(maxConcurrency))
	for _, task := range tasks {
		tm.AddTask(task)
	}
	// starting save tasks stats periodic routine
	go task.SaveTaskStatsPeriodic(ctx)
	//blocking call , will exit only if github-audit stops or crashes
	tm.startScheduling(ctx)
}

// NewTaskManager creates new instance of the task manager.
func NewTaskManager(maxConcurrency int64) *TaskManager {
	tm := new(TaskManager)
	tm.maxConcurrency = make(chan struct{}, maxConcurrency)
	tm.taskQueue = make([]*task.Task, 0)
	return tm
}

// startScheduling starts the tasks and keep running them at regular interval.
func (tm *TaskManager) startScheduling(ctx context.Context) {
	// task manager runs in loop at regular interval of 1 seconds to schedule tasks.
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			tm.runReadyTasks()
		case <-ctx.Done():
			log.Infof("stopped task manager")
			return
		}
	}
}

// AddTask adds a task to task queue.
func (tm *TaskManager) AddTask(t *task.Task) {
	tm.taskQueueMutex.Lock()
	defer tm.taskQueueMutex.Unlock()
	tm.taskQueue = append(tm.taskQueue, t)
}

// getReadyTasks returns the tasks that are ready to be executed.
func (tm *TaskManager) getReadyTasks() []*task.Task {
	readyTaskQueue := make([]*task.Task, 0)
	tm.taskQueueMutex.Lock()
	defer tm.taskQueueMutex.Unlock()

	if len(tm.taskQueue) < 1 {
		log.Debugf("No scheduled tasks.")
		return readyTaskQueue
	}
	for _, task := range tm.taskQueue {
		currentTime := time.Now()
		// if time passed the nextruntime , immediately schedule it
		isTaskReady := currentTime.After(task.NextRunTime)
		isRunning := atomic.LoadInt64(&task.IsRunning)
		if !isTaskReady || isRunning == 1 {
			continue
		}
		atomic.StoreInt64(&task.IsRunning, 1)
		task.ScheduleNextRunOfTask()
		task.PreviousRunTime = currentTime
		readyTaskQueue = append(readyTaskQueue, task)
	}
	return readyTaskQueue
}

// runReadyTasks looks at the readyTasks and schedules them. If there are
// more than one ready task , these will execute in parallel. The max number
// of tasks that will be executed at any point is determined by maxConcurrencyGuard.
func (tm *TaskManager) runReadyTasks() {
	readyTasks := tm.getReadyTasks()
	if len(readyTasks) < 1 {
		log.Debug("no ready tasks to schedule")
		return
	}

	for _, t := range readyTasks {
		tm.maxConcurrency <- struct{}{}
		go func(t *task.Task) {
			log.Debugf("running task with id %v", t.ID)
			err := t.Start()
			if err != nil {
				log.Error("error[%v] in running task with id %v", err, t.ID)
			}
			atomic.StoreInt64(&t.IsRunning, 0)
			<-tm.maxConcurrency
			log.Debugf("completed task with id %v", t.ID)
		}(t)
	}
}
