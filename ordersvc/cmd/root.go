package cmd

import (
	"github.com/spf13/cobra"

	"ordersvc/cmd/api"
	"ordersvc/cmd/cron"
	"ordersvc/cmd/schema"
)

func Execute() {
	var cmd = &cobra.Command{
		Short: "Order Services",
		Long:  "Order Service - Backend",
	}

	// Add command
	cmd.AddCommand(api.API)
	cmd.AddCommand(schema.Migrate)
	cmd.AddCommand(cron.Cron)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
