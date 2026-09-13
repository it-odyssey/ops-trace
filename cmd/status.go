package cmd

import (
	"fmt"

	"github.com/it-odyssey/waketrail/internal/state"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current Waketrail recording session",
	RunE: func(cmd *cobra.Command, args []string) error {
		active, err := state.HasActiveSession()
		if err != nil {
			return err
		}

		if !active {
			fmt.Println("No active recording.")
			return nil
		}

		session, err := state.LoadSession()
		if err != nil {
			return err
		}

		fmt.Printf("Recording: %s\n", session.Name)
		fmt.Printf("Started:   %s\n", session.StartedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
