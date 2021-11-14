package main

import (
	"chat/v1/config"
	"chat/v1/pkg/controller"
	"chat/v1/pkg/router"
	"log"
	"os"
)

func main() {
	port := config.GetAPIPort()
	ctrl := controller.NewController()
	router := router.NewRouter(ctrl)

	if err := router.Run(port); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
