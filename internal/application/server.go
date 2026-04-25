package application

import "context"

type Server interface {
	Run(ctx context.Context) error
}
