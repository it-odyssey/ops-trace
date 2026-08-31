package cmd

import (
	"fmt"
	"time"

	"github.com/Jeff-Fontenot/flight-recorder/internal/state"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [session-name]",
	Short: "Start a new flight recording session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		active, err := state.HasActiveSession()
		if err != nil {
			return err
		}

		if active {
			session, err := state.LoadSession()
			if err != nil {
				return err
			}

			fmt.Printf("A recording is already active: %s\n", session.Name)
			return nil
		}
		session := state.Session{
			Name:      args[0],
			StartedAt: time.Now(),
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
