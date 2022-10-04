package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rinuccia/jwt-auth/routes"
	"github.com/rinuccia/jwt-auth/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	// Run server
	srv := new(server.Server)
	router := gin.New()
	routes.InitRoutes(router)

	go func() {
		if err := srv.Run(router, os.Getenv("PORT")); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	// Waiting signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
