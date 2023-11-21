package config

import (
	"errors"
	"github.com/Okira-E/patchi/pkg/utils"
	"github.com/Okira-E/patchi/safego"
	"os"
	"runtime"

	"github.com/Okira-E/patchi/pkg/types"
)

// getCOnfigFilePathBasedOnOS returns the config file path based on the OS.
func getConfigFilePathBasedOnOS() (string, safego.Option[error]) {
	var osUserName string

	if runtime.GOOS == "windows" {
		osUserName = os.Getenv("USERNAME")
		return "C:\\Users\\" + osUserName + "\\AppData\\Roaming\\patchi\\config.json", safego.None[error]()
	} else if runtime.GOOS == "darwin" {
		osUserName = os.Getenv("USER")
		return "/Users/" + osUserName + "/Library/Application Support/patchi/config.json", safego.None[error]()
	} else if runtime.GOOS == "linux" {
		osHomeDir := os.Getenv("HOME")
		return osHomeDir + "/.config/patchi/config.json", safego.None[error]()
	} else {
		err := errors.New("unsupported OS")
		return "", safego.Some[error](err)
	}
}

// doesConfigFileExists checks if the config file exists.
func doesConfigFileExists() (bool, safego.Option[error]) {
	filePath, errOpt := getConfigFilePathBasedOnOS()
	if errOpt.IsSome() {
		return false, errOpt
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, safego.None[error]()
	}

	return true, safego.None[error]()
}

// createConfigFile creates the config file.
func createConfigFile() safego.Option[error] {
	filePath, errOption := getConfigFilePathBasedOnOS()
	if errOption.IsSome() {
		return errOption
	}

	// Create the directory.
	dirPath := filePath[:len(filePath)-len("/config.json")]
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return safego.Some[error](err)
	}

	// Create the file inside the directory.
	file, err := os.Create(filePath)
	defer file.Close()

	_, err = file.Write([]byte(`{}`))
	if err != nil {
		return safego.Some(err)
	}

	return safego.None[error]()
}

// SetupUserConfig sets up the user config. If the config file doesn't exist, it will create it.
func SetupUserConfig() safego.Option[error] {
	found, errOpt := doesConfigFileExists()
	if errOpt.IsSome() {
		return errOpt
	}

	if !found {
		errOption := createConfigFile()
		if errOption.IsSome() {
			return errOption
		}
	}

	return safego.None[error]()
}

// GetUserConfig gets the user config.
func GetUserConfig() (types.UserConfig, safego.Option[error]) {
	var userConfig types.UserConfig

	filePath, errOption := getConfigFilePathBasedOnOS()
	if errOption.IsSome() {
		return types.UserConfig{}, errOption
	}

	errOpt := utils.ReadJSONFile(filePath, &userConfig)
	if errOpt.IsSome() {
		return types.UserConfig{}, errOpt
	}

	return userConfig, safego.None[error]()
}

func WriteUserConfig(userConfig types.UserConfig) safego.Option[error] {
	filePath, errOption := getConfigFilePathBasedOnOS()
	if errOption.IsSome() {
		return errOption
	}

	errOpt := utils.WriteToJSONFile(filePath, userConfig)
	if errOpt.IsSome() {
		return errOpt
	}

	return safego.None[error]()
}
