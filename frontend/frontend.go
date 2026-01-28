package frontend

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

func GetDist() fs.FS {
	fs, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	return fs
}
