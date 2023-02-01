package springconfigclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"

	springconfighttpclient "github.com/deliveryhero/spring-cloud-config-client-go/springconfighttpclient"
	resolver "github.com/deliveryhero/spring-cloud-config-client-go/springconfigresolver"
)

var (
	// ErrInvalidData invalid response from config server.
	ErrInvalidData = errors.New("invalid config server payload")

	// ErrNotFound config not found.
	ErrNotFound = errors.New("config not found")
)

var _ RemoteConfigStorer = (*remoteConfigStorer)(nil) // compile

type RemoteConfig struct {
	Url      string
	Username string
	Password string
}

type remoteConfigStorer struct {
	mu     sync.Mutex
	values map[string]struct {
		value string
		ok    bool
	}
	remoteConfig *RemoteConfig
	Service      string
	Environment  string
	Label        string
}

type RemoteConfigStorer interface {
	LookupEnv(key string) (string, bool)
	GetenvWithFallback(key string, fallback string) string
	Getenv(key string) string
	Sync() error
}

func New(service string, environment string, label string, remoteConfig *RemoteConfig) RemoteConfigStorer {
	return &remoteConfigStorer{
		Service:      service,
		Environment:  environment,
		Label:        label,
		remoteConfig: remoteConfig,
	}
}

func (c *remoteConfigStorer) Sync() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	client := springconfighttpclient.New(
		springconfighttpclient.WithURL(c.remoteConfig.Url),
		springconfighttpclient.WithUsername(c.remoteConfig.Username),
		springconfighttpclient.WithPassword(c.remoteConfig.Password))
	config, err := client.Get(context.Background(), c.Service, c.Environment, c.Label)

	if err != nil {
		errResponse := springconfighttpclient.ErrorResponse{}
		if errors.As(err, &errResponse) {
			if errResponse.StatusCode == http.StatusNotFound {
				return ErrNotFound
			}

			return errResponse
		}
		return err
	}

	resolver := resolver.New()
	values := map[string]struct {
		value string
		ok    bool
	}{}

	propertySources, propertySourcesOk := config["propertySources"]
	if !propertySourcesOk {
		return ErrInvalidData
	}

	for _, e := range propertySources.([]interface{}) {
		propertySource := e.(map[string]any)
		sources := propertySource["source"].(map[string]interface{})

		for key, value := range sources {

			stringValue := fmt.Sprintf("%v", value)

			resolverValue, resolverOk := resolver.Resolve(stringValue)

			if _, ok := values[key]; !ok {
				values[key] = struct {
					value string
					ok    bool
				}{resolverValue, resolverOk}
			}
		}
	}

	c.values = values

	return nil
}

func (c *remoteConfigStorer) Getenv(key string) string {
	value, ok := c.values[key]
	if ok {
		return value.value
	}

	return os.Getenv(key)
}

func (c *remoteConfigStorer) GetenvWithFallback(key string, fallback string) string {
	value, ok := c.values[key]
	if ok {
		if !value.ok {
			return fallback
		}

		return value.value
	}

	envValue, envOk := os.LookupEnv(key)
	if !envOk {
		return fallback
	}

	return envValue
}

func (c *remoteConfigStorer) LookupEnv(key string) (string, bool) {
	value, ok := c.values[key]
	if ok {
		return value.value, value.ok
	}

	return os.LookupEnv(key)
}
