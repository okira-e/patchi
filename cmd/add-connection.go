package cmd

import (
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/spf13/cobra"
)

var AddConnectionCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new connection.",
	Long:  `Add a new database connection to the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		errOpt := config.AddDbConnection()
		if errOpt.IsSome() {
			utils.Abort(errOpt.Unwrap().Error())
		}

		utils.PrintInColor(colors.Green, "Connection added successfully.", false)
	},
}
