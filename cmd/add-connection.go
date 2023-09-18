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
		_, errOpt := config.GetUserConfig()
		if errOpt.IsSome() {
			logger.PrintInColor(colors.Red, "Config file not found. Run `patchi init` to create one.")
			return
		}

		errOpt = config.AddConnection()
		if errOpt.IsSome() {
			logger.PrintInColor(colors.Red, errOpt.Unwrap().Error())
		}

		logger.PrintInColor(colors.Green, "Connection added successfully.")
	},
}
