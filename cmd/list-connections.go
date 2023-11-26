package cmd

import (
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/spf13/cobra"
)

var ListConnectionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the connections.",
	Long:  `List all database connections that are stored in the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.PrintStoredConnections()
	},
}
