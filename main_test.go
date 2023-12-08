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

func TestParseErrTooLong(t *testing.T) {
	var s string
	for len(s) < 256 {
		s += "x"
	}

	_, err := Parse(s)
	if err == nil {
		t.Fatal("should give a too-long error but gave nil")
	}
	if err != ErrDomainTooLong {
		t.Fatalf("should give a too-long error but gave:\n %v", err)
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		domain   string
		expected string
	}{
		{"", ""},
		{"COM", ""},
		{"example.COM", "example.com"},
		{"WwW.example.COM", "example.com"},
		{"example", ""},
		{"example.example", "example.example"},
		{"b.example.example", "example.example"},
		{"a.b.example.example", "example.example"},
		{"biz", ""},
		{"domain.biz", "domain.biz"},
		{"b.domain.biz", "domain.biz"},
		{"a.b.domain.biz", "domain.biz"},
		{"uk.com", ""},
		{"example.uk.com", "example.uk.com"},
		{"b.example.uk.com", "example.uk.com"},
		{"c.kobe.jp", ""},
		{"b.c.kobe.jp", "b.c.kobe.jp"},
		{"a.b.c.kobe.jp", "b.c.kobe.jp"},
		{"city.kobe.jp", "city.kobe.jp"},
		{"www.city.kobe.jp", "city.kobe.jp"},
		{"食狮.com.cn", "食狮.com.cn"},
		{"食狮.公司.cn", "食狮.公司.cn"},
		{"www.食狮.公司.cn", "食狮.公司.cn"},
		{"xn--85x722f.com.cn", "xn--85x722f.com.cn"},
		{"xn--85x722f.xn--55qx5d.cn", "xn--85x722f.xn--55qx5d.cn"},
		{"www.xn--85x722f.xn--55qx5d.cn", "xn--85x722f.xn--55qx5d.cn"},
	}

	for _, tc := range testCases {
		t.Run(tc.domain, func(t *testing.T) {
			result, _ := Get(tc.domain)
			if result != tc.expected {
				t.Errorf("psl.Get(%s) = %s; want %s", tc.domain, result, tc.expected)
			}
		})
	}
}
