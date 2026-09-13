package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	dockercollector "github.com/it-odyssey/waketrail/internal/collectors/docker"
)

func dockerSnapshotPath() (string, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")

	if stateHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		stateHome = filepath.Join(
			home,
			".local",
			"state",
		)
	}

	dir := filepath.Join(
		stateHome,
		"waketrail",
	)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(
		dir,
		"docker-snapshot.json",
	), nil
}

func SaveDockerSnapshot(
	containers []dockercollector.ContainerState,
) error {
	path, err := dockerSnapshotPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		containers,
		"",
		"  ",
	)
	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0644,
	)
}

func LoadDockerSnapshot() (
	[]dockercollector.ContainerState,
	error,
) {
	path, err := dockerSnapshotPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var containers []dockercollector.ContainerState

	if err := json.Unmarshal(
		data,
		&containers,
	); err != nil {
		return nil, err
	}

	return containers, nil
}
