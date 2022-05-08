package conman

import (
	"github.com/nathan-osman/i5/container"
)

// Info returns information about running containers.
func (c *Conman) Info() []container.ContainerData {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ret := []container.ContainerData{}
	for _, v := range c.idMap {
		con := v.(*container.Container)
		if con.ID != "" {
			ret = append(ret, con.ContainerData)
		}
	}
	return ret
}
