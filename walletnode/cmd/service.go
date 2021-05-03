package cmd

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errors"
	"golang.org/x/sync/errgroup"
)

type service interface {
	Run(context.Context) error
}

func runServices(ctx context.Context, services ...service) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, service := range services {
		service := service

		group.Go(func() (err error) {
			defer errors.Recover(func(recErr error) { err = recErr })
			return service.Run(ctx)
		})
	}

	return group.Wait()
}
