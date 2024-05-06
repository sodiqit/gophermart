package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/sodiqit/gophermart/gen/proto/auth"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/infra"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var srv *grpc.Server

func RunServer(ctx context.Context, deps *infra.AppContainer) error {
	listen, err := net.Listen("tcp", deps.Config.GRPCAddress)
	if err != nil {
		log.Fatal(err)
	}

	logger := deps.Logger
	config := deps.Config

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Errorw("Recovered from panic", "panic", p)

			return status.Errorf(codes.Internal, "Internal server error")
		}),
	}

	srv = grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logger), []logging.Option{}...),
	))

	authv1.RegisterAuthServiceServer(srv, deps.AuthContainer.GRPCServer)

	logger.Infow("start gRPC server", "port", config.GRPCAddress)

	return srv.Serve(listen)
}

func StopServer() error {
	fmt.Println("\nGracefully shutdown gRPC server")
	srv.GracefulStop()
	return nil
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		args := []any{"msg", msg}
		l.Log(zapcore.Level(lvl), append(args, fields...))
	})
}
