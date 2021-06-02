package cmd

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
)

type service interface {
	Run(context.Context) error
}

func runServices(ctx context.Context, services ...service) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, service := range services {
		service := service

		group.Go(func() error {
			return service.Run(ctx)
		})
	}

	return group.Wait()
}
