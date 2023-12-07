package psl

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/idna"
)

var DomainTooShort = errors.New("Domain name too short.")
var DomainTooLong = errors.New("Domain name too long. It should be no more than 255 chars.")
var LabelStartsWithDash = errors.New("Domain label starts with a dash.")
var LabelEndsWithDash = errors.New("Domain label ends with a dash.")
var LabelTooLong = errors.New("Domain name label should be at most 63 chars long.")
var LabelTooShort = errors.New("Domain name label should be at least 1 character long.")
var LabelInvalidChar = errors.New("Domain name label can only contain alphanumeric characters or dashes.")

type Rule struct {
	Rule       string
	Suffix     string
	PunySuffix string
	Wildcard   bool
	Exception  bool
}

func parseRule(rule string) *Rule {
	re := regexp.MustCompile(`^(\*\.|\!)`)
	return &Rule{
		Rule:       rule,
		Suffix:     re.ReplaceAllString(rule, ""),
		PunySuffix: "",
		Wildcard:   strings.HasPrefix(rule, "*"),
		Exception:  strings.HasPrefix(rule, "!"),
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

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

func findRule(domain string) *Rule {
	punyDomain, err := idna.ToASCII(domain)
	if err != nil {
		panic(err)
	}
	var matchedRule *Rule
	for _, rule := range rules() {
		if rule.PunySuffix == "" {
			punySuffix, err := idna.ToASCII(rule.Suffix)
			if err != nil {
				continue
			}
			rule.PunySuffix = punySuffix

		}
		if !endsWith(punyDomain, "."+rule.PunySuffix) && punyDomain != rule.PunySuffix {
			continue
		}
		matchedRule = rule
		break
	}
	return matchedRule

}

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

type Domain struct {
	Input     string
	Tld       string
	Sld       string
	Domain    string
	Subdomain string
	Listed    bool
}

func (d *Domain) handlePunycode() error {
	if !strings.Contains(d.Input, "xn--") {
		return nil
	}
	if d.Domain == "" {
		d.Domain, _ = idna.ToASCII(d.Domain)
	}
	if d.Subdomain == "" {
		d.Subdomain, _ = idna.ToASCII(d.Subdomain)
	}
	return nil
}

// Parse domain.
func Parse(input string) (*Domain, error) {
	domain := strings.ToLower(input)
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}

	parsed := &Domain{
		Input:     input,
		Tld:       "",
		Sld:       "",
		Domain:    "",
		Subdomain: "",
		Listed:    false,
	}

	if err := validate(domain); err != nil {
		return parsed, err
	}

	domainParts := strings.Split(domain, ".")
	if domainParts[len(domainParts)-1] == "local" {
		return parsed, nil // assuming 'parsed' is a variable defined in your function that you want to return
	}

	rule := findRule(domain)
	if rule == nil {
		if len(domainParts) == 1 {
			parsed.Domain = domain
			return parsed, nil
		}
		parsed.Tld = domainParts[len(domainParts)-1]
		parsed.Sld = domainParts[len(domainParts)-2]
		parsed.Domain = strings.Join([]string{parsed.Sld, parsed.Tld}, ".")
		if len(domainParts) > 2 {
			parsed.Subdomain = strings.Join(domainParts[:len(domainParts)-2], ".")
		}
		parsed.handlePunycode()
		return parsed, nil
	}

	parsed.Listed = true
	tldParts := strings.Split(rule.Suffix, ".")
	privateParts := domainParts[:len(domainParts)-len(tldParts)]

	var x string
	if rule.Exception {
		x, tldParts = tldParts[0], tldParts[1:]
		privateParts = append(privateParts, x)
	}

	parsed.Tld = strings.Join(tldParts, ".")

	if len(privateParts) == 0 {
		parsed.handlePunycode()
		return parsed, nil
	}

	if rule.Wildcard {
		x, privateParts = privateParts[len(privateParts)-1], privateParts[:len(privateParts)-1]
		tldParts = append([]string{x}, tldParts...)
		parsed.Tld = strings.Join(tldParts, ".")
	}

	if len(privateParts) == 0 {
		parsed.handlePunycode()
		return parsed, nil
	}

	parsed.Sld, privateParts = privateParts[len(privateParts)-1], privateParts[:len(privateParts)-1]
	parsed.Domain = strings.Join([]string{parsed.Sld, parsed.Tld}, ".")

	if len(privateParts) > 0 {
		parsed.Subdomain = strings.Join(privateParts, ".")
	}

	parsed.handlePunycode()
	return parsed, nil
}

// Get domain.
func Get(domain string) (string, error) {
	if domain == "" {
		return "", errors.New("Empty domain.")
	}
	parsed, err := Parse(domain)
	return parsed.Domain, err
}

// Check whether domain belongs to a known public suffix.
func IsValid(domain string) bool {
	parsed, err := Parse(domain)
	if err != nil {
		return false
	}
	return parsed.Domain != ""
}
