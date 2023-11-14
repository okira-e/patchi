package config

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/safego"
	"github.com/Okira-E/patchi/pkg/types"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/pkg/vars"
	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"strconv"
)

// AddDbConnection adds a new connection to the config file.
// It returns an error if the connection already exists.
func AddDbConnection() safego.Option[error] {
	userConfig, errOpt := GetUserConfig()
	if errOpt.IsSome() {
		return errOpt
	}

	// -- Take the connection details from the user.

	// Dialect
	dialectPrmpt := promptui.Select{
		Label: "Dialect",
		Items: vars.SupportedDatabases,
	}
	_, dialect, err := dialectPrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}

	// Name
	NamePrmpt := promptui.Prompt{
		Label: "Connection name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("name cannot be empty")
			}

			return nil
		},
	}
	connectionName, err := NamePrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}

	if _, ok := userConfig.DbConnections[connectionName]; ok {
		return safego.Some[error](fmt.Errorf("connection with name %s already exists", connectionName))
	}

	// Host
	HostPrmpt := promptui.Prompt{
		Label: "Host",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("host cannot be empty")
			}

			return nil
		},
	}
	host, err := HostPrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}

	// Port
	PortPrmpt := promptui.Prompt{
		Label: "Port",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("port cannot be empty")
			}
			// Check if the port is a number.
			intPort, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("port must be a number")
			}

			if intPort < 0 || intPort > 65535 {
				return fmt.Errorf("port must be between 0 and 65535")
			}

			return nil
		},
	}
	stringPort, err := PortPrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}
	port, err := strconv.Atoi(stringPort)
	if err != nil {
		return safego.Some[error](err)
	}

	// User
	userPrmpt := promptui.Prompt{
		Label: "User",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("user cannot be empty")
			}

			return nil
		},
	}
	user, err := userPrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}

	// Password
	passwordPrmpt := promptui.Prompt{
		Label: "Password",
		Mask:  '*',
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("password cannot be empty")
			}

			return nil
		},
	}
	password, err := passwordPrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}

	// Database name
	databasePrmpt := promptui.Prompt{
		Label: "Database name",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("database name cannot be empty")
			}

			return nil
		},
	}
	database, err := databasePrmpt.Run()
	if err != nil {
		return safego.Some[error](err)
	}

	// -- Add the connection to the config file.

	if len(userConfig.DbConnections) == 0 {
		userConfig.DbConnections = make(map[string]*types.DbConnectionInfo)
	}

	userConfig.DbConnections[connectionName] = &types.DbConnectionInfo{
		Dialect:  dialect,
		Name:     connectionName,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}

	configFilePathBasedOnOs, errOpt := getConfigFilePathBasedOnOS()
	if errOpt.IsSome() {
		return errOpt
	}

	utils.WriteToJSONFile(configFilePathBasedOnOs, userConfig)

	return safego.None[error]()
}

// RmConnection removes a connection from the config file.
// It returns an error if the connection does not exist.
func RmConnection(connectionName string) safego.Option[error] {
	userConfig, errOpt := GetUserConfig()
	if errOpt.IsSome() {
		return errOpt
	}

	if len(userConfig.DbConnections) == 0 {
		return safego.Some[error](fmt.Errorf("no connections to remove"))
	}

	// -- Remove the connection from the config file.
	configFilePathBasedOnOs, errOpt := getConfigFilePathBasedOnOS()
	if errOpt.IsSome() {
		return errOpt
	}

	delete(userConfig.DbConnections, connectionName)
	utils.WriteToJSONFile(configFilePathBasedOnOs, userConfig)

	return safego.None[error]()
}

// PrintStoredConnections prints all the stored connections in the config file.
func PrintStoredConnections() {
	userConfig, errOpt := GetUserConfig()
	if errOpt.IsSome() {
		log.Fatalf("Error getting user config: %s", errOpt.Unwrap())
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
