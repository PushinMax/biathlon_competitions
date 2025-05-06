package main

import (
	"biathlon/internal/handler"
	"biathlon/internal/repository"
	"biathlon/internal/service"
	"fmt"
)

func main() {
	fileName := "configs/events"

	repo := repository.New("configs/config.json")
	service := service.New(repo)
	handler := handler.New(service, fileName)
	
	if err := handler.Start(); err != nil {
		fmt.Println(err)
	}
}