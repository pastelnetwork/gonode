package fake

import (
	"context"
	"embed"
	"net/http"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services"
	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"
)

//go:embed data/*.json
var fs embed.FS

// service represents a service that returns a static data.
type service struct {
	*services.Common
	topMasterNodes          models.TopMasterNodes
	idTicketRecords         models.IDTicketRecords
	storageFeeGetNetworkFee *models.StorageFeeGetNetworkFee
}

func (service *service) Handle(_ context.Context, r *http.Request, method string, params []string) (resp interface{}, err error) {
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
		return newTicketsListIDMine(service.idTicketRecords, *node), nil
	case "masternode_status":
		return newMasterNodeStatusByNode(node), nil
	}

	return nil, services.ErrNotFoundMethod
}

func (service *service) loadFiles() error {
	files := map[string]interface{}{
		"data/storagefee_getnetworkfee.json": &service.storageFeeGetNetworkFee,
		"data/masternode_top.json":           &service.topMasterNodes,
		"data/tickets_list_id_mine.json":     &service.idTicketRecords,
	}

	for filepath, obj := range files {
		if err := service.UnmarshalFile(fs, filepath, obj); err != nil {
			return err
		}
	}
	return nil
}

// New returns a new fake Service instance.
func New() (services.Service, error) {
	service := &service{
		Common: services.NewCommon(),
	}

	if err := service.loadFiles(); err != nil {
		return nil, err
	}
	return service, nil
}
