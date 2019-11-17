package dockmon

import (
	"fmt"
	"strings"

	"github.com/nathan-osman/i5/util"
)

// Database represents database requirements for a container.
type Database struct {
	Driver   string
	Name     string
	User     string
	Password string
}

// ParseDatabaseString parses a string in the "k1=v1 k2=v2 ..." format.
func ParseDatabaseString(dbStr string) (*Database, error) {
	d := &Database{}
	p, err := util.NewParser(strings.NewReader(dbStr))
	if err != nil {
		return nil, err
	}
	for {
		k, v, err := p.ParseNextEntry()
		if err != nil {
			if err != util.ErrEOF {
				return nil, err
			}
			break
		}
		switch k {
		case "driver":
			d.Driver = v
		case "name":
			d.Name = v
		case "user":
			d.User = v
		case "password":
			d.Password = v
		default:
			return nil, fmt.Errorf("invalid key \"%s\"", k)
		}
	}
	return d, nil
}
