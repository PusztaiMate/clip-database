package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PusztaiMate/clip-database/api"
	"github.com/PusztaiMate/clip-database/db"
	"github.com/PusztaiMate/clip-database/service"
)

const (
	address = ":8080"
)

func main() {
	logger := log.New(os.Stdout, "[CLIP-SERVER]", log.LstdFlags)
	store := db.NewInMemoryClipStore()
	service := service.NewClipperService(store)
	server := api.NewServer(service, logger)

	logger.Println("Starting up server on ", address)
	logger.Fatal(http.ListenAndServe(address, server))
}
