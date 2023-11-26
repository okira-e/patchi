package prompts

import (
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
	"sort"

	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/manifoldco/promptui"
)

// PromptForDbConnections Prompts the user to choose which two databases to compare against.
func PromptForDbConnections(userConfig types.UserConfig) (*types.DbConnectionInfo, *types.DbConnectionInfo, safego.Option[string]) {
	if len(userConfig.DbConnections) == 0 {
		return &types.DbConnectionInfo{}, &types.DbConnectionInfo{}, safego.Some("no connections found")
	} else if len(userConfig.DbConnections) == 1 {
		return &types.DbConnectionInfo{}, &types.DbConnectionInfo{}, safego.Some("only one connection found. Please add another connection to compare against")
	}

	utils.PrintInColor(colors.Cyan, "Choose the connections from the list below to compare:")

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
	_, firstSelectedConnectionName, err := firstConnectionToChoosePrmpt.Run()
	if err != nil {
		return &types.DbConnectionInfo{}, &types.DbConnectionInfo{}, safego.Some(err.Error())
	}

	utils.PrintInColor(colors.Cyan, "Choose the second database (The dialect must be the same as the first one):")

	// Filter out the first selected connection from the list of connections for the second
	// prompt (We don't want to compare a database with itself.)
	// Also, filter out the connections that have a different dialect than the first selected connection.

	filteredConnectionNamesForSecondPrompt := []string{}
	for connectionName, connection := range userConfig.DbConnections {
		if connectionName != firstSelectedConnectionName && connection.Dialect == userConfig.DbConnections[firstSelectedConnectionName].Dialect {
			filteredConnectionNamesForSecondPrompt = append(filteredConnectionNamesForSecondPrompt, connectionName)
		}
	}

	if len(filteredConnectionNamesForSecondPrompt) == 0 {
		return &types.DbConnectionInfo{}, &types.DbConnectionInfo{}, safego.Some("no connections with the same dialect remaining")
	}

	secondSelectedConnectionPrmpt := promptui.Select{
		Label: "Choose connections",
		Items: filteredConnectionNamesForSecondPrompt,
		Size:  10,
	}
	_, secondSelectedConnectionName, err := secondSelectedConnectionPrmpt.Run()
	if err != nil {
		return &types.DbConnectionInfo{}, &types.DbConnectionInfo{}, safego.Some(err.Error())
	}

	return userConfig.DbConnections[firstSelectedConnectionName], userConfig.DbConnections[secondSelectedConnectionName], safego.None[string]()
}
