package cmd

import (
	"fmt"
	"time"

	"github.com/it-odyssey/waketrail/internal/state"
	"github.com/it-odyssey/waketrail/internal/storage"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [session-name]",
	Short: "Start a new flight recording session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.Open()
		if err != nil {
			return err
		}
		defer store.Close()

		startedAt := time.Now()

		sessionID, err := store.CreateSession(args[0], startedAt)
		if err != nil {
			return err
		}

		session := state.Session{
			ID:        sessionID,
			Name:      args[0],
			StartedAt: startedAt,
		}

		if err := state.SaveSession(session); err != nil {
			return err
		}

		fmt.Printf("Recording session: %s\n", session.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
