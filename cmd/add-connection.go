package cmd

import (
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/utils/logger"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/spf13/cobra"
)

var AddConnectionCmd = &cobra.Command{
	Use:   "add-connection",
	Short: "Add a new connection.",
	Long:  `Add a new database connection to the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		errOpt := config.AddConnection()
		if errOpt.IsSome() {
			logger.PrintInColor(colors.Red, errOpt.Unwrap().Error())
			return
		}

		logger.PrintInColor(colors.Green, "Connection added successfully.")
	},
}
