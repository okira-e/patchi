package cmd

import (
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/utils/logger"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/spf13/cobra"
)

var ListConnectionsCmd = &cobra.Command{
	Use:   "list-connections",
	Short: "List all the connections.",
	Long:  `List all database connections that are stored in the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, errOpt := config.GetUserConfig()
		if errOpt.IsSome() {
			logger.PrintInColor(colors.Red, "Config file not found. Run `patchi init` to create one.")
			return
		}

		config.PrintStoredConnections()
	},
}
