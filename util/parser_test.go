package util

import (
	"strings"
	"testing"
)

type stringPair struct {
	k string
	v string
}

func TestParser(t *testing.T) {
	for _, c := range []struct {
		input  string
		output []stringPair
		err    error
	}{
		{input: "", err: ErrEOF},
		{input: "k=", output: []stringPair{{k: "k", v: ""}}, err: ErrEOF},
		{input: "k=v", output: []stringPair{{k: "k", v: "v"}}, err: ErrEOF},
		{input: "k=\"v\"", output: []stringPair{{k: "k", v: "v"}}, err: ErrEOF},
		{input: " k = v ", output: []stringPair{{k: "k", v: "v"}}, err: ErrEOF},
		{input: " k = \"v\" ", output: []stringPair{{k: "k", v: "v"}}, err: ErrEOF},
		{input: "k1=v1 k2=v2", output: []stringPair{{k: "k1", v: "v1"}, {k: "k2", v: "v2"}}, err: ErrEOF},
		{input: "k", err: ErrEqualsExpected},
		{input: "k=\"", err: ErrQuoteExpected},
		{input: "k=\"\\\"", err: ErrQuoteExpected},
	} {
		p, err := NewParser(strings.NewReader(c.input))
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range c.output {
			k, v, err := p.ParseNextEntry()
			if err != nil {
				t.Fatal(err)
			}
			if o.k != k {
				t.Fatalf("%s != %s", o.k, k)
			}
			if o.v != v {
				t.Fatalf("%s != %s", o.v, v)
			}
		}
		if _, _, err := p.ParseNextEntry(); err != c.err {
			t.Fatalf("%v != %v", err, c.err)
		}
	}
}
