package main

import (
	"github.com/Okira-E/patchi/cmd"
	"github.com/Okira-E/patchi/pkg/config"
	"log"
)

func main() {
	// -- User config setup
	errOpt := config.SetupUserConfig()
	if errOpt.IsSome() {
		log.Fatalf("Error setting up user config: %s", errOpt.Unwrap())
	}

	cmd.Execute()
}
