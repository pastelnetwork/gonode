package services

import (
	"context"

	"github.com/pastelnetwork/go-commons/log"
	letters "github.com/pastelnetwork/walletnode/rest/gen/letters"
)

type lettersService struct {
	commonService
}

// Add a new letter and return its ID.
func (service *lettersService) Add(ctx context.Context, p *letters.AddPayload) (res string, err error) {
	log.WithContext(ctx).Debug("letters.add")
	return
}

// List all stored letters.
func (service *lettersService) List(ctx context.Context, p *letters.ListPayload) (res letters.StoredLetterCollection, err error) {
	log.WithContext(ctx).Debug("letters.list")
	return
}

// NewLetters returns the letters service implementation.
func NewLetters() letters.Service {
	return &lettersService{}
}
