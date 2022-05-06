package conman

import (
	"github.com/nathan-osman/i5/container"
)

// Info stores information for a specific container.
type Info struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Domains  []string `json:"domains"`
	Disabled bool     `json:"disabled"`
}

// Info returns information about running containers.
func (c *Conman) Info() []*Info {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ret := []*Info{}
	for _, v := range c.idMap {
		con := v.(*container.Container)
		if con.ID != "" {
			ret = append(ret, &Info{
				ID:       con.ID,
				Name:     con.Name,
				Domains:  con.Domains,
				Disabled: con.Disabled,
			})
		}
	}
	return ret
}
