package main

import (
	"biathlon/internal/handler"
	"biathlon/internal/repository"
	"biathlon/internal/service"
	"fmt"
	"flag"
)

func main() {
	eventsPath := flag.String("event", "configs/events", "path to file with events")
	cfg := flag.String("config", "configs/config.json", "path to config file")
	
	flag.Parse()

	repo := repository.New(*cfg)
	service := service.New(repo)
	handler := handler.New(service, *eventsPath)
	
	if err := handler.Start(); err != nil {
		fmt.Println(err)
	}
}