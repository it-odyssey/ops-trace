package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDetectGitRepository(t *testing.T) {
	repoDir := createTestRepository(t)

	context, err := Detect(repoDir)
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	if !context.IsRepository {
		t.Fatal("IsRepository = false, want true")
	}

	if context.Root != repoDir {
		t.Errorf(
			"Root = %q, want %q",
			context.Root,
			repoDir,
		)
	}

	if context.Branch == "" {
		t.Error("Branch is empty")
	}

	if context.Commit == "" {
		t.Error("Commit is empty")
	}

	if context.Dirty {
		t.Error("Dirty = true, want false")
	}
}

func TestDetectDirtyRepository(t *testing.T) {
	repoDir := createTestRepository(t)

	filePath := filepath.Join(repoDir, "changed.txt")

	if err := os.WriteFile(
		filePath,
		[]byte("changed"),
		0644,
	); err != nil {
		t.Fatalf("write changed file: %v", err)
	}

	context, err := Detect(repoDir)
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	if !context.Dirty {
		t.Error("Dirty = false, want true")
	}
}

func TestDetectOutsideRepository(t *testing.T) {
	dir := t.TempDir()

	context, err := Detect(dir)
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	if context.IsRepository {
		t.Fatal("IsRepository = true, want false")
	}
}

func createTestRepository(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()

	runTestGit(t, repoDir, "init")
	runTestGit(t, repoDir, "config", "user.name", "WakeTrail Test")
	runTestGit(t, repoDir, "config", "user.email", "test@waketrail.local")

	filePath := filepath.Join(repoDir, "README.md")

	if err := os.WriteFile(
		filePath,
		[]byte("# Test Repository\n"),
		0644,
	); err != nil {
		t.Fatalf("write README: %v", err)
	}

	runTestGit(t, repoDir, "add", "README.md")
	runTestGit(t, repoDir, "commit", "-m", "initial commit")

	return repoDir
}

func runTestGit(
	t *testing.T,
	workingDirectory string,
	args ...string,
) {
	t.Helper()

	commandArgs := append(
		[]string{"-C", workingDirectory},
		args...,
	)

	cmd := exec.Command("git", commandArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf(
			"git %v failed: %v\n%s",
			args,
			err,
			string(output),
		)
	}
}
