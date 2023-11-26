package cmd

import (
	"fmt"
	"github.com/Okira-E/patchi/pkg/utils"
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
			utils.Abort(fmt.Sprintf("Error printing help: %s", err.Error()))
		}
	},
}

func Execute() {
	RmConnectionCmd.Flags().String("name", "", "Name of the connection to remove.")

	rootCmd.AddCommand(ListConnectionsCmd)
	rootCmd.AddCommand(AddConnectionCmd)
	rootCmd.AddCommand(RmConnectionCmd)
	rootCmd.AddCommand(StartCmd)

	err := rootCmd.Execute()
	if err != nil {
		utils.Abort(fmt.Sprintf("Error executing root command: %s", err.Error()))
	}
}
