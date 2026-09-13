package git

import (
	"errors"
	"os/exec"
	"strings"
)

type Context struct {
	IsRepository bool
	Root         string
	Branch       string
	Commit       string
	Dirty        bool
}

func Detect(workingDirectory string) (Context, error) {
	root, err := runGit(
		workingDirectory,
		"rev-parse",
		"--show-toplevel",
	)

	if err != nil {
		if isNotRepository(err) {
			return Context{
				IsRepository: false,
			}, nil
		}

		return Context{}, err
	}

	branch, err := runGit(
		workingDirectory,
		"branch",
		"--show-current",
	)
	if err != nil {
		return Context{}, err
	}

	commit, err := runGit(
		workingDirectory,
		"rev-parse",
		"HEAD",
	)
	if err != nil {
		return Context{}, err
	}

	status, err := runGit(
		workingDirectory,
		"status",
		"--porcelain",
	)
	if err != nil {
		return Context{}, err
	}

	return Context{
		IsRepository: true,
		Root:         root,
		Branch:       branch,
		Commit:       commit,
		Dirty:        status != "",
	}, nil
}

func runGit(workingDirectory string, args ...string) (string, error) {
	commandArgs := append(
		[]string{"-C", workingDirectory},
		args...,
	)

	cmd := exec.Command("git", commandArgs...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func isNotRepository(err error) bool {
	var exitErr *exec.ExitError

	if !errors.As(err, &exitErr) {
		return false
	}

	return exitErr.ExitCode() == 128
}
