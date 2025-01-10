package cmd

import (
	"github.com/spf13/cobra"

	"shopsvc/cmd/schema"
)

func Execute() {
	var cmd = &cobra.Command{
		Short: "Shop Services",
		Long:  "Shop Service - Backend",
	}

	// Add command
	cmd.AddCommand(schema.Migrate)
	cmd.AddCommand(schema.Seeder)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
