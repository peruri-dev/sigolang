//go:build !dist

package frontend

import (
	"io/fs"
	"os"
)

func GetDist() fs.FS {
	fs := os.DirFS("./frontend/dist")
	return fs
}
