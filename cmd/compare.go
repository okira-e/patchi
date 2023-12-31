package cmd

import (
	"fmt"

	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/prompts"
	"github.com/Okira-E/patchi/pkg/tui"
	"github.com/Okira-E/patchi/pkg/tui/patchi_renderer"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
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
			utils.Abort(errOpt.Unwrap().Error())
		}

		firstDbConnectionInfo, secondDbConnectionInfo, errMsgOpt := prompts.PromptForDbConnections(userConfig)
		if errMsgOpt.IsSome() {
			utils.Abort(errOpt.Unwrap().Error())
		}

		firstDbConnection, errOpt := firstDbConnectionInfo.Connect()
		if errOpt.IsSome() {
			utils.Abort(fmt.Sprintf("Error connecting to %s: %s", firstDbConnectionInfo.DatabaseName, errOpt.Unwrap()))
		}
		defer func() {
			err := firstDbConnection.Close()
			if err != nil {
				utils.Abort(fmt.Sprintf("Error closing connection to %s: %s", firstDbConnectionInfo.DatabaseName, err))
			}
		}()

		secondDbConnection, errOpt := secondDbConnectionInfo.Connect()
		if errOpt.IsSome() {
			utils.Abort(fmt.Sprintf("Error connecting to %s: %s", secondDbConnectionInfo.DatabaseName, errOpt.Unwrap()))
		}
		defer func() {
			err := secondDbConnection.Close()
			if err != nil {
				utils.Abort(fmt.Sprintf("Error closing connection to %s: %s", secondDbConnectionInfo.DatabaseName, err))
			}
		}()

		if err := firstDbConnection.Ping(); err != nil {
			utils.Abort(fmt.Sprintf("Failed to ping the \"%s\" database", firstDbConnectionInfo.Name))
		}

		if err := secondDbConnection.Ping(); err != nil {
			utils.Abort(fmt.Sprintf("Failed to ping the \"%s\" database", secondDbConnectionInfo.Name))
		}

		params := &patchi_renderer.PatchiRendererParams {
			FirstDb: types.DbConnection{
				Info:          firstDbConnectionInfo,
				SqlConnection: firstDbConnection,
			},
			SecondDb: types.DbConnection{
				Info:          secondDbConnectionInfo,
				SqlConnection: secondDbConnection,
			},
		}

		tui.RenderTui(params)
	},
}
