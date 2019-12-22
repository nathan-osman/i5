// +build ignore

package main

import (
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	if err := vfsgen.Generate(http.Dir("ui/build"), vfsgen.Options{
		Filename:     "ui/vfsdata.go",
		PackageName:  "ui",
		VariableName: "Assets",
	}); err != nil {
		panic(err)
	}
}
