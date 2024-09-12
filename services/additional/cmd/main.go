package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	additionpb "example/gen/go/additional/v1"
	"example/services/additional"
)

type Config struct {
	GrpcServerPort int `env:"GRPC_SERVER_PORT,required"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	gracefulShutdown(func() {
		cancel()
	})

	loggerCfg := zap.NewProductionConfig()
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatalln("could not build logger from config", err)
		return
	}
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Fatalln("could not sync logger", err)
		}
	}()

	cfg := new(Config)
	err = env.Parse(cfg)
	if err != nil {
		logger.Fatal("could not parse config", zap.Error(err))
	}

	grpcServerListen := fmt.Sprintf("0.0.0.0:%v", cfg.GrpcServerPort)

	var group errgroup.Group

	// GRPC endpoints
	{
		const defaultGrpcMessageSize = 2048

		grpcServer := grpc.NewServer(
			grpc.MaxRecvMsgSize(defaultGrpcMessageSize),
			grpc.MaxSendMsgSize(defaultGrpcMessageSize),
		)

		additionServer := additional.New()
		additionpb.RegisterAdditionServiceServer(grpcServer, additionServer)

		logger.Info("starting additional server")

		group.Go(func() error {
			<-ctx.Done()
			grpcServer.GracefulStop()
			return nil
		})
		group.Go(func() error {
			lis, err := net.Listen("tcp", grpcServerListen)
			if err != nil {
				return fmt.Errorf("failed to listen: %w", err)
			}

			return grpcServer.Serve(lis)
		})
	}

	logger.Info("The additional-service was terminated with: %v", zap.Error(group.Wait()))
}

func gracefulShutdown(action func()) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		action()
	}()
}
