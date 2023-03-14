package main

import (
	"context"
	"fmt"
	"github.com/pastelnetwork/gonode/common/log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/pastelnetwork/gonode/integration/fakes/common/register"
	"github.com/pastelnetwork/gonode/integration/fakes/common/storage"
	"github.com/pastelnetwork/gonode/integration/fakes/dd-server/handler"
	"github.com/pastelnetwork/gonode/integration/fakes/dd-server/server"
)

func main() {
	store := storage.New(1*time.Minute, 2*time.Minute)
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
		if err := router.Run(":51052"); err != nil {
			fmt.Printf("Run Server Failed: err %v\n", err)
			return
		}
	}()

	grpc := server.New(server.NewConfig(),
		"service",
		handler.NewDDService(store),
	)

	if err := grpc.Run(log.ContextWithPrefix(context.Background(), "")); err != nil {
		panic("unable to run grpc server")
	}

}

func healthGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "fake-dd-server")
	}
}
