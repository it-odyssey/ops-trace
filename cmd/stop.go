package cmd

import (
	"fmt"
	"time"

	"github.com/it-odyssey/waketrail/internal/state"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current flight recording session",
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

		duration := time.Since(session.StartedAt)

		if err := state.ClearSession(); err != nil {
			return err
		}

		fmt.Printf("Stopped recording: %s\n", session.Name)
		fmt.Printf("Duration: %s\n", duration.Round(time.Second))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
