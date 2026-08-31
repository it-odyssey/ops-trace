package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:    "record [command]",
	Short:  "Record a shell command",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Captured: %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
