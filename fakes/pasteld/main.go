package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bitwurx/jrpc2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pastelnetwork/gonode/fakes/common/register"
	"github.com/pastelnetwork/gonode/fakes/common/storage"
	"github.com/pastelnetwork/gonode/fakes/pasteld/handler"
)

func main() {
	store := storage.New(1*time.Minute, 2*time.Minute)
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

	// start the server instance
	s.Start()
}

func healthGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "fake-pasteld")
	}
}
