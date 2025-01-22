package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver
	"github.com/oklog/run"
	"google.golang.org/grpc"

	"apify-poi-data/config"
	database "apify-poi-data/db"
	sqlcdb "apify-poi-data/db/sqlc"
	"apify-poi-data/pkg/health"
)

var (
	cfg       *config.Config
	tlsConfig *tls.Config
)

var (
	grpcServer *grpc.Server
	httpServer *http.Server
)

var (
	db *sqlcdb.Database
)

func init() {
	var err error
	// Load Configuration
	cfg = config.LoadConfig()
	tlsConfig, err = loadTLSCredentials()
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}
}

func Run() {
	// Connect to database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful Shutdown Signal handler
	setupSignalHandler(cancel)

	// Health service
	healthService, err := health.NewHealthHandler()
	if err != nil {
		panic(fmt.Errorf("failed to create health service: %v", err))
	}

	var g run.Group

	// Health check for the database
	if err := healthCheck(cfg.Database.URL); err != nil {
		panic(fmt.Errorf("database health check failed: %v", err))
	}

	db, err = sqlcdb.NewDatabase(ctx, cfg.Database.URL)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database; err=%v", err))
	}
	defer db.Close()

	err = database.RunMigrations(
		cfg.Database.URL,
		cfg.Database.MigrationPath,
		cfg.Database.DatabaseName,
		cfg.Database.DatabaseVersion,
	)

	if err != nil {
		panic(fmt.Errorf("failed to run migrations: %v", err))
	}

	// Start gRPC server
	g.Add(func() error {
		log.Println("Starting GRPC server...")
		grpcServer, listener, err := setupGRPCServer(cfg.Ports.GRPCPort, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to setup GRPC server: %w", err)
		}
		return grpcServer.Serve(listener)
	}, func(err error) {
		log.Println("Shutting down GRPC server...")
		ShutdownGRPCServer(grpcServer)
		cancel()
	})

	// Start HTTP server
	g.Add(func() error {
		log.Println("Starting HTTP server...")
		httpServer, err = SetupHTTPMux(ctx, cfg.Ports.GRPCPort, cfg.Ports.HTTPPort, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to setup HTTP server: %w", err)
		}
		return httpServer.ListenAndServe()
	}, func(err error) {
		log.Println("Shutting down HTTP server...")
		err = httpServer.Shutdown(ctx)
		if err != nil {
			return
		}
		cancel()
	})

	// Start Health service
	g.Add(func() error {
		mux, err := healthService.ServeHealthcheckMux()
		if err != nil {
			return fmt.Errorf("failed to serve healthcheck: %w", err)
		}

		healthServer := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Ports.HealthPort),
			Handler: mux,
		}
		log.Printf("Health server listening on port %d...", cfg.Ports.HealthPort)
		return healthServer.ListenAndServe()
	}, func(err error) {
		log.Println("Shutting down Health server...")
		cancel()
	})

	// Context cancellation listener
	g.Add(func() error {
		<-ctx.Done()
		log.Println("Shutting down servers...")
		return nil
	}, func(err error) {
		cancel()
	})

	// Run application
	log.Println("Application started...")
	if err := g.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}

}

func healthCheck(databaseURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	if err = conn.Ping(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	log.Println("Database is up and running")
	return nil
}

func setupSignalHandler(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received interrupt signal; shutting down...")
		cancel()
	}()
}
