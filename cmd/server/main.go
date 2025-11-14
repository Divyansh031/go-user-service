package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Divyansh031/user-service/internal/config"
	"github.com/Divyansh031/user-service/internal/grpc/handlers"
	"github.com/Divyansh031/user-service/internal/storage/scylla"
	pb "github.com/Divyansh031/user-service/api/proto/user/v1"
	"google.golang.org/grpc"
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
	reflection.Register(grpcServer)

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
	mux := http.NewServeMux()

	// Basic health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// REST API endpoints info
	mux.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"User Service API - Use gRPC for operations"}`))
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("HTTP REST server listening", "port", cfg.HTTP.Port)

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server error", "error", err)
	}
}
