package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "patchi",
	Short: "This is a tool for migrating database environments.",
	Long: `
Patchi connects to 2 of your databases and shows you the differences between them. Useful for 
migrating database environments.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Panicf("Error printing help: %s", err.Error())
		}
	},
}

func Execute() {
	rootCmd.AddCommand(ListConnectionsCmd)
	rootCmd.AddCommand(AddConnectionCmd)
	rootCmd.AddCommand(RmConnectionCmd)
	rootCmd.AddCommand(StartCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Panicf("Error executing root command: %s", err.Error())
	}
}
