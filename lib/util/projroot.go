package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetCurrentProjectRoot() (string, error) {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	// The `b` variable contains the absolute path of the current source file.
	// We go up one directory to reach the project root (assuming the executable is directly in the root).
	// Adjust this if your executable is in a subdirectory (e.g., "bin").
	projectRoot := filepath.Join(filepath.Dir(b), "..", "..")

	// Clean the path to remove any redundant separators or ".." components.
	projectRoot = filepath.Clean(projectRoot)

	return projectRoot, nil
}

func GetExecutablePath() string {
	executablePath, err := os.Executable()
	if err != nil {
		executablePath = "sigolang"
	}
	return executablePath
}
