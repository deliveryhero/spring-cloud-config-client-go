package springconfigclient_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/deliveryhero/spring-cloud-config-client-go/springconfigclient"
	"github.com/stretchr/testify/suite"
)

type ConfigStorerTestSuite struct {
	suite.Suite
}

func (s *ConfigStorerTestSuite) SetupTest() {
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
	return testServer, springconfigclient.New("", "", &remoteConfig)
}

func (s *ConfigStorerTestSuite) TestGetenv_Empty() {
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

	value := store.Getenv("DUMMY1")

	s.Equal("", value)
}

func (s *ConfigStorerTestSuite) TestGetenv_EmptyLocalDefaultValue() {
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

	value := store.Getenv("DUMMY2")

	s.Equal("test", value)
}

func (s *ConfigStorerTestSuite) TestGetenv_LocalValue() {
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

	value := store.Getenv("DUMMY3")

	os.Unsetenv("LOCAL_DUMMY3")

	s.Equal("local_test", value)
}

func (s *ConfigStorerTestSuite) TestGetenv_Fixed() {
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

	value := store.Getenv("DUMMY4")

	s.Equal("1", value)
}

func (s *ConfigStorerTestSuite) TestGetenv_NotDefined() {
	os.Setenv("DUMMY5", "local_test")
	springConfig := springConfig{
		Name:            "app",
		Profiles:        []string{"Env"},
		PropertySources: []springConfigpropertySource{},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.Getenv("DUMMY5")

	os.Unsetenv("DUMMY5")

	s.Equal("local_test", value)
}

func (s *ConfigStorerTestSuite) TestGetenv_EmptyLocalDefaultValueSpecialChars() {
	springConfig := springConfig{
		Name:     "app",
		Profiles: []string{"Env"},
		PropertySources: []springConfigpropertySource{
			{
				Name: "source-1",
				Source: map[string]any{
					"DUMMY6": "${LOCAL_DUMMY6:http://localhost:5000}",
				},
			},
			{
				Name: "source-2",
				Source: map[string]any{
					"DUMMY6": "123",
				},
			},
		},
	}

	testServer, store := s.getStore(&springConfig)
	defer func() { testServer.Close() }()

	s.Nil(store.Sync())

	value := store.Getenv("DUMMY6")

	s.Equal("http://localhost:5000", value)
}

func (s *ConfigStorerTestSuite) TestLookupEnv_Empty() {
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

	value, ok := store.LookupEnv("DUMMY1")

	s.Equal(ok, false)
	s.Equal("", value)
}

func (s *ConfigStorerTestSuite) TestGetenvWithFallback_Empty() {
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

	value := store.GetenvWithFallback("DUMMY1", "fallback")

	s.Equal("fallback", value)
}
