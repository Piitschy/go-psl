package psl

import (
	"testing"
)

func TestParse(t *testing.T) {
	google := Domain{
		"www.google.com",
		"com",
		"google",
		"google.com",
		"www",
		true,
	}
	parsed, err := Parse("www.google.com")
	if err != nil {
		t.Fail()
	}
	if *parsed != google {
		t.Fatal("www.google.com not parsed right")
	}
}

func TestParseTooLong(t *testing.T) {
	var s string
	for len(s) < 256 {
		s += "x"
	}

	_, err := Parse(s)
	if err == nil {
		t.Fatal("should give a too-long error but gave nil")
	}
	if err != DomainTooLong {
		t.Fatalf("should give a too-long error but gave:\n %v", err)
	}
}
