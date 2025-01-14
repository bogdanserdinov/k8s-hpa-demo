package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	subtractionpb "example/gen/go/subtraction/v1"
	"example/pkg/greceful/shutdown"
	"example/pkg/http/cors"
	service "example/services/subtraction/v1"
)

type Config struct {
	GrpcServerPort        int `env:"GRPC_SERVER_PORT,required"`
	HealthcheckServerPort int `env:"HEALTHCHECK_SERVER_PORT,required"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown.Graceful(func() {
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

	{ // GRPC endpoints
		const defaultGrpcMessageSize = 2048

		grpcServer := grpc.NewServer(
			grpc.MaxRecvMsgSize(defaultGrpcMessageSize),
			grpc.MaxSendMsgSize(defaultGrpcMessageSize),
		)

		subtractionServer := service.New()
		subtractionpb.RegisterSubtractionServiceServer(grpcServer, subtractionServer)

		logger.Info("starting subtraction grpc server")

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

	{ // HTTP endpoints.
		httpServerListen := fmt.Sprintf("0.0.0.0:%v", cfg.HealthcheckServerPort)

		router := mux.NewRouter()
		router.HandleFunc("/heathz", func(writer http.ResponseWriter, request *http.Request) {}).Methods(http.MethodGet)

		httpServer := &http.Server{
			Addr:              httpServerListen,
			Handler:           cors.Allow(router),
			ReadHeaderTimeout: 2 * time.Second,
		}

		logger.Info("starting subtraction healthz server")

		group.Go(func() error {
			<-ctx.Done()
			return httpServer.Shutdown(ctx)
		})
		group.Go(func() error {
			err := httpServer.ListenAndServe()
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			return err
		})
	}

	logger.Info("The subtraction-service was terminated with: %v", zap.Error(group.Wait()))
}
