package docker

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type ContainerState struct {
	Name   string
	State  string
	Status string
	Health string
}

func Detect() ([]ContainerState, error) {
	cmd := exec.Command(
		"docker",
		"ps",
		"-a",
		"--format",
		"{{.Names}}\t{{.State}}\t{{.Status}}",
	)

	output, err := cmd.Output()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, errors.New("docker CLI not found")
		}

		return nil, err
	}

	return parseDockerPS(output)
}

func parseDockerPS(output []byte) ([]ContainerState, error) {
	var containers []ContainerState

	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		fields := strings.SplitN(line, "\t", 3)

		if len(fields) != 3 {
			continue
		}

		status := fields[2]

		containers = append(containers, ContainerState{
			Name:   fields[0],
			State:  fields[1],
			Status: status,
			Health: healthFromStatus(status),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return containers, nil
}

func healthFromStatus(status string) string {
	switch {
	case strings.Contains(status, "(healthy)"):
		return "healthy"

	case strings.Contains(status, "(unhealthy)"):
		return "unhealthy"

	case strings.Contains(status, "(health: starting)"):
		return "starting"

	default:
		return ""
	}
}
