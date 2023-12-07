Here's an edited version of the README for your `psl` project, now rewritten in Go:

---

# psl (Public Suffix List)

## What is the Public Suffix List?

The Public Suffix List is a cross-vendor initiative to provide an accurate list of domain name suffixes.

Initially created to meet the needs of browser manufacturers, it's now a community-maintained resource available for any software use. It lists all known public suffixes, which are domain parts where Internet users can directly register names. Examples include ".com", ".co.uk", and "pvt.k12.wy.us".

Source: http://publicsuffix.org

## Installation

### Go

```sh
go get github.com/[your-github-username]/psl
```

## API

### `psl.Parse(domain string)`

Parse a domain based on the Public Suffix List. Returns a `struct` with the following properties:

* `TLD`: Top level domain (the public suffix).
* `SLD`: Second level domain (the first private part of the domain name).
* `Domain`: The combination of `SLD` and `TLD`.
* `Subdomain`: Any optional parts left of the domain.

#### Example:

```go
package main

import (
    "fmt"
    "github.com/[your-github-username]/psl"
)

func main() {
    parsed, _ := psl.Parse("www.google.com")
    fmt.Println(parsed.TLD) // 'com'
    fmt.Println(parsed.SLD) // 'google'
    fmt.Println(parsed.Domain) // 'google.com'
    fmt.Println(parsed.Subdomain) // 'www'
}
```

### `psl.Get(domain string)`

Get the domain name, `SLD` + `TLD`. Returns an empty string if not valid.

#### Example:

```go
// [Similar examples as provided in the JavaScript version, rewritten in Go]
```

### `psl.IsValid(domain string)`

Check whether a domain has a valid Public Suffix. Returns a `bool` indicating the validity.

#### Example

```go
// [Similar examples as provided in the JavaScript version, rewritten in Go]
```

## Testing and Building

Tests are written using Go's built-in testing framework. To run tests:

```sh
go test
```

Feel free to fork if you see possible improvements!

## Acknowledgements

* Mozilla Foundation's [Public Suffix List](https://publicsuffix.org/)
* Thanks to Rob Stradling of [Comodo](https://www.comodo.com/) for providing test data.
* Inspired by [weppos/publicsuffix-ruby](https://github.com/weppos/publicsuffix-ruby)

## License

[MIT License](LICENSE.md)

---

Remember to replace `[your-github-username]` with your actual GitHub username, and adjust any paths or URLs as necessary for your project.