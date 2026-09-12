# WakeTrail

**A black box flight recorder for DevOps work.**

WakeTrail records the sequence of commands, state changes, and failures that occur during an engineering session so you can reconstruct exactly what happened when troubleshooting infrastructure or systems.

The goal is simple:

> When something breaks, you should not have to rely on memory to figure out what changed.

## Why WakeTrail?

Traditional shell history tells you what commands were entered, but not the full operational context around them.

WakeTrail is designed to capture a richer timeline, including:

- command executed
- timestamp
- working directory
- exit code
- execution duration
- Git repository state
- branch and commit
- changed files
- service and container health
- incident markers
- environment changes

Future adapters are planned for tools such as:

- Docker / Docker Compose
- systemd
- Kubernetes
- Terraform / OpenTofu
- Ansible
- AWS CLI
- Azure CLI
- GitHub CLI

## Current Status

WakeTrail is in early development.

The current CLI supports:

```bash
waketrail start <session-name>
waketrail status
waketrail stop

The Bash integration can currently detect user commands and capture exit codes.

Persistent command event storage and the Omarchy plugin are under active development.

Project Goals

WakeTrail is being built around a few core principles:

Local first — no cloud account required
Low friction — work normally while recording
Forensic visibility — reconstruct what happened after a failure
Modular architecture — collectors and integrations remain independent
Open source — designed for real-world DevOps and Linux workflows
Omarchy native — includes an Omarchy plugin for session control and timeline visibility
Planned Architecture
WakeTrail
│
├── CLI / recorder core
├── shell integrations
├── local event storage
├── collectors
│   ├── Git
│   ├── Docker
│   ├── systemd
│   ├── Kubernetes
│   └── Terraform
│
└── Omarchy plugin
    ├── recording indicator
    ├── session controls
    ├── recent events
    └── incident timeline
Development

WakeTrail is currently developed in Go.

Run locally with:

go run .

Run all tests with:

go test ./...
Roadmap

Initial milestones:

- [x] CLI foundation
- [x] Session start / status / stop
- [x] Bash command capture prototype
- [ ] Persistent command events
- [ ] Execution duration tracking
- [ ] Git context capture
- [ ] SQLite event storage
- [ ] Session timeline output
- [ ] Incident markers
- [ ] Omarchy plugin
- [ ] Docker / systemd collectors
- [ ] Terraform / Kubernetes adapters

Project

WakeTrail is an IT Odyssey project.

Seek Always A New Horizon.

License

License information will be added before the first tagged release.