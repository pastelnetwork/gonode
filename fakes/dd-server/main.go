package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/pastelnetwork/gonode/fakes/common/register"
	"github.com/pastelnetwork/gonode/fakes/common/storage"
	"github.com/pastelnetwork/gonode/fakes/dd-server/dupedetection"
	"github.com/pastelnetwork/gonode/fakes/dd-server/server"
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
		dupedetection.NewDDService(store),
	)

	if err := grpc.Run(context.Background()); err != nil {
		panic("unable to run grpc server")
	}

}

func healthGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "fake-dd-server")
	}
}
