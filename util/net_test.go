package util

import (
	"testing"
)

func TestParseHost(t *testing.T) {
	for _, tt := range []struct {
		Addr string
		Host string
	}{
		{Addr: ":"},
		{Addr: "host:", Host: "host"},
	} {
		if host := ParseHost(tt.Addr); host != tt.Host {
			t.Fatalf("%s != %s", host, tt.Host)
		}
	}
}

func TestParsePort(t *testing.T) {
	for _, tt := range []struct {
		Addr  string
		Def   int
		Port  int
		Error bool
	}{
		{Addr: ":", Def: 0, Port: 0},
		{Addr: "host", Def: 80, Port: 80},
		{Addr: "host:80", Def: 0, Port: 80},
		{Addr: "host:http", Def: 0, Port: 80},
		{Addr: "host:invalid", Def: 0, Port: 0, Error: true},
	} {
		p, err := ParsePort(tt.Addr, tt.Def)
		if p != tt.Port {
			t.Fatalf("%d != %d", p, tt.Port)
		}
		if (err != nil) != tt.Error {
			if err != nil {
				t.Fatal(err)
			} else {
				t.Fatal("error expected")
			}
		}
	}
}
