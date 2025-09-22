package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	executableName string
	once           sync.Once
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

func GetExecutableName() string {
	if executableName != "" {
		return executableName
	}

	once.Do(func() {
		executableName = "sigolang"
		ex, err := os.Executable()

		if err == nil {
			dir := filepath.Dir(ex)
			if strings.Contains(dir, "go-build") {
				return
			}
			if exeName := filepath.Base(ex); exeName != "" {
				executableName = exeName
			}
		}
	})

	return executableName
}
