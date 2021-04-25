package cmd

import (
	"context"

	"github.com/pastelnetwork/go-commons/errors"
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
			defer errors.Recover(func(rec error) { err = rec })
			return service.Run(ctx)
		})
	}

	return group.Wait()
}
