package static

import (
	"context"
	"embed"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services"
)

//go:embed data/*.json
var fs embed.FS

// service represents a service that returns a static data.
type service struct {
	*services.Common
}

func (service *service) Handle(ctx context.Context, method string, params []string) (interface{}, error) {
	switch service.RoutePath(method, params) {
	case "masternode_top":
		return service.topMasterNode()
	case "storagefee_getnetworkfee":
		return service.storageFeeGetNetworkFee()
	}

	return nil, services.ErrNotFoundMethod
}

func (service *service) topMasterNode() (interface{}, error) {
	var response TopMasterNodes

	if err := service.UnmarshalFile(fs, "data/masternode_top.json", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (service *service) storageFeeGetNetworkFee() (interface{}, error) {
	var response StorageFeeGetNetworkFee

	if err := service.UnmarshalFile(fs, "data/storagefee_getnetworkfee.json", &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// New returns a new Service instance.
func New() services.Service {
	return &service{
		Common: services.NewCommon(),
	}
}
