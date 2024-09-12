package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	additionpb "example/gen/go/addition/v1"
)

type Config struct {
	AdditionalServerAddr string `env:"ADDITIONAL_SERVER_ADDR,required"`
	HttpPublicServerPort int    `env:"HTTP_PUBLIC_SERVER_PORT,required"`
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

	httpServerListen := fmt.Sprintf("0.0.0.0:%v", cfg.HttpPublicServerPort)
	log.Println(httpServerListen)

	const defaultGrpcMessageSize = 2048

	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultGrpcMessageSize),
			grpc.MaxCallSendMsgSize(defaultGrpcMessageSize),
		),
	}

	if err := additionpb.RegisterAdditionServiceHandlerFromEndpoint(ctx, grpcMux, cfg.AdditionalServerAddr, opts); err != nil {
		log.Fatalf("failed to register addition service: %v", err)
	}

	router := mux.NewRouter()
	router.PathPrefix("/").Handler(
		http.StripPrefix("/", grpcMux),
	)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	httpServer := &http.Server{
		Addr:              httpServerListen,
		Handler:           allowCORS(router),
		ReadHeaderTimeout: 2 * time.Second,
	}

	logger.Info("starting gateway server")

	var group errgroup.Group

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

	logger.Info("The gateway-service was terminated with: %v", zap.Error(group.Wait()))
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		h.ServeHTTP(w, r)
	})
}

func gracefulShutdown(actions func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		actions()
	}()
}
