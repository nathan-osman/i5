package ui

import (
	"embed"
	"io/fs"
)

var (
	//go:embed build/*
	content    embed.FS
	Content, _ = fs.Sub(content, "build")
)
