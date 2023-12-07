

# psl (Public Suffix List)

This is a Go port from [lupomontero/psl](https://github.com/lupomontero/psl).

## What is the Public Suffix List?

The Public Suffix List is a cross-vendor initiative to provide an accurate list of domain name suffixes.

Initially created to meet the needs of browser manufacturers, it's now a community-maintained resource available for any software use. It lists all known public suffixes, which are domain parts where Internet users can directly register names. Examples include ".com", ".co.uk", and "pvt.k12.wy.us".

Source: http://publicsuffix.org

## Installation

### Go

```sh
go get github.com/Piitschy/psl
```

## API

### `psl.Parse(domain string)`

Parse a domain based on the Public Suffix List. Returns a `struct` with the following properties:

* `Tld`: Top level domain (the public suffix).
* `Sld`: Second level domain (the first private part of the domain name).
* `Domain`: The combination of `Sld` and `Tld`.
* `Subdomain`: Any optional parts left of the domain.

#### Example:

```go
package main

import (
    "fmt"
    "github.com/Piitschy/psl"
)

func main() {
    parsed, _ := psl.Parse("www.google.com")
    fmt.Println(parsed.Tld) // 'com'
    fmt.Println(parsed.Sld) // 'google'
    fmt.Println(parsed.Domain) // 'google.com'
    fmt.Println(parsed.Subdomain) // 'www'
}
```

### `psl.Get(domain string)`

Get the domain name, `Sld` + `Tld`. Returns an empty string if not valid.

#### Example:

```go
package main

import (
    "fmt"
    "github.com/Piitschy/psl"
)

func main() {
    // null input.
    fmt.Println(psl.Get("")) // should print an empty string

    // Mixed case.
    fmt.Println(psl.Get("COM")) // should print an empty string
    fmt.Println(psl.Get("example.COM")) // 'example.com'
    fmt.Println(psl.Get("WwW.example.COM")) // 'example.com'

    // Unlisted TLD.
    fmt.Println(psl.Get("example")) // should print an empty string
    fmt.Println(psl.Get("example.example")) // 'example.example'
    fmt.Println(psl.Get("b.example.example")) // 'example.example'
    fmt.Println(psl.Get("a.b.example.example")) // 'example.example'

    // TLD with only 1 rule.
    fmt.Println(psl.Get("biz")) // should print an empty string
    fmt.Println(psl.Get("domain.biz")) // 'domain.biz'
    fmt.Println(psl.Get("b.domain.biz")) // 'domain.biz'
    fmt.Println(psl.Get("a.b.domain.biz")) // 'domain.biz'

    // TLD with some 2-level rules.
    fmt.Println(psl.Get("uk.com")) // should print an empty string
    fmt.Println(psl.Get("example.uk.com")) // 'example.uk.com'
    fmt.Println(psl.Get("b.example.uk.com")) // 'example.uk.com'

    // More complex TLD.
    fmt.Println(psl.Get("c.kobe.jp")) // should print an empty string
    fmt.Println(psl.Get("b.c.kobe.jp")) // 'b.c.kobe.jp'
    fmt.Println(psl.Get("a.b.c.kobe.jp")) // 'b.c.kobe.jp'
    fmt.Println(psl.Get("city.kobe.jp")) // 'city.kobe.jp'
    fmt.Println(psl.Get("www.city.kobe.jp")) // 'city.kobe.jp'

    // IDN labels.
    fmt.Println(psl.Get("食狮.com.cn")) // '食狮.com.cn'
    fmt.Println(psl.Get("食狮.公司.cn")) // '食狮.公司.cn'
    fmt.Println(psl.Get("www.食狮.公司.cn")) // '食狮.公司.cn'

    // Same as above, but punycoded.
    fmt.Println(psl.Get("xn--85x722f.com.cn")) // 'xn--85x722f.com.cn'
    fmt.Println(psl.Get("xn--85x722f.xn--55qx5d.cn")) // 'xn--85x722f.xn--55qx5d.cn'
    fmt.Println(psl.Get("www.xn--85x722f.xn--55qx5d.cn")) // 'xn--85x722f.xn--55qx5d.cn'
}

```

### `psl.IsValid(domain string)`

Check whether a domain has a valid Public Suffix. Returns a `bool` indicating the validity.

#### Example

```go
package main

import (
    "fmt"
    "github.com/Piitschy/psl"
)

func main() {
    fmt.Println(psl.IsValid("google.com"))       // true
    fmt.Println(psl.IsValid("www.google.com"))   // true
    fmt.Println(psl.IsValid("x.yz"))             // false
}
```

## Testing and Building

Tests are written using Go's built-in testing framework. To run tests:

```sh
go test
```

Feel free to fork if you see possible improvements!

## Acknowledgements

* Mozilla Foundation's [Public Suffix List](https://publicsuffix.org/)
* Inspired by [lupomontero/psl](https://github.com/lupomontero/psl)

## License

[MIT License](LICENSE.md)