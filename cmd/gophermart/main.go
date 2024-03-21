package main

import (
	"fmt"

	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/infra/http"
)

func main() {
	config := config.ParseConfig()

	err := http.RunServer(config)
	if err != nil {
		fmt.Println("Error while running gophermart server", err.Error())
	}
}
