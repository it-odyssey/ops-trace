package cmd

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/it-odyssey/waketrail/internal/storage"
	"github.com/spf13/cobra"
)

type displayEvent struct {
	OccurredAt time.Time
	Kind       string

	CommandEvent  *storage.CommandEvent
	TimelineEvent *storage.TimelineEvent
}

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

		commandEvents, err := store.CommandEventsForSession(session.ID)
		if err != nil {
			return err
		}

		timelineEvents, err := store.TimelineEventsForSession(session.ID)
		if err != nil {
			return err
		}

		events := make([]displayEvent, 0, len(commandEvents)+len(timelineEvents))

		for i := range commandEvents {
			event := &commandEvents[i]

			events = append(events, displayEvent{
				OccurredAt:   event.StartedAt,
				Kind:         "command",
				CommandEvent: event,
			})
		}

		for i := range timelineEvents {
			event := &timelineEvents[i]

			events = append(events, displayEvent{
				OccurredAt:    event.OccurredAt,
				Kind:          "timeline",
				TimelineEvent: event,
			})
		}

		sort.Slice(events, func(i, j int) bool {
			return events[i].OccurredAt.Before(events[j].OccurredAt)
		})

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

		fmt.Printf("Events: %d\n", len(events))
		fmt.Println()

		if len(events) == 0 {
			fmt.Println("No events recorded for this session.")
			return nil
		}

		for _, event := range events {
			switch event.Kind {
			case "command":
				if err := printCommandEvent(store, *event.CommandEvent); err != nil {
					return err
				}

			case "timeline":
				printTimelineEvent(*event.TimelineEvent)
			}
		}

		return nil
	},
}

func printCommandEvent(store *storage.Store, event storage.CommandEvent) error {
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

	gitContext, err := store.GitContextForCommandEvent(event.ID)

	if err != nil && !errors.Is(err, storage.ErrGitContextNotFound) {
		return err
	}

	if err == nil {
		commit := gitContext.CommitSHA

		if len(commit) > 7 {
			commit = commit[:7]
		}

		dirty := "clean"
		if gitContext.Dirty {
			dirty = "dirty"
		}

		fmt.Println()
		fmt.Println("          git:")
		fmt.Printf(
			"            repo:   %s\n",
			gitContext.RepositoryRoot,
		)
		fmt.Printf(
			"            branch: %s\n",
			gitContext.Branch,
		)
		fmt.Printf(
			"            commit: %s\n",
			commit,
		)
		fmt.Printf(
			"            state:  %s\n",
			dirty,
		)
	}

	fmt.Println()

	return nil
}

func printTimelineEvent(event storage.TimelineEvent) {
	switch event.EventType {
	case "note":
		fmt.Printf(
			"%s  📝 NOTE\n",
			event.OccurredAt.Format("15:04:05"),
		)

		fmt.Printf(
			"          %s\n",
			event.Summary,
		)

	default:
		fmt.Printf(
			"%s  • %s\n",
			event.OccurredAt.Format("15:04:05"),
			event.EventType,
		)

		fmt.Printf(
			"          source: %s\n",
			event.Source,
		)

		fmt.Printf(
			"          %s\n",
			event.Summary,
		)
	}

	fmt.Println()
}

func formatTimelineTime(value time.Time) string {
	return value.Format("2006-01-02 15:04:05")
}

func init() {
	rootCmd.AddCommand(showCmd)
}
