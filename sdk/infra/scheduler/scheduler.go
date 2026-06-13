package scheduler

import (
	"log"
	"sync"
	"time"

	"hei-gin/sdk/kernel/plugin"

	"github.com/robfig/cron/v3"
)

type Task interface {
	Name() string
	Run()
}

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

type Scheduler struct {
	cron *cron.Cron
}

var defaultScheduler = New()

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(cron.Recover(cron.DefaultLogger)),
		),
	}
}

func Register(spec string, task Task) (cron.EntryID, error) {
	return defaultScheduler.Register(spec, task)
}

func RegisterInterval(d time.Duration, task Task) (cron.EntryID, error) {
	return defaultScheduler.RegisterInterval(d, task)
}

func (s *Scheduler) Register(spec string, task Task) (cron.EntryID, error) {
	id, err := s.cron.AddJob(spec, &taskWrapper{task: task})
	if err != nil {
		log.Printf("[Scheduler] Failed to register task %s with spec %q: %v", task.Name(), spec, err)
		return 0, err
	}
	log.Printf("[Scheduler] Registered task: %s [spec=%q]", task.Name(), spec)
	return id, nil
}

func (s *Scheduler) RegisterInterval(d time.Duration, task Task) (cron.EntryID, error) {
	return s.Register("@every "+d.String(), task)
}

func Start() { defaultScheduler.Start() }
func Stop()  { defaultScheduler.Stop() }

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Printf("[Scheduler] Started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Printf("[Scheduler] Stopped")
}

type schedulerPlugin struct{ plugin.NoopPlugin }

var registerOnce sync.Once

func (m *schedulerPlugin) Name() string { return "scheduler" }
func (m *schedulerPlugin) Start() error { Start(); return nil }
func (m *schedulerPlugin) Stop() error  { Stop(); return nil }

func RegisterPlugin() {
	registerOnce.Do(func() {
		plugin.Register(&schedulerPlugin{})
	})
}
