package main

import (
	"log"

	"github.com/polyk005/notesync"
	"github.com/polyk005/notesync/pkg/handler"
	"github.com/polyk005/notesync/pkg/repository"
	"github.com/polyk005/notesync/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(notesync.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
