package util

import (
	"reflect"
	"testing"
)

const (
	stringMapStr  = ""
	stringMapStr1 = "1"
	stringMapStr2 = "2"
)

func TestStringMapInsertRemoveHas(t *testing.T) {
	s := StringMap{}
	s.Insert(stringMapStr, stringMapStr)
	if !s.Has(stringMapStr) {
		t.Fatal("string expected in map")
	}
	s.Remove(stringMapStr)
	if s.Has(stringMapStr) {
		t.Fatal("string found in map")
	}
}

func TestStringMapDifference(t *testing.T) {
	s1 := StringMap{}
	s1.Insert(stringMapStr, stringMapStr)
	s1.Insert(stringMapStr1, stringMapStr1)
	s2 := StringMap{}
	s2.Insert(stringMapStr, stringMapStr)
	s2.Insert(stringMapStr2, stringMapStr2)
	for _, c := range []struct {
		map1 StringMap
		map2 StringMap
		diff StringMap
	}{
		{map1: s1, map2: s2, diff: StringMap{stringMapStr1: stringMapStr1}},
		{map1: s2, map2: s1, diff: StringMap{stringMapStr2: stringMapStr2}},
	} {
		diff := c.map1.Difference(c.map2)
		if !reflect.DeepEqual(diff, c.diff) {
			t.Fatalf("%+v != %+v", diff, c.diff)
		}
	}
}
