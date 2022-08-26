package process

import "context"

type Process interface {
	Start(ctx context.Context) error
	Stop() error
	Wait()
}
