package psl

import (
	"errors"
	"regexp"
	"strings"

	"github.com/Piitschy/psl/data"
	"github.com/Piitschy/psl/utils"
	"golang.org/x/net/idna"
)

var (
	ErrDomainTooShort      = errors.New("Domain name too short")
	ErrDomainTooLong       = errors.New("Domain name too long. It should be no more than 255 chars")
	ErrLabelStartsWithDash = errors.New("Domain label starts with a dash")
	ErrLabelEndsWithDash   = errors.New("Domain label ends with a dash")
	ErrLabelTooLong        = errors.New("Domain name label should be at most 63 chars long")
	ErrLabelTooShort       = errors.New("Domain name label should be at least 1 character long")
	ErrLabelInvalidChar    = errors.New("Domain name label can only contain alphanumeric characters or dashes")
)

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
	//convert to Rule
	var ruleList []*Rule
	for _, rule := range data.Rules {
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
	rules := rules()
	for _, rule := range rules {
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
	}
	return matchedRule

}

func validate(input string) error {
	ascii, err := idna.ToASCII(strings.ToLower(input))
	if err != nil {
		return ErrLabelInvalidChar
	}
	if len(ascii) < 1 {
		return ErrDomainTooShort
	}
	if len(ascii) > 255 {
		return ErrDomainTooLong
	}
	labels := strings.Split(ascii, ".")
	for _, label := range labels {
		if len(label) < 1 {
			return ErrLabelTooShort
		}
		if len(label) > 63 {
			return ErrLabelTooLong
		}
		if strings.HasPrefix(label, "-") {
			return ErrLabelStartsWithDash
		}
		if strings.HasSuffix(label, "-") {
			return ErrLabelEndsWithDash
		}
		if !regexp.MustCompile(`^[a-z0-9\-]+$`).MatchString(label) {
			return ErrLabelInvalidChar
		}
	}
	return nil
}

type Domain struct {
	Input     string
	TLD       string
	SLD       string
	Domain    string
	Subdomain string
	Listed    bool
}

func (d *Domain) handlePunycode() *Domain {
	if !strings.Contains(strings.ToLower(d.Input), "xn--") {
		return d
	}
	if d.Domain != "" {
		d.Domain, _ = idna.ToASCII(d.Domain)
	}
	if d.Subdomain != "" {
		d.Subdomain, _ = idna.ToASCII(d.Subdomain)
	}
	return d
}

// Parse domain.
func Parse(input string) (*Domain, error) {
	domain := strings.ToLower(input)
	if strings.HasPrefix(domain, "http") {
		domain = strings.Split(domain, "/")[2]
	}
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}

	parsed := &Domain{
		Input:     input,
		TLD:       "",
		SLD:       "",
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
		if len(domainParts) < 2 {
			return parsed, nil
		}
		parsed.TLD = utils.Popp(&domainParts)
		parsed.SLD = utils.Popp(&domainParts)
		parsed.Domain = strings.Join([]string{parsed.SLD, parsed.TLD}, ".")
		if len(domainParts) > 0 {
			parsed.Subdomain = utils.Popp(&domainParts)
		}
		return parsed.handlePunycode(), nil
	}

	parsed.Listed = true
	tldParts := strings.Split(rule.Suffix, ".")
	privateParts := domainParts[:len(domainParts)-len(tldParts)]

	var x string
	if rule.Exception {
		x, tldParts = tldParts[0], tldParts[1:]
		privateParts = append(privateParts, x)
	}

	parsed.TLD = strings.Join(tldParts, ".")

	if len(privateParts) == 0 {
		return parsed.handlePunycode(), nil
	}

	if rule.Wildcard {
		x, privateParts = privateParts[len(privateParts)-1], privateParts[:len(privateParts)-1]
		tldParts = append([]string{x}, tldParts...)
		parsed.TLD = strings.Join(tldParts, ".")
	}

	if len(privateParts) == 0 {
		return parsed.handlePunycode(), nil
	}

	parsed.SLD = utils.Popp(&privateParts)
	parsed.Domain = strings.Join([]string{parsed.SLD, parsed.TLD}, ".")

	if len(privateParts) > 0 {
		parsed.Subdomain = strings.Join(privateParts, ".")
	}

	return parsed.handlePunycode(), nil
}

// Get domain.
func Get(domain string) (string, error) {
	if domain == "" {
		return "", errors.New("empty domain")
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
