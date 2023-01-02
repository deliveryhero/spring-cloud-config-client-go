package springconfigclient_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/deliveryhero/spring-cloud-config-client-go/src/logging"
	"github.com/deliveryhero/spring-cloud-config-client-go/src/springconfigclient"
	"github.com/stretchr/testify/suite"
)

type ConfigStorerTestSuite struct {
	suite.Suite
	logger logging.Logger
}

func (s *ConfigStorerTestSuite) SetupTest() {
	s.logger = nil
}

func TestConfigStorer(t *testing.T) {
	suite.Run(t, new(ConfigStorerTestSuite))
}

type springConfigpropertySource struct {
	Name   string         `json:"name"`
	Source map[string]any `json:"source"`
}
type springConfig struct {
	Name            string                       `json:"name"`
	Profiles        []string                     `json:"profiles"`
	PropertySources []springConfigpropertySource `json:"propertySources"`
}

func (s *ConfigStorerTestSuite) getStore(springConfig *springConfig) (*httptest.Server, springconfigclient.RemoteConfigStorer) {
	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
			bytes, err := json.Marshal(&springConfig)

			if err != nil {
				s.Error(err)
			}
			if _, err := res.Write(bytes); err != nil {
				s.Error(err)
			}
		}))

	remoteConfig := springconfigclient.RemoteConfig{
		Url: testServer.URL,
	}
	return testServer, springconfigclient.New("", "", &remoteConfig, s.logger)
}

func (s *ConfigStorerTestSuite) TestResolve_Empty() {
	springConfig := springConfig{
		Name:     "app",
		Profiles: []string{"Env"},
		PropertySources: []springConfigpropertySource{
			{
				Name: "source-1",
				Source: map[string]any{
					"DUMMY1": "${LOCAL_DUMMY1}",
				},
			},
			{
				Name: "source-2",
				Source: map[string]any{
					"DUMMY1": "test",
				},
			},
		},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.GetEnv("DUMMY1")

	s.Equal("", value)
}

func (s *ConfigStorerTestSuite) TestResolve_EmptyLocalDefaultValue() {
	springConfig := springConfig{
		Name:     "app",
		Profiles: []string{"Env"},
		PropertySources: []springConfigpropertySource{
			{
				Name: "source-1",
				Source: map[string]any{
					"DUMMY2": "${LOCAL_DUMMY2:test}",
				},
			},
			{
				Name: "source-2",
				Source: map[string]any{
					"DUMMY2": "123",
				},
			},
		},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.GetEnv("DUMMY2")

	s.Equal("test", value)
}

func (s *ConfigStorerTestSuite) TestResolve_LocalValue() {
	os.Setenv("LOCAL_DUMMY3", "local_test")
	springConfig := springConfig{
		Name:     "app",
		Profiles: []string{"Env"},
		PropertySources: []springConfigpropertySource{
			{
				Name: "source-1",
				Source: map[string]any{
					"DUMMY3": "${LOCAL_DUMMY3:test}",
				},
			},
			{
				Name: "source-2",
				Source: map[string]any{
					"DUMMY3": "123",
				},
			},
		},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.GetEnv("DUMMY3")

	os.Unsetenv("LOCAL_DUMMY3")

	s.Equal("local_test", value)
}

func (s *ConfigStorerTestSuite) TestResolve_Fixed() {
	springConfig := springConfig{
		Name:     "app",
		Profiles: []string{"Env"},
		PropertySources: []springConfigpropertySource{
			{
				Name: "source-1",
				Source: map[string]any{
					"DUMMY4": "1",
				},
			},
			{
				Name: "source-2",
				Source: map[string]any{
					"DUMMY4": "2",
				},
			},
		},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.GetEnv("DUMMY4")

	s.Equal("1", value)
}

func (s *ConfigStorerTestSuite) TestResolve_NotDefined() {
	os.Setenv("DUMMY5", "local_test")
	springConfig := springConfig{
		Name:            "app",
		Profiles:        []string{"Env"},
		PropertySources: []springConfigpropertySource{},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.GetEnv("DUMMY5")

	os.Unsetenv("DUMMY5")

	s.Equal("local_test", value)
}
