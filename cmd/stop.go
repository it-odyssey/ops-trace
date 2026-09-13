package cmd

import (
	"fmt"
	"time"

	"github.com/it-odyssey/waketrail/internal/state"
	"github.com/it-odyssey/waketrail/internal/storage"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current WakeTrail recording session",
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

		store, err := storage.Open()
		if err != nil {
			return err
		}
		defer store.Close()

		endedAt := time.Now()

		if err := store.EndSession(session.ID, endedAt); err != nil {
			return err
		}

		duration := endedAt.Sub(session.StartedAt)

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
