package app

import (
	"Classroom/Auth/internal/config"
	"fmt"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type app struct {
	logger *slog.Logger
	srv    *grpc.Server
	addr   string
}

type Controller interface {
	Init(srv *grpc.Server)
}

func New(logger *slog.Logger, conf *config.Config, controllers ...Controller) *app {
	srv := grpc.NewServer()

	for _, c := range controllers {
		c.Init(srv)
	}

	addr := fmt.Sprintf(":%d", conf.Port)
	return &app{logger: logger, srv: srv, addr: addr}
}

// Запускает gRPC сервер
func (a *app) Start() {
	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	a.logger.Info("starting grpc server", "addr", lis.Addr())
	if err := a.srv.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}

// Останавливает gRPC сервер
func (a *app) Stop() {
	a.srv.GracefulStop()
	a.logger.Info("grpc server stopped")
}
