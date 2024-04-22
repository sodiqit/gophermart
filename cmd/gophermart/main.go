package main

import (
	"fmt"

	_ "net/http/pprof"

	_ "github.com/sodiqit/gophermart/docs"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/infra/http"
)

//	@Title			GopherMart API
//	@Description	Сервис накопительный системы.
//	@Version		1.0

//	@BasePath	/api/

//	@SecurityDefinitions.apikey	ApiKeyAuth
//	@In							header
//	@Name						Authorization

func main() {
	config := config.ParseConfig()

	err := http.RunServer(config)
	if err != nil {
		fmt.Println("Error while running gophermart server", err.Error())
	}
}
