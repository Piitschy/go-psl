package psl

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
)

var DomainTooShort = errors.New("Domain name too short.")
var DomainTooLong = errors.New("Domain name too long. It should be no more than 255 chars.")
var LabelStartsWithDash = errors.New("Domain label starts with a dash.")
var LabelEndsWithDash = errors.New("Domain label ends with a dash.")
var LabelTooLong = errors.New("Domain name label should be at most 63 chars long.")
var LabelTooShort = errors.New("Domain name label should be at least 1 character long.")
var LabelInvalidChar = errors.New("Domain name label can only contain alphanumeric characters or dashes.")

type Rule struct {
	rule       string
	suffix     string
	punySuffix int
	wildcard   bool
	exception  bool
}

func parseRule(rule string) *Rule {
	re := regexp.MustCompile(`^(\*\.|\!)`)
	return &Rule{
		rule:       rule,
		suffix:     re.ReplaceAllString(rule, ""),
		punySuffix: -1,
		wildcard:   strings.HasPrefix(rule, "*"),
		exception:  strings.HasPrefix(rule, "!"),
	}
}

func rules() []*Rule {
	//open ./data/rules.json
	dat, err := os.ReadFile("./data/rules.json")
	if err != nil {
		panic(err)
	}
	//parse json
	var rules []string
	err = json.Unmarshal(dat, &rules)
	if err != nil {
		panic(err)
	}
	//convert to Rule
	var ruleList []*Rule
	for _, rule := range rules {
		ruleList = append(ruleList, parseRule(rule))
	}
	return ruleList
}

func findRule() {}

func validate(input string) error {
	ascii := strings.ToLower(input)
	if len(ascii) < 1 {
		return DomainTooShort
	}
	if len(ascii) > 255 {
		return DomainTooLong
	}
	labels := strings.Split(ascii, ".")
	for _, label := range labels {
		if len(label) < 1 {
			return LabelTooShort
		}
		if len(label) > 63 {
			return LabelTooLong
		}
		if strings.HasPrefix(label, "-") {
			return LabelStartsWithDash
		}
		if strings.HasSuffix(label, "-") {
			return LabelEndsWithDash
		}
		if !regexp.MustCompile(`^[a-z0-9\-]+$`).MatchString(label) {
			return LabelInvalidChar
		}
	}
	return nil
}

func Parse(input string) (string, error) {
	return "", nil
}

func Get(domain string) (string, error) {
	return "", nil
}

func IsValid(domain string) bool {
	return false
}
