package loop

import (
	"context"
)

type Runner interface {
	Run(ctx context.Context) (err error)
}

func (l *Loop) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error)

	for _, runner := range l.runners {
		go func(runner Runner) {
			errCh <- runner.Run(ctx)
		}(runner)
	}

	for range l.runners {
		workerErr := <-errCh
		if ctx.Err() == nil && workerErr != nil {
			err = workerErr
			cancel()
		}
	}

	return err
}
