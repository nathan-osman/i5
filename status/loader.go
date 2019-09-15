package status

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/nathan-osman/i5/assets"
)

type vfsgenLoader struct{}

func (v *vfsgenLoader) Abs(base, name string) string {
	return name
}

func (v *vfsgenLoader) Get(path string) (io.Reader, error) {
	f, err := assets.Assets.Open(fmt.Sprintf("templates/%s", path))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
