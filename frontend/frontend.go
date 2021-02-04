package frontend

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
)

//go:embed build/*
var assets embed.FS

func GetFilesystem() http.FileSystem {
	return http.FS(&scopedFilesystem{
		backend: assets,
		scope:   "build/",
	})
}

type scopedFilesystem struct {
	backend embed.FS
	scope   string
}

func (s scopedFilesystem) Open(name string) (fs.File, error) {
	return s.backend.Open(path.Join(s.scope, name))
}
