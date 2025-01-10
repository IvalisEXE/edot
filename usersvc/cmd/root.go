package cmd

import (
	"github.com/spf13/cobra"

	"usersvc/cmd/api"
	"usersvc/cmd/schema"
)

func Execute() {
	var cmd = &cobra.Command{
		Short: "User Services",
		Long:  "User Service - Backend",
	}

	// Add command
	cmd.AddCommand(api.API)
	cmd.AddCommand(schema.Migrate)
	cmd.AddCommand(schema.Seeder)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
