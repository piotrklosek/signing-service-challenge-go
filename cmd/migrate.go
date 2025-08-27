package cmd

import (
	"fmt"
	"log"

	"github.com/piotrklosek/signing-service-challenge-go/internal/config"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence/postgres"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations (Postgres only)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		if cfg.DBType != "postgres" {
			log.Fatal("migrate command is only supported for Postgres backend")
		}

		if err := postgres.RunMigrations(cfg.Postgres.DSN); err != nil {
			log.Fatalf("migration failed: %v", err)
		}

		fmt.Println("Migrations applied successfully")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
