package cmd

import (
	"fmt"
	"time"

	dockercollector "github.com/it-odyssey/waketrail/internal/collectors/docker"
	"github.com/it-odyssey/waketrail/internal/state"
	"github.com/it-odyssey/waketrail/internal/storage"
	"github.com/spf13/cobra"
)

var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observe external system state",
}

var observeDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Observe Docker container state",

	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := dockercollector.Detect()
		if err != nil {
			return err
		}

		previous, err := storage.LoadDockerSnapshot()
		if err != nil {
			return err
		}

		if previous == nil {
			if err := storage.SaveDockerSnapshot(
				current,
			); err != nil {
				return err
			}

			fmt.Printf(
				"Stored Docker baseline: %d containers\n",
				len(current),
			)

			return nil
		}

		transitions := dockercollector.Compare(
			previous,
			current,
		)

		active, err := state.HasActiveSession()
		if err != nil {
			return err
		}

		if active {
			session, err := state.LoadSession()
			if err != nil {
				return err
			}

			store, err := storage.Open()
			if err != nil {
				return err
			}
			defer store.Close()

			for _, transition := range transitions {
				event := storage.TimelineEvent{
					SessionID:  &session.ID,
					EventType:  transition.EventType,
					Source:     "docker",
					Summary:    transition.Summary,
					OccurredAt: time.Now(),
				}

				if _, err := store.InsertTimelineEvent(
					event,
				); err != nil {
					return err
				}
			}
		}

		if err := storage.SaveDockerSnapshot(
			current,
		); err != nil {
			return err
		}

		fmt.Printf(
			"Docker observation complete: %d transitions\n",
			len(transitions),
		)

		return nil
	},
}

func init() {
	observeCmd.AddCommand(
		observeDockerCmd,
	)

	rootCmd.AddCommand(
		observeCmd,
	)
}
