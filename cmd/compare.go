package cmd

import (
	"fmt"
	"log"

	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/datasource"
	"github.com/Okira-E/patchi/pkg/tui"
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

		firstDbConnectionName, secondDbConnectionName, errMsgOpt := utils.PromptForDbConnections(userConfig)
		if errMsgOpt.IsSome() {
			utils.PrintInColor(colors.Red, errMsgOpt.Unwrap())

			return
		}

		firstDbConnection, errOpt := datasource.GetDataSource(userConfig.DbConnections[firstDbConnectionName])
		if errOpt.IsSome() {
			log.Fatalf("Error connecting to %s: %s", firstDbConnectionName, errOpt.Unwrap())
		}
		defer func() {
			err := firstDbConnection.Close()
			if err != nil {
				log.Fatalf("Error closing connection to %s: %s", firstDbConnectionName, err)
			}
		}()

		secondDbConnection, errOpt := datasource.GetDataSource(userConfig.DbConnections[secondDbConnectionName])
		if errOpt.IsSome() {
			log.Fatalf("Error connecting to %s: %s", secondDbConnectionName, errOpt.Unwrap())
		}
		defer func() {
			err := secondDbConnection.Close()
			if err != nil {
				log.Fatalf("Error closing connection to %s: %s", secondDbConnectionName, err)
			}
		}()

		if err := firstDbConnection.Ping(); err != nil {
			utils.Abort(fmt.Sprintf("Failed to ping the \"%s\" database", userConfig.DbConnections[firstDbConnectionName].Name))
		}

		if err := secondDbConnection.Ping(); err != nil {
			utils.Abort(fmt.Sprintf("Failed to ping the \"%s\" database", userConfig.DbConnections[secondDbConnectionName].Name))
		}

		params := &tui.HomeParams{
			FirstDb: tui.DbConnectionInfo{
				Info:          userConfig.DbConnections[firstDbConnectionName],
				SqlConnection: firstDbConnection,
			},
			SecondDb: tui.DbConnectionInfo{
				Info:          userConfig.DbConnections[secondDbConnectionName],
				SqlConnection: secondDbConnection,
			},
		}

		tui.GlobalRenderer(params)
	},
}
