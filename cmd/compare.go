package cmd

import (
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/datasource"
	"github.com/Okira-E/patchi/pkg/safego"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils/logger"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"sort"
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
			logger.PrintInColor(colors.Red, errOpt.Unwrap().Error())

			return
		}

		firstDbConnectionInfo, secondDbConnectionInfo, errMsgOpt := promptForDbConnections(userConfig)
		if errMsgOpt.IsSome() {
			logger.PrintInColor(colors.Red, errMsgOpt.Unwrap())

			return
		}

		firstDbConnection, errOpt := datasource.GetDataSource(userConfig.DbConnections[firstDbConnectionInfo])
		if errOpt.IsSome() {
			log.Panicf("Error connecting to %s: %s", firstDbConnectionInfo, errOpt.Unwrap())
		}
		defer firstDbConnection.Close()

		secondDbConnection, errOpt := datasource.GetDataSource(userConfig.DbConnections[secondDbConnectionInfo])
		if errOpt.IsSome() {
			log.Panicf("Error connecting to %s: %s", secondDbConnectionInfo, errOpt.Unwrap())
		}
		defer secondDbConnection.Close()

		props := &types.CompareRootRendererProps{
			FirstDb: types.DbConnection{
				Info:          userConfig.DbConnections[firstDbConnectionInfo],
				SqlConnection: firstDbConnection,
			},
			SecondDb: types.DbConnection{
				Info:          userConfig.DbConnections[secondDbConnectionInfo],
				SqlConnection: secondDbConnection,
			},
		}

		_ = props
	},
}

func promptForDbConnections(userConfig types.UserConfig) (string, string, safego.Option[string]) {
	logger.PrintInColor(colors.Cyan, "Choose the connections from the list below to compare:")

	allConnectionNames := []string{}
	for connectionName, _ := range userConfig.DbConnections {
		allConnectionNames = append(allConnectionNames, connectionName)
	}
	// Sort allConnectionNames
	sort.Strings(allConnectionNames)

	firstConnectionToChoosePrmpt := promptui.Select{
		Label: "Choose connections",
		Items: allConnectionNames,
		Size:  10,
	}
	_, firstSelectedConnection, err := firstConnectionToChoosePrmpt.Run()
	if err != nil {
		return "", "", safego.Some(err.Error())
	}

	logger.PrintInColor(colors.Cyan, "Choose the second database (The dialect must be the same as the first one):")

	// Filter out the first selected connection from the list of connections for the second prompt.
	// Also, filter out the connections that have a different dialect than the first selected connection.

	filteredConnectionNamesForSecondPrompt := []string{}
	for connectionName, connection := range userConfig.DbConnections {
		if connectionName != firstSelectedConnection && connection.Dialect == userConfig.DbConnections[firstSelectedConnection].Dialect {
			filteredConnectionNamesForSecondPrompt = append(filteredConnectionNamesForSecondPrompt, connectionName)
		}
	}

	if len(filteredConnectionNamesForSecondPrompt) == 0 {
		return "", "", safego.Some("no connections with the same dialect remaining")
	}

	secondSelectedConnectionPrmpt := promptui.Select{
		Label: "Choose connections",
		Items: filteredConnectionNamesForSecondPrompt,
		Size:  10,
	}
	_, secondSelectedConnection, err := secondSelectedConnectionPrmpt.Run()
	if err != nil {
		return "", "", safego.Some(err.Error())
	}

	return firstSelectedConnection, secondSelectedConnection, safego.None[string]()
}
