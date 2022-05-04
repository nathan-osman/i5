package conman

import (
	"github.com/nathan-osman/i5/dockmon"
)

// Info stores information for a specific container.
type Info struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Domain  string `json:"domain"`
	Running bool   `json:"running"`
}

// Info returns information about running containers.
func (c *Conman) Info() []*Info {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ret := []*Info{}
	for d, v := range c.domainMap {
		con := v.(*dockmon.Container)
		if con.ID != "" {
			ret = append(ret, &Info{
				ID:      con.ID,
				Name:    con.Name,
				Domain:  d,
				Running: con.Running,
			})
		}
	}
	return ret
}
