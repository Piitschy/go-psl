package utils

import (
	"testing"
)

func TestPop(t *testing.T) {
	s := []string{"a", "b", "c"}
	val, s := Pop(s)
	if val != "c" {
		t.Errorf("should return 'c' but gave %v", val)
	}
	if s[0] != "a" || s[1] != "b" {
		t.Fail()
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should remove 'c' but doesn't")
		}
	}()
	_ = s[2]
}

func TestPopp(t *testing.T) {
	s := []string{"a", "b", "c"}
	val := Popp(&s)
	if val != "c" {
		t.Errorf("should return 'c' but gave %v", val)
	}
	if s[0] != "a" || s[1] != "b" {
		t.Fail()
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("should remove 'c' but doesn't")
		}
	}()
	_ = s[2]
}
