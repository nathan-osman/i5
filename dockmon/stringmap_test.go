package dockmon

import (
	"reflect"
	"testing"
)

const (
	stringMapKey  = ""
	stringMapKey1 = "1"
	stringMapKey2 = "2"
)

func TestStringMapInsertRemoveHas(t *testing.T) {
	s := StringMap{}
	s.Insert(stringMapKey, nil)
	if !s.Has(stringMapKey) {
		t.Fatal("string expected in map")
	}
	s.Remove(stringMapKey)
	if s.Has(stringMapKey) {
		t.Fatal("string found in map")
	}
}

func TestStringMapDifference(t *testing.T) {
	s1 := StringMap{}
	s1.Insert(stringMapKey, nil)
	s1.Insert(stringMapKey1, nil)
	s2 := StringMap{}
	s2.Insert(stringMapKey, nil)
	s2.Insert(stringMapKey2, nil)
	for _, c := range []struct {
		map1 StringMap
		map2 StringMap
		diff []string
	}{
		{map1: s1, map2: s2, diff: []string{stringMapKey1}},
		{map1: s2, map2: s1, diff: []string{stringMapKey2}},
	} {
		diff := c.map1.Difference(c.map2)
		if !reflect.DeepEqual(diff, c.diff) {
			t.Fatalf("%+v != %+v", diff, c.diff)
		}
	}
}
