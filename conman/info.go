package conman

import (
	"github.com/nathan-osman/i5/dockmon"
)

// Info stores information for a specific container.
type Info struct {
	Name    string
	Running bool
}

// Info returns information about running containers.
func (c *Conman) Info() []*Info {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ret := []*Info{}
	for _, v := range c.domainMap {
		con := v.(*dockmon.Container)
		if con.ID != "" {
			ret = append(ret, &Info{
				Name:    con.Name,
				Running: con.Running,
			})
		}
	}
	return ret
}
