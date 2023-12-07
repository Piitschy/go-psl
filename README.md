psl (Public Suffix List)
Go CI

psl is a Go domain name parser based on the Public Suffix List.

This Go implementation is tested against the test data hosted by Mozilla and kindly provided by Comodo.

What is the Public Suffix List?
The Public Suffix List is a cross-vendor initiative to provide an accurate list of domain name suffixes.

Initially created to meet the needs of browser manufacturers, it's now a community-maintained resource available for any software use. It lists all known public suffixes, which are domain parts where Internet users can directly register names. Examples include ".com", ".co.uk", and "pvt.k12.wy.us".

Source: http://publicsuffix.org

Installation
Go
sh
Copy code
go get github.com/Piitschy/psl
API
psl.Parse(domain string)
Parse a domain based on the Public Suffix List. Returns a struct with the following properties:

TLD: Top level domain (the public suffix).
SLD: Second level domain (the first private part of the domain name).
Domain: The combination of SLD and TLD.
Subdomain: Any optional parts left of the domain.
Example:
go
Copy code
package main

import (
    "fmt"
    "github.com/Piitschy/psl"
)

func main() {
    parsed, _ := psl.Parse("www.google.com")
    fmt.Println(parsed.TLD) // 'com'
    fmt.Println(parsed.SLD) // 'google'
    fmt.Println(parsed.Domain) // 'google.com'
    fmt.Println(parsed.Subdomain) // 'www'
}
psl.Get(domain string)
Get the domain name, SLD + TLD. Returns an empty string if not valid.

Example:
go
Copy code
// [Similar examples as provided in the JavaScript version, rewritten in Go]
psl.IsValid(domain string)
Check whether a domain has a valid Public Suffix. Returns a bool indicating the validity.

Example
go
Copy code
// [Similar examples as provided in the JavaScript version, rewritten in Go]
Testing and Building
Tests are written using Go's built-in testing framework. To run tests:

sh
Copy code
go test
Feel free to fork if you see possible improvements!

Acknowledgements
Mozilla Foundation's Public Suffix List
Rewrite from lupomontero/psl

