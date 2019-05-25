package util

// StringSet provides a set implementation for strings.
type StringSet struct {
	values map[string]interface{}
}

// NewStringSet creates a new set of strings.
func NewStringSet() *StringSet {
	return &StringSet{
		values: map[string]interface{}{},
	}
}

// Insert adds a value to the set, regardless of whether the value already exists.
func (s *StringSet) Insert(v string) {
	s.values[v] = nil
}

// Remove removes a value from the set.
func (s *StringSet) Remove(v string) {
	delete(s.values, v)
}

// Has checks for the provided value in the set.
func (s *StringSet) Has(v string) bool {
	_, ok := s.values[v]
	return ok
}

// Difference returns a list of strings that are in this set but not other.
func (s *StringSet) Difference(other *StringSet) []string {
	diff := []string{}
	for v := range s.values {
		if !other.Has(v) {
			diff = append(diff, v)
		}
	}
	return diff
}
