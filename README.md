![Version](https://img.shields.io/badge/version-1.0.16-orange.svg)
[![GolangCI Lint](https://github.com/deliveryhero/spring-cloud-config-client-go/actions/workflows/go-lint.yml/badge.svg)](https://github.com/deliveryhero/sc-payment-service/actions/workflows/go-lint.yml)
[![Golang Tests](https://github.com/deliveryhero/spring-cloud-config-client-go/actions/workflows/go-test.yml/badge.svg)](https://github.com/deliveryhero/sc-payment-service/actions/workflows/go-test.yml) 
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit) 
![Test Coverage](https://img.shields.io/badge/coverage-84.7%25-orange.svg)

# spring-cloud-config-client-go
Spring Cloud Config Client is GO client for [Spring Cloud Config](https://docs.spring.io/spring-cloud-config/docs/current/reference/html/). Supports [Property Overrides](https://docs.spring.io/spring-cloud-config/docs/current/reference/html/#property-overrides) feature.

## Sample Usage

```go
package main

import (
	"fmt"

	configclient "github.com/deliveryhero/spring-cloud-config-client-go/springconfigclient"
)

func main() {
	c := configclient.RemoteConfig{
		Url:      "https://remote-url.com",
		Username: "username",
		Password: "pass",
	}
	a := configclient.New("sample-api", "prod", &c)

	if err := a.Sync(); err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("ENV_KEY value is: ", a.GetEnv("ENV_KEY"))
}
```

## Installation

You can add this package via;

```bash
go get github.com/deliveryhero/spring-cloud-config-client-go
```

---


## Tests

To run tests, use `rake test` or;

```bash
go test -p 1 -v -race ./...
```

## Publishing New Release

Prerequisites

- You need to be in `main` branch
- You need to be ready to bump to a new version

Use `rake publish[revision]` task to bump new version and push newly created
tag and updated code to remote and verify go package. (all in one!)

- `rake publish`: `0.0.0` -> `0.0.1`, default revision is `patch`
- `rake publish[minor]`: `0.0.0` -> `0.1.0`
- `rake publish[major]`: `0.0.0` -> `1.0.0`

---

## Contributor(s)

* [Erhan Akpınar](https://github.com/erhanakp) - Creator, maintainer
* [Hakan Kutluay](https://github.com/hakankutluay) - Contributor

---

## Contribute

All PR’s are welcome!

1. `fork` (https://github.com/deliveryhero/spring-cloud-config-client-go/fork)
2. Create your `branch` (`git checkout -b my-feature`)
3. `commit` yours (`git commit -am 'add some functionality'`)
4. `push` your `branch` (`git push origin my-feature`)
5. Than create a new **Pull Request**!

This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].


[coc]: https://github.com/deliveryhero/spring-cloud-config-client-go/blob/main/CODE_OF_CONDUCT.md