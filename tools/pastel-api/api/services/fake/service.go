package fake

import (
	"context"
	"embed"
	"net/http"
	"strings"

	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services"
	"github.com/pastelnetwork/gonode/tools/pastel-api/api/services/fake/models"
)

//go:embed data/*.json
var fs embed.FS

// service represents a service that returns a static data.
type service struct {
	*services.Common
	topMasterNodes          models.TopMasterNodes
	idTickets               models.IDTickets
	storageFeeGetNetworkFee *models.StorageFeeGetNetworkFee
}

func (service *service) Handle(_ context.Context, r *http.Request, method string, params []string) (resp interface{}, err error) {
	routePath := service.RoutePath(method, params)
	node := service.node(r)

	switch routePath {
	case "masternode_top":
		return newTopMasterNodes(service.topMasterNodes), nil
	case "storagefee_getnetworkfee":
		return newStorageFeeGetNetworkFee(service.storageFeeGetNetworkFee), nil
	case "tickets_list_id_mine":
		return newTicketsListIDMine(service.idTickets, node), nil
	case "masternode_list-conf":
		return newMasterNodeListConfByNode(node), nil
	case "masternode_status":
		if node == nil {
			return nil, services.ErrNotMasterNode
		}
		return newMasterNodeStatusByNode(node), nil
	default:
		switch {
		case strings.HasPrefix(routePath, "pastelid_sign_") && len(params) == 4:
			return newPastelID(params[2]).sign(params[1]), nil
		case strings.HasPrefix(routePath, "pastelid_verify_") && len(params) == 4:
			return newPastelID(params[3]).verify(params[1], params[2]), nil
		}
	}

	return nil, services.ErrNotFoundMethod
}

func (service *service) node(r *http.Request) *models.MasterNode {
	user, _, _ := r.BasicAuth()
	if user == "" {
		return nil
	}

	return service.topMasterNodes.LastBlock().ByPort(user)
}

func (service *service) loadFiles() error {
	files := map[string]interface{}{
		"data/storagefee_getnetworkfee.json": &service.storageFeeGetNetworkFee,
		"data/masternode_top.json":           &service.topMasterNodes,
		"data/tickets_list_id_mine.json":     &service.idTickets,
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
