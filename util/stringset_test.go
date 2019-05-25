package util

import (
	"reflect"
	"testing"
)

const (
	stringSetValue  = ""
	stringSetValue1 = "1"
	stringSetValue2 = "2"
)

func TestStringSetInsertRemoveHas(t *testing.T) {
	s := NewStringSet()
	s.Insert(stringSetValue)
	if !s.Has(stringSetValue) {
		t.Fatal("string expected in set")
	}
	s.Remove(stringSetValue)
	if s.Has(stringSetValue) {
		t.Fatal("string found in set")
	}
}

func TestStringSetDifference(t *testing.T) {
	s1 := NewStringSet()
	s1.Insert(stringSetValue)
	s1.Insert(stringSetValue1)
	s2 := NewStringSet()
	s2.Insert(stringSetValue)
	s2.Insert(stringSetValue2)
	for _, c := range []struct {
		set1 *StringSet
		set2 *StringSet
		diff []string
	}{
		{set1: s1, set2: s2, diff: []string{stringSetValue1}},
		{set1: s2, set2: s1, diff: []string{stringSetValue2}},
	} {
		diff := c.set1.Difference(c.set2)
		if !reflect.DeepEqual(diff, c.diff) {
			t.Fatalf("%+v != %+v", diff, c.diff)
		}
	}
}
