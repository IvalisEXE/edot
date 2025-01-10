package schema

import (
	"fmt"
	"warehousesvc/internal/infrastructure/repository"
	"warehousesvc/pkg/config"
	"warehousesvc/pkg/database"

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
			&repository.Warehouse{},
			&repository.WarehouseStockItem{},
			&repository.WarehouseStockTransfer{},
		)
		if err != nil {
			panic(err)
		}
	},
}

// Seed the database
var Seeder = &cobra.Command{
	Use:   "seed",
	Short: "Run database seeders",
	Run: func(cmd *cobra.Command, args []string) {

		// Create a new database connection
		database.Load(cfg.DB)
		defer database.Close()

		var seeders = map[string]database.Seeder{
			"warehouse":            repository.InitSeeder(database.DB),
			"warehouse_stock_item": repository.InitSeeder(database.DB),
		}

		// Seed the database
		if len(args) > 0 {
			// Seed specific seeders
			for _, arg := range args {
				if seeder, ok := seeders[arg]; ok {
					if err := seeder.Seed(); err != nil {
						panic(err)
					}
				} else {
					fmt.Printf("Seeder '%s' not found\n", arg)
				}
			}
		} else {
			// Seed all
			for _, s := range seeders {
				if err := s.Seed(); err != nil {
					panic(err)
				}
			}
		}
	},
}
