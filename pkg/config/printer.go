package config

import (
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"os"
)

// PrintStoredConnections prints all the stored connections in the config file.
func PrintStoredConnections() {
	userConfig, errOpt := GetUserConfig()
	if errOpt.IsSome() {
		log.Panicf("Error getting user config: %s", errOpt.Unwrap())
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Connection Name", "Dialect", "Host", "Port", "User", "Password", "Database Name"})

	if len(userConfig.DbConnections) != 0 {
		for connectionName, connection := range userConfig.DbConnections {
			maskedPassword := utils.MaskString(connection.Password)

			t.AppendRow([]any{connectionName, connection.Dialect, connection.Host, connection.Port, connection.User, maskedPassword, connection.Database})
		}
	} else {
		t.AppendRow([]any{"No connections stored"})
	}

	t.Render()
}
