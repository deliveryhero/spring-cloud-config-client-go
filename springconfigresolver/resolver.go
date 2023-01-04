package resolver

import (
	"os"
	"regexp"
	"strings"
)

type resolver struct {
	matcher *regexp.Regexp
}

type Resolver interface {
	Resolve(value string) (string, bool)
}

var _ Resolver = (*resolver)(nil) // compile time proof

func New() Resolver {
	matcher := regexp.MustCompile(`\{(.*?)\}`)
	return &resolver{
		matcher: matcher,
	}
}

func (r *resolver) Resolve(str string) (string, bool) {
	if len(str) > 0 && str[0] == '$' {
		result := r.matcher.FindStringSubmatch(str)
		if len(result) == 0 {
			return str, false
		}

		v := strings.Split(result[1], ":")
		var defaultValue *string
		envValueName := v[0]
		if len(v) > 1 {
			defaultPart := strings.Join(v[1:], ":")
			defaultValue = &defaultPart
		}

		envValue, envOk := os.LookupEnv(envValueName)

		if envValue == "" && defaultValue != nil {
			return *defaultValue, true
		}

		return envValue, envOk
	}

	return str, true
}
