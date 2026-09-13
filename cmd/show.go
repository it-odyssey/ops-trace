package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/it-odyssey/waketrail/internal/storage"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [session-name]",
	Short: "Show a WakeTrail session timeline",
	Args:  cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.Open()
		if err != nil {
			return err
		}
		defer store.Close()

		var session storage.SessionRecord

		if len(args) == 1 {
			session, err = store.SessionByName(args[0])
		} else {
			session, err = store.LatestSession()
		}

		if errors.Is(err, storage.ErrSessionNotFound) {
			fmt.Println("No matching WakeTrail session found.")
			return nil
		}

		if err != nil {
			return err
		}

		events, err := store.CommandEventsForSession(session.ID)
		if err != nil {
			return err
		}

		fmt.Printf("WakeTrail Session: %s\n", session.Name)
		fmt.Printf("Started: %s\n", formatTimelineTime(session.StartedAt))

		if session.EndedAt != nil {
			fmt.Printf("Ended:   %s\n", formatTimelineTime(*session.EndedAt))
			fmt.Printf(
				"Duration: %s\n",
				session.EndedAt.Sub(session.StartedAt).Round(time.Second),
			)
		} else {
			fmt.Println("Status:  Recording")
		}

		fmt.Printf("Commands: %d\n", len(events))
		fmt.Println()

		if len(events) == 0 {
			fmt.Println("No command events recorded for this session.")
			return nil
		}

		for _, event := range events {
			status := "✓"

			if event.ExitCode != 0 {
				status = "✗"
			}

			duration := event.EndedAt.Sub(event.StartedAt)

			fmt.Printf(
				"%s  %s  %s\n",
				event.StartedAt.Format("15:04:05"),
				status,
				event.Command,
			)

			fmt.Printf(
				"          exit: %d  duration: %s\n",
				event.ExitCode,
				duration.Round(time.Millisecond),
			)

			fmt.Printf(
				"          cwd: %s\n",
				event.Cwd,
			)

			fmt.Println()
		}

		return nil
	},
}

func formatTimelineTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}

func init() {
	rootCmd.AddCommand(showCmd)
}
