package scheduler

import (
	"context"
	"time"
)

type Runner interface {
	RunDueChecks(context.Context)
}

type Scheduler struct {
	runner   Runner
	interval time.Duration
}

func New(runner Runner, interval time.Duration) Scheduler {
	return Scheduler{runner: runner, interval: interval}
}

func (s Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	s.runner.RunDueChecks(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.runner.RunDueChecks(ctx)
		}
	}
}
