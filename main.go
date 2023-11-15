package main

import (
	"fmt"
	"github.com/Okira-E/patchi/cmd"
	"github.com/Okira-E/patchi/pkg/config"
	"github.com/Okira-E/patchi/pkg/utils"
)

func main() {
	// -- User config setup
	errOpt := config.SetupUserConfig()
	if errOpt.IsSome() {
		utils.Abort(fmt.Sprintf("Error setting up user config: %s", errOpt.Unwrap()))
	}

	cmd.Execute()
}
