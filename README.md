![Version](https://img.shields.io/badge/version-1.0.8-orange.svg)
[![GolangCI Lint](https://github.com/deliveryhero/spring-cloud-config-client-go/actions/workflows/go-lint.yml/badge.svg)](https://github.com/deliveryhero/sc-payment-service/actions/workflows/go-lint.yml)
[![Golang Tests](https://github.com/deliveryhero/spring-cloud-config-client-go/actions/workflows/go-test.yml/badge.svg)](https://github.com/deliveryhero/sc-payment-service/actions/workflows/go-test.yml) 
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit) 
![Test Coverage](https://img.shields.io/badge/coverage-76.3%25-orange.svg)

# spring-cloud-config-client-go
Spring Cloud Config Client

## Installation

This is a private repo, you need to fix your git/ssh or token operations before
injecting to you go app.

If you are using `git+ssh`, I’m assuming that you have already authorized your
ssh-key key on [GitHub](https://github.com/settings/tokens) side under
**Personal access tokens** page via **Configure SSO** combo-box.

Fix your git configuration, run this:

```bash
git config --global --add url."git@github.com:".insteadOf "https://github.com/"
```

command above adds few lines to your `~/.gitconfig`:

```ini
[url "git@github.com:"]
	insteadOf = https://github.com/
```

Now, add `GOPRIVATE` environment value to your shell environment.

```bash
export GOPRIVATE="github.com/deliveryhero"
```

Now you can add this package via;

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

1. `fork` (https://github.com/deliveryhero/sc-honeylogger/fork)
1. Create your `branch` (`git checkout -b my-feature`)
1. `commit` yours (`git commit -am 'add some functionality'`)
1. `push` your `branch` (`git push origin my-feature`)
1. Than create a new **Pull Request**!

This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].


[coc]: https://github.com/deliveryhero/sc-honeylogger/blob/main/CODE_OF_CONDUCT.md