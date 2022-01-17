package register

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pastelnetwork/gonode/fakes/pasteld/storage"
)

// Register is suppoed to provider handler functions to register mock
// responses to requests
type Register interface {
	Cleanup() gin.HandlerFunc
	Register() gin.HandlerFunc
}

type register struct {
	store storage.Store
}

// New returns  new instance of register
func New(store storage.Store) Register {
	return &register{
		store: store,
	}
}

func (r *register) mock(method string, params []string, data []byte) error {
	key := method + "*" + strings.Join(params, "*")
	return r.store.Set(key, data)
}

// Cleanup clears up storage
func (r *register) Cleanup() gin.HandlerFunc {
	return func(c *gin.Context) {
		r.store.Reset(1*time.Minute, 2*time.Minute)
		c.JSON(http.StatusOK, nil)
	}
}

// Register is handler func to register any mock request & response
// method and params are supposed to be passed as query params
func (r *register) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Query("method")
		params := c.Query("params")

		if method == "" || params == "" {
			c.JSON(http.StatusBadRequest, "method or params empty")
		} else {
			data, err := c.GetRawData()
			if err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			} else {
				if err := r.mock(method, strings.Split(params, ","), data); err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
				} else {
					c.JSON(http.StatusOK, nil)
				}
			}
		}
	}
}
