package cmd

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type service interface {
	Run(context.Context) error
}

func run(ctx context.Context, services ...service) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, service := range services {
		service := service
		group.Go(func() error {
			return service.Run(ctx)
		})
	}

	err := group.Wait()
	return err
}
