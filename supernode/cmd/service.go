package cmd

import (
	"context"

	"github.com/pastelnetwork/gonode/common/errgroup"
	"github.com/pastelnetwork/gonode/common/log"
)

type service interface {
	Run(context.Context) error
}

func runServices(ctx context.Context, services ...service) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, service := range services {
		service := service

		group.Go(func() error {
			err := service.Run(ctx)
			log.WithContext(ctx).WithError(err).Error("Server stop error")
			return err
		})
	}

	return group.Wait()
}
