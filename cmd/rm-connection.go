package cmd

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/utils/logger"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var RmConnectionCmd = &cobra.Command{
	Use:   "rm-connection",
	Short: "Removes a connection from the config file.",
	Long:  "Removes a stored database connection from the config file.",
	Run: func(cmd *cobra.Command, args []string) {
		_, errOpt := config.GetUserConfig()
		if errOpt.IsSome() {
			logger.PrintInColor(colors.Red, "Config file not found. Run `patchi init` to create one.")
			return
		}

		namePrmpt := promptui.Prompt{
			Label: "Connection name",
			Validate: func(s string) error {
				if s == "" {
					return fmt.Errorf("name cannot be empty")
				}

				return nil
			},
		}
		connectionName, err := namePrmpt.Run()
		if err != nil {
			logger.PrintInColor(colors.Red, err.Error())
			return
		}

		errOpt = config.RmConnection(connectionName)
		if errOpt.IsSome() {
			logger.PrintInColor(colors.Red, errOpt.Unwrap().Error())
			return
		}

		logger.PrintInColor(colors.Green, "Connection "+connectionName+" removed successfully.")
	},
}
