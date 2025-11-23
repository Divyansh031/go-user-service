package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"syscall"


	pb "github.com/Divyansh031/user-service/api/proto/user/v1"
	"github.com/Divyansh031/user-service/internal/config"
	"github.com/Divyansh031/user-service/internal/grpc/handlers"
	"github.com/Divyansh031/user-service/internal/storage/scylla"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load configuration
	cfg := config.MustLoad()

	// Setup logger
	logLevel := slog.LevelInfo
	if cfg.Log.Level == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting user service", "env", cfg.Env)

	// Initialize database
	slog.Info("Initializing ScyllaDB", "hosts", cfg.ScyllaDB.Hosts, "keyspace", cfg.ScyllaDB.Keyspace)
	db, err := scylla.NewScyllaDB(cfg.ScyllaDB.Hosts, cfg.ScyllaDB.Port, cfg.ScyllaDB.Keyspace, cfg.ScyllaDB.Consistency)
	if err != nil {
		slog.Error("Failed to initialize ScyllaDB", "error", err)
		log.Fatal(err)
	}
	defer db.Close()

	slog.Info("ScyllaDB initialized successfully")

	// Start gRPC server
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		slog.Error("Failed to listen on gRPC port", "port", cfg.GRPC.Port, "error", err)
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	userServiceServer := handlers.NewUserServiceServer(db)
	pb.RegisterUserServiceServer(grpcServer, userServiceServer)

	// Register reflection for grpcurl
	reflection.Register(grpcServer)  // Allows grpcurl to inspect your API

	slog.Info("gRPC server listening", "port", cfg.GRPC.Port)

	// Start gRPC server in goroutine
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	// Start HTTP/REST server
	go startRESTServer(cfg)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	slog.Info("Shutdown signal received, gracefully shutting down...")

	// Shutdown gRPC server
	grpcServer.GracefulStop()

	slog.Info("User service stopped")
}

// startRESTServer starts a simple HTTP REST server
func startRESTServer(cfg *config.Config) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", cfg.GRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("Failed to connect to gRPC server", "error", err)
		return
	}
	defer conn.Close()

	// Gateway mux (This expects /v1/users)
	gwMux := runtime.NewServeMux()

	client := pb.NewUserServiceClient(conn)
	if err := pb.RegisterUserServiceHandlerClient(ctx, gwMux, client); err != nil {
		slog.Error("Failed to register gateway handler", "error", err)
		return
	}

	// NEW MUX — This will handle /api/v1/users
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", gwMux))
	mux.Handle("/", gwMux) // fallback for /v1/users (if someone uses directly)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: mux,
	}

	slog.Info("HTTP REST server listening", "port", cfg.HTTP.Port)
	slog.Info("REST API → POST http://localhost:8080/api/v1/users")

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server error", "error", err)
	}
}