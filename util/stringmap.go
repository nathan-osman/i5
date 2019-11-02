package util

// StringMap provides a convenient interface for a map that uses strings for keys and stores interface{} values.
type StringMap map[string]interface{}

// Insert adds a key to the map and sets its value to the one provided.
func (s StringMap) Insert(k string, v interface{}) {
	s[k] = v
}

// Remove removes a value from the map.
func (s StringMap) Remove(k string) {
	delete(s, k)
}

// Has checks for the provided value in the map.
func (s StringMap) Has(k string) bool {
	_, ok := s[k]
	return ok
}

// Intersection returns a list of keys that are present in both maps.
func (s StringMap) Intersection(other StringMap) StringMap {
	ret := StringMap{}
	for k, v := range s {
		if other.Has(k) {
			ret[k] = v
		}
	}
	return ret
}

// Difference returns a list of keys that are in this map but not another.
func (s StringMap) Difference(other StringMap) StringMap {
	diff := StringMap{}
	for k, v := range s {
		if !other.Has(k) {
			diff[k] = v
		}
	}
	return diff
}
