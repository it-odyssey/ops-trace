package storage

import (
	"testing"

	dockercollector "github.com/it-odyssey/waketrail/internal/collectors/docker"
)

func TestDockerSnapshotSaveAndLoad(t *testing.T) {
	t.Setenv(
		"XDG_STATE_HOME",
		t.TempDir(),
	)

	expected := []dockercollector.ContainerState{
		{
			Name:   "traefik",
			State:  "running",
			Status: "Up 5 minutes (healthy)",
			Health: "healthy",
		},
		{
			Name:   "postgres",
			State:  "running",
			Status: "Up 5 minutes",
		},
	}

	if err := SaveDockerSnapshot(expected); err != nil {
		t.Fatalf(
			"SaveDockerSnapshot() returned error: %v",
			err,
		)
	}

	actual, err := LoadDockerSnapshot()
	if err != nil {
		t.Fatalf(
			"LoadDockerSnapshot() returned error: %v",
			err,
		)
	}

	if len(actual) != len(expected) {
		t.Fatalf(
			"len(actual) = %d, want %d",
			len(actual),
			len(expected),
		)
	}

	if actual[0].Name != "traefik" {
		t.Errorf(
			"actual[0].Name = %q, want %q",
			actual[0].Name,
			"traefik",
		)
	}

	if actual[0].Health != "healthy" {
		t.Errorf(
			"actual[0].Health = %q, want %q",
			actual[0].Health,
			"healthy",
		)
	}
}

func TestLoadDockerSnapshotWhenMissing(t *testing.T) {
	t.Setenv(
		"XDG_STATE_HOME",
		t.TempDir(),
	)

	containers, err := LoadDockerSnapshot()
	if err != nil {
		t.Fatalf(
			"LoadDockerSnapshot() returned error: %v",
			err,
		)
	}

	if containers != nil {
		t.Errorf(
			"containers = %#v, want nil",
			containers,
		)
	}
}
