package main

import (
	"biathlon/internal/handler"
	"biathlon/internal/service"
	"fmt"
)

func main() {
	fileName := "configs/events"

	service := service.New()
	handler := handler.New(service, fileName)
	
	if err := handler.Start(); err != nil {
		fmt.Println(err)
	}
}