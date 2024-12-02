package loop

import (
	"context"
)

func (l *Loop) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error)

	for _, r := range l.runners {
		go func(r runner) {
			errCh <- r.Run(ctx)
		}(r)
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
