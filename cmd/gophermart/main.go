package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	_ "github.com/sodiqit/gophermart/docs"
	"github.com/sodiqit/gophermart/internal/server/config"
	"github.com/sodiqit/gophermart/internal/server/infra"
	"github.com/sodiqit/gophermart/internal/server/infra/grpc"
	"github.com/sodiqit/gophermart/internal/server/infra/http"
)

//	@Title			GopherMart API
//	@Description	Сервис накопительный системы.
//	@Version		1.0

//	@BasePath	/api/

//	@SecurityDefinitions.apikey	ApiKeyAuth
//	@In							header
//	@Name						Authorization

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Println("Build version:", checkVarBuild(buildVersion))
	fmt.Println("Build date:", checkVarBuild(buildDate))
	fmt.Println("Build commit:", checkVarBuild(buildCommit))

	config := config.ParseConfig()

	ctx, cancel := context.WithCancel(context.Background())

	deps, err := infra.NewAppContainer(ctx, config)

	if err != nil {
		log.Fatalf("Error while init app deps: %s", err)
	}

	c := make(chan os.Signal, 1)

	quit := make(chan struct{})

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		grpc.StopServer()
		err := http.StopServer(context.Background())
		if err != nil {
			log.Fatalln("Error while shutdown server", err.Error())
		}
		close(quit)
	}()

	go func() {
		err := grpc.RunServer(ctx, deps)
		if err != nil {
			log.Fatalln("Error while run grpc server", err.Error())
		}
	}()

	err = http.RunServer(ctx, deps)
	if err != nil && !errors.Is(err, netHttp.ErrServerClosed) {
		fmt.Println("Error while running gophermart server", err.Error())
		close(quit)
	}

	<-quit
}

func checkVarBuild(s string) string {
	if s == "" {
		return "N/A"
	}

	return s
}
