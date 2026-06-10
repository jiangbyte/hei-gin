package scheduler

import (
	"hei-gin/sdk/plugin"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Task represents a scheduled background task.
type Task interface {
	// Name returns the task name for logging.
	Name() string
	// Run executes the task.
	Run()
}

// taskWrapper adapts a Task to cron.Job with panic recovery and logging.
type taskWrapper struct {
	task Task
}

func (w *taskWrapper) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Scheduler] Task %s panicked: %v", w.task.Name(), r)
		}
	}()
	log.Printf("[Scheduler] Running task: %s", w.task.Name())
	w.task.Run()
}

// Scheduler manages cron-based scheduled tasks.
type Scheduler struct {
	cron *cron.Cron
}

var defaultScheduler = New()

// New creates a new Scheduler with seconds field enabled and panic recovery.
func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.Recover(cron.DefaultLogger)),
		),
	}
}

// Register adds a task with a cron expression to the default scheduler.
//
// The spec supports two formats:
//   - 6-field: "sec min hour dom mon dow" (default, seconds field enabled)
//   - 5-field: "min hour dom mon dow" (prepend "0 " automatically)
//
// Predefined schedules:
//
//	"@every 5m"  — runs every 5 minutes
//	"@daily"     — runs at midnight
//	"@hourly"    — runs at the start of every hour
func Register(spec string, task Task) (cron.EntryID, error) {
	return defaultScheduler.Register(spec, task)
}

// RegisterInterval registers a task that runs at a fixed interval on the default scheduler.
func RegisterInterval(d time.Duration, task Task) (cron.EntryID, error) {
	return defaultScheduler.RegisterInterval(d, task)
}

// Register adds a task with a cron expression to this scheduler.
func (s *Scheduler) Register(spec string, task Task) (cron.EntryID, error) {
	id, err := s.cron.AddJob(spec, &taskWrapper{task: task})
	if err != nil {
		log.Printf("[Scheduler] Failed to register task %s with spec %q: %v", task.Name(), spec, err)
		return 0, err
	}
	log.Printf("[Scheduler] Registered task: %s [spec=%q]", task.Name(), spec)
	return id, nil
}

// RegisterInterval registers a task that runs at a fixed interval.
func (s *Scheduler) RegisterInterval(d time.Duration, task Task) (cron.EntryID, error) {
	return s.Register("@every "+d.String(), task)
}

// Start starts the default scheduler.
func Start() { defaultScheduler.Start() }

// Stop gracefully stops the default scheduler and waits for running tasks to finish.
func Stop() { defaultScheduler.Stop() }

// Start starts the scheduler. It is safe to call multiple times.
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Printf("[Scheduler] Started")
}

// Stop gracefully stops the scheduler, waiting for running tasks to complete.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Printf("[Scheduler] Stopped")
}


// ---- plugin registration ----

type schedulerPlugin struct{ plugin.NoopPlugin }

func (m *schedulerPlugin) Name() string { return "scheduler" }

func (m *schedulerPlugin) Start() error { Start(); return nil }

func (m *schedulerPlugin) Stop() error { Stop(); return nil }

func init() { plugin.Register(&schedulerPlugin{}) }
