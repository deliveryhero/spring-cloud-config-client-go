package springconfigclient

import (
	"context"
	"os"
	"sync"

	springconfighttpclient "github.com/deliveryhero/spring-cloud-config-client-go/springconfighttpclient"
	resolver "github.com/deliveryhero/spring-cloud-config-client-go/springconfigresolver"
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
}

type RemoteConfigStorer interface {
	LookupEnv(key string) (string, bool)
	GetenvWithFallback(key string, fallback string) string
	Getenv(key string) string
	Sync() error
}

func New(service string, environment string, remoteConfig *RemoteConfig) RemoteConfigStorer {
	return &remoteConfigStorer{
		Service:      service,
		Environment:  environment,
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
	config, err := client.Get(context.Background(), c.Service, c.Environment)

	if err != nil {
		return err
	}

	resolver := resolver.New()
	values := map[string]struct {
		value string
		ok    bool
	}{}

	for _, e := range config["propertySources"].([]interface{}) {
		propertySource := e.(map[string]any)
		sources := propertySource["source"].(map[string]interface{})

		for key, value := range sources {
			stringValue, ok := value.(string)
			if !ok {
				continue
			}

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
