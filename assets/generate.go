// +build ignore

package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/shurcooL/vfsgen"
)

// https://github.com/shurcooL/vfsgen/issues/68

type FilterFile struct {
	http.File
}

func (f FilterFile) Readdir(count int) ([]os.FileInfo, error) {
	var files []os.FileInfo
	i, err := f.File.Readdir(count)
	if err != nil {
		return nil, err
	}
	for _, v := range i {
		if v.IsDir() || !strings.HasSuffix(v.Name(), ".go") {
			files = append(files, v)
		}
	}
	return files, nil
}

type FilterFS struct {
	http.FileSystem
}

func (f FilterFS) Open(name string) (http.File, error) {
	v, err := f.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return FilterFile{v}, err
}

func main() {
	if err := vfsgen.Generate(FilterFS{http.Dir("assets")}, vfsgen.Options{
		Filename:     "assets/vfsdata.go",
		PackageName:  "assets",
		VariableName: "Assets",
	}); err != nil {
		panic(err)
	}
}
