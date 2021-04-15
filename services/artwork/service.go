package artwork

const logPrefix = "[artwork]"

type Service struct {
}

func (service *Service) Register(artwork *Artwork) error {
	return nil
}

func NewService() *Service {
	return &Service{}
}
