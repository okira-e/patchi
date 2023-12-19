package cmd

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var RmConnectionCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a connection from the config file.",
	Long:  "Removes a stored database connection from the config file.",
	Run: func(cmd *cobra.Command, args []string) {
		// FEAT: Add a way to remove all connections.
		// FEAT: Make removing a connection a selection prompt.

		var err error
		var connectionName string

		if len(args) > 0 {
			connectionName = args[0]
		} else {
			namePrmpt := promptui.Prompt{
				Label: "Connection name",
				Validate: func(s string) error {
					if s == "" {
						return fmt.Errorf("name cannot be empty")
					}

					return nil
				},
			}

			connectionName, err = namePrmpt.Run()
			if err != nil {
				utils.Abort(err.Error())
			}
		}

		errOpt := config.RmConnection(connectionName)
		if errOpt.IsSome() {
			utils.Abort(errOpt.Unwrap().Error())
		}

		utils.PrintInColor(colors.Green, "Connection "+connectionName+" removed successfully.", false)
	},
}
