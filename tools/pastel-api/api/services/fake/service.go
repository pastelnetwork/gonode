package fake

import (
	"context"
	"embed"
	"net/http"
	"sync"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services"
	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"
)

//go:embed data/*.json
var fs embed.FS

// service represents a service that returns a static data.
type service struct {
	*services.Common
	once                    sync.Once
	topMasterNodes          models.TopMasterNodes
	idTicketRecords         models.IDTicketRecords
	storageFeeGetNetworkFee *models.StorageFeeGetNetworkFee
}

func (service *service) Handle(_ context.Context, r *http.Request, method string, params []string) (interface{}, error) {
	service.once.Do(func() {
		service.loadFiles()
	})

	routePath := service.RoutePath(method, params)

	user, _, _ := r.BasicAuth()
	if user == "" {
		return nil, services.ErrUnauthorized
	}

	// walletnode

	switch routePath {
	case "masternode_top":
		return newTopMasterNodes(service.topMasterNodes), nil
	case "storagefee_getnetworkfee":
		return newStorageFeeGetNetworkFee(service.storageFeeGetNetworkFee), nil
	}

	// supernode

	node := service.topMasterNodes.LastBlock().ByPort(user)
	if node == nil {
		return nil, services.ErrNotFoundMethod
	}

	switch routePath {
	case "tickets_list_id_mine":
		return newTicketsListIDMine(*node, *personalPastelIDTicket), nil
	case "masternode_status":
		return newMasterNodeStatusByNode(node), nil
	}

	return nil, services.ErrNotFoundMethod
}

func (service *service) loadFiles() error {
	if err := service.UnmarshalFile(fs, "data/storagefee_getnetworkfee.json", &service.storageFeeGetNetworkFee); err != nil {
		return err
	}
	if err := service.UnmarshalFile(fs, "data/masternode_top.json", &service.topMasterNodes); err != nil {
		return err
	}
	if err := service.UnmarshalFile(fs, "data/tickets_list_id_mine.json", &service.idTicketRecords); err != nil {
		return err
	}

	return nil

}

// New returns a new fake Service instance.
func New() services.Service {
	return &service{
		Common: services.NewCommon(),
	}
}
