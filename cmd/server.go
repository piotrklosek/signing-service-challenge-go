package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence/inmemory"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence/mongo"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence/postgres"
	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/logger"
	"go.uber.org/zap"

	"github.com/piotrklosek/signing-service-challenge-go/internal/api"
	"github.com/piotrklosek/signing-service-challenge-go/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the Signature Service API server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		var (
			deviceRepo    persistence.DeviceRepository
			signatureRepo persistence.SignatureRepository
			userRepo      persistence.UserRepository
			err           error
		)
		shutdownCtx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
		// create logger
		logger := logger.CreateLogger("server")

		switch cfg.DBType {
		case "postgres":
			deviceRepo, signatureRepo, userRepo, err = postgres.NewRepositories(cfg.Postgres.DSN)
		case "mongo":
			deviceRepo, signatureRepo, userRepo, err = mongo.NewRepositories(cfg.Mongo.URI, cfg.Mongo.Database)
		default: // inmemory
			inMemoryStore, err := inmemory.NewMemoryStore(cfg.InMemory.DBFilePath)
			if err != nil {
				logger.Fatal("enabled to create inmemory store", zap.Error(err))
				panic(err)
			}
			deviceRepo = inMemoryStore.DeviceRepo
			signatureRepo = inMemoryStore.SignatureRepo
			userRepo = inMemoryStore.UserRepo
			defer inMemoryStore.SaveOnShutdown(shutdownCtx)
		}

		if err != nil {
			logger.Fatal("failed to initialize storage", zap.Error(err))
		}

		router := api.NewRouter(deviceRepo, signatureRepo, userRepo)

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.Port),
			Handler: router,
		}

		go func() {
			logger.Info("Signature Service running",
				zap.String("port", cfg.Port), zap.String("dbType", cfg.DBType))
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("server failed", zap.Error(err))
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("Shutting down server...")
		// dump db backup

		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Fatal("Server forced to shutdown", zap.Error(err))
		}
		logger.Info("Server stopped gracefully")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Flags specific to the server command
	serverCmd.Flags().String("port", "8080", "Port for HTTP server")
	_ = viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))

	serverCmd.Flags().String("db.type", "memory", "Database type (memory|postgres|mongo)")
	_ = viper.BindPFlag("db.type", serverCmd.Flags().Lookup("db.type"))

	serverCmd.Flags().String("db.postgres.dsn", "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable", "Postgres driver connection string (DSN)")
	_ = viper.BindPFlag("db.postgres.dsn", serverCmd.Flags().Lookup("db.postgres.dsn"))

	serverCmd.Flags().String("db.inmemory.filepath", "./database/inmemory/database.json", "Database file path")
	_ = viper.BindPFlag("db.inmemory.dbfilepath", serverCmd.Flags().Lookup("db.inmemory.filepath"))
}
