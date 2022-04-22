package ui

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
)

var (
	//go:embed dist
	content    embed.FS
	Content, _ = fs.Sub(content, "dist")
)

// EmbedFileSystem implements the ServeFileSystem interface.
type EmbedFileSystem struct {
	http.FileSystem
}

func (e EmbedFileSystem) Exists(prefix, path string) bool {
	f, err := e.Open(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	defer f.Close()
	return true
}
