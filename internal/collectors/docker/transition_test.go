package docker

import "testing"

func TestCompareDetectsFailure(t *testing.T) {
	previous := []ContainerState{
		{
			Name:   "traefik",
			State:  "running",
			Health: "healthy",
		},
	}

	current := []ContainerState{
		{
			Name:   "traefik",
			State:  "running",
			Health: "unhealthy",
		},
	}

	transitions := Compare(previous, current)

	if len(transitions) != 1 {
		t.Fatalf(
			"len(transitions) = %d, want 1",
			len(transitions),
		)
	}

	if transitions[0].EventType != EventFailure {
		t.Errorf(
			"EventType = %q, want %q",
			transitions[0].EventType,
			EventFailure,
		)
	}

	if transitions[0].Summary !=
		"traefik: running/healthy -> running/unhealthy" {
		t.Errorf(
			"Summary = %q",
			transitions[0].Summary,
		)
	}
}

func TestCompareDetectsRecovery(t *testing.T) {
	previous := []ContainerState{
		{
			Name:  "api",
			State: "restarting",
		},
	}

	current := []ContainerState{
		{
			Name:  "api",
			State: "running",
		},
	}

	transitions := Compare(previous, current)

	if len(transitions) != 1 {
		t.Fatalf(
			"len(transitions) = %d, want 1",
			len(transitions),
		)
	}

	if transitions[0].EventType != EventRecovery {
		t.Errorf(
			"EventType = %q, want %q",
			transitions[0].EventType,
			EventRecovery,
		)
	}
}

func TestCompareDetectsStateChange(t *testing.T) {
	previous := []ContainerState{
		{
			Name:  "worker",
			State: "created",
		},
	}

	current := []ContainerState{
		{
			Name:  "worker",
			State: "running",
		},
	}

	transitions := Compare(previous, current)

	if len(transitions) != 1 {
		t.Fatalf(
			"len(transitions) = %d, want 1",
			len(transitions),
		)
	}

	if transitions[0].EventType != EventStateChange {
		t.Errorf(
			"EventType = %q, want %q",
			transitions[0].EventType,
			EventStateChange,
		)
	}
}

func TestCompareIgnoresUnchangedContainer(t *testing.T) {
	previous := []ContainerState{
		{
			Name:   "postgres",
			State:  "running",
			Health: "healthy",
		},
	}

	current := []ContainerState{
		{
			Name:   "postgres",
			State:  "running",
			Health: "healthy",
		},
	}

	transitions := Compare(previous, current)

	if len(transitions) != 0 {
		t.Fatalf(
			"len(transitions) = %d, want 0",
			len(transitions),
		)
	}
}

func TestCompareDetectsNewContainer(t *testing.T) {
	current := []ContainerState{
		{
			Name:  "nginx",
			State: "running",
		},
	}

	transitions := Compare(nil, current)

	if len(transitions) != 1 {
		t.Fatalf(
			"len(transitions) = %d, want 1",
			len(transitions),
		)
	}

	if transitions[0].EventType != EventStateChange {
		t.Errorf(
			"EventType = %q, want %q",
			transitions[0].EventType,
			EventStateChange,
		)
	}

	if transitions[0].Summary != "nginx appeared: running" {
		t.Errorf(
			"Summary = %q",
			transitions[0].Summary,
		)
	}
}

func TestCompareDetectsRemovedContainer(t *testing.T) {
	previous := []ContainerState{
		{
			Name:  "nginx",
			State: "exited",
		},
	}

	transitions := Compare(previous, nil)

	if len(transitions) != 1 {
		t.Fatalf(
			"len(transitions) = %d, want 1",
			len(transitions),
		)
	}

	if transitions[0].Summary != "nginx disappeared" {
		t.Errorf(
			"Summary = %q",
			transitions[0].Summary,
		)
	}
}
