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
	Resolve(value string) string
}

var _ Resolver = (*resolver)(nil) // compil

func New() Resolver {
	matcher := regexp.MustCompile(`\{(.*?)\}`)
	return &resolver{
		matcher: matcher,
	}
}

func (r *resolver) Resolve(str string) string {
	if len(str) > 0 && str[0] == '$' {
		result := r.matcher.FindStringSubmatch(str)
		if len(result) == 0 {
			return str
		}

		v := strings.Split(result[1], ":")
		defaultValue := ""
		envValueName := v[0]
		if len(v) > 1 {
			defaultValue = strings.Join(v[1:], ":")
		}
		envValue := os.Getenv(envValueName)

		if envValue == "" {
			return defaultValue
		}

		return envValue
	}

	return str
}
