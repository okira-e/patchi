package cmd

import (
	"fmt"
	"log"

	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/tui"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "compare",
	Short: "CompareRoot 2 databases.",
	Long: `
Patchi connects to 2 of your databases and shows you the differences between them. Useful for
migrating database environments.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		userConfig, errOpt := config.GetUserConfig()
		if errOpt.IsSome() {
			utils.PrintInColor(colors.Red, errOpt.Unwrap().Error())

			return
		}

		firstDbConnectionInfo, secondDbConnectionInfo, errMsgOpt := utils.PromptForDbConnections(userConfig)
		if errMsgOpt.IsSome() {
			utils.PrintInColor(colors.Red, errMsgOpt.Unwrap())

			return
		}

		firstDbConnection, errOpt := firstDbConnectionInfo.Connect()
		if errOpt.IsSome() {
			log.Fatalf("Error connecting to %s: %s", firstDbConnectionInfo, errOpt.Unwrap())
		}
		defer func() {
			err := firstDbConnection.Close()
			if err != nil {
				log.Fatalf("Error closing connection to %s: %s", firstDbConnectionInfo, err)
			}
		}()

		secondDbConnection, errOpt := secondDbConnectionInfo.Connect()
		if errOpt.IsSome() {
			log.Fatalf("Error connecting to %s: %s", secondDbConnectionInfo, errOpt.Unwrap())
		}
		defer func() {
			err := secondDbConnection.Close()
			if err != nil {
				log.Fatalf("Error closing connection to %s: %s", secondDbConnectionInfo, err)
			}
		}()

		if err := firstDbConnection.Ping(); err != nil {
			utils.Abort(fmt.Sprintf("Failed to ping the \"%s\" database", firstDbConnectionInfo.Name))
		}

		if err := secondDbConnection.Ping(); err != nil {
			utils.Abort(fmt.Sprintf("Failed to ping the \"%s\" database", secondDbConnectionInfo.Name))
		}

		params := &tui.HomeParams{
			FirstDb: types.DbConnection{
				Info:          firstDbConnectionInfo,
				SqlConnection: firstDbConnection,
			},
			SecondDb: types.DbConnection{
				Info:          secondDbConnectionInfo,
				SqlConnection: secondDbConnection,
			},
		}

		tui.GlobalRenderer(params)
	},
}
