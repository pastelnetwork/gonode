package services

import "context"

type Common struct{}

func (service *Common) Run(ctx context.Context) error {
	return nil
}

func NewCommon() *Common {
	return &Common{}
}
