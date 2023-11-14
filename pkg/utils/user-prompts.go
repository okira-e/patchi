package utils

import (
	"sort"

	"github.com/Okira-E/patchi/pkg/safego"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/vars/colors"
	"github.com/manifoldco/promptui"
)

// PromptForDbConnections Prompts the user to choose which two databases to compare against.
func PromptForDbConnections(userConfig types.UserConfig) (string, string, safego.Option[string]) {
	PrintInColor(colors.Cyan, "Choose the connections from the list below to compare:")

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

	PrintInColor(colors.Cyan, "Choose the second database (The dialect must be the same as the first one):")

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
