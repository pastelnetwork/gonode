package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bitwurx/jrpc2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pastelnetwork/gonode/integration/fakes/common/register"
	"github.com/pastelnetwork/gonode/integration/fakes/common/storage"
	"github.com/pastelnetwork/gonode/integration/fakes/pasteld/handler"
)

func main() {
	store := storage.New(10*time.Minute, 10*time.Minute)
	rpcHandler := handler.New(store)
	registrationHandler := register.New(store)

	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	router.GET("/health", healthGET())
	router.POST("/register", registrationHandler.Register())
	router.POST("/cleanup", registrationHandler.Cleanup())

	// Start Server
	go func() {
		if err := router.Run(":29932"); err != nil {
			fmt.Printf("Run Server Failed: err %v\n", err)
			return
		}
	}()

	// create a new server instance
	s := jrpc2.NewServer(":19932", "/", nil)

	// register the add method
	s.Register("masternode", jrpc2.Method{Method: rpcHandler.HandleMasternode})
	s.Register("storagefee", jrpc2.Method{Method: rpcHandler.HandleStorageFee})
	s.Register("tickets", jrpc2.Method{Method: rpcHandler.HandleTickets})
	s.Register("getblockcount", jrpc2.Method{Method: rpcHandler.HandleGetBlockCount})
	s.Register("getblock", jrpc2.Method{Method: rpcHandler.HandleGetBlock})
	s.Register("pastelid", jrpc2.Method{Method: rpcHandler.HandlePastelid})
	s.Register("z_getbalance", jrpc2.Method{Method: rpcHandler.HandleZGetBalance})
	s.Register("z_getoperationstatus", jrpc2.Method{Method: rpcHandler.HandleGetOperationStatus})
	s.Register("z_sendmanywithchangetosender", jrpc2.Method{Method: rpcHandler.HandleSendMany})
	s.Register("getrawtransaction", jrpc2.Method{Method: rpcHandler.HandleGetRawTransaction})

	// start the server instance
	s.Start()
}

func healthGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "fake-pasteld")
	}
}
