package loop

import "context"

type runner interface {
	Run(ctx context.Context) (err error)
}
