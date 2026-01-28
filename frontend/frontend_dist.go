//go:build dist

package frontend

import (
	"embed"
	"io/fs"
	"log/slog"
	"os"
)

//go:embed all:dist
var dist embed.FS

func init() {
	os.Setenv("VITE_DIST", "true")
	slog.Info("Embedding frontend dist")
}

func GetDist() fs.FS {
	fs, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	return fs
}
