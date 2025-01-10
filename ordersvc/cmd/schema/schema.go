package schema

import (
	"ordersvc/internal/infrastructure/repository"
	"ordersvc/pkg/config"
	"ordersvc/pkg/database"

	"github.com/spf13/cobra"
)

var cfg *config.Container

func init() {
	// Load configuration
	cfg = config.New(".")
}

// Migrate schema the database
var Migrate = &cobra.Command{
	Use:   "migrate",
	Short: "Run migration",
	Run: func(cmd *cobra.Command, args []string) {

		// Create a new database connection
		database.Load(cfg.DB)
		defer database.Close()

		err := database.DB.AutoMigrate(
			&repository.Order{},
		)
		if err != nil {
			panic(err)
		}
	},
}
