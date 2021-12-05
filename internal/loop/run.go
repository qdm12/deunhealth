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

	var subLoops int

	subLoops++
	go func() {
		errCh <- l.unhealthy.Run(ctx)
	}()

	for i := 0; i < subLoops; i++ {
		workerErr := <-errCh
		if ctx.Err() == nil && workerErr != nil {
			err = workerErr
			cancel()
		}
	}

	return err
}
