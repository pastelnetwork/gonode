package artwork

// const logPrefix = "[artwork]"

// Service represent artwork service.
type Service struct {
}

// Register regusters a new artwork.
func (service *Service) Register(artwork *Artwork) error {
	return nil
}

// NewService returns a new Service instance.
func NewService() *Service {
	return &Service{}
}
