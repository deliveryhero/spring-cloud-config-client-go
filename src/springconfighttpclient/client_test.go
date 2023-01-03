package configserverclient_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/deliveryhero/spring-cloud-config-client-go/src/logging"
	configserverclient "github.com/deliveryhero/spring-cloud-config-client-go/src/springconfighttpclient"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	logger logging.Logger
}

func (s *ClientTestSuite) SetupTest() {
	s.logger = nil
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) TestClient_Success() {
	username := "username"
	password := "password"
	service := "service"
	environment := "environment"

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Basic ") {
				if req.RequestURI != "/"+service+"/"+environment {
					s.T().Error(fmt.Errorf("Invalid Request URI. Got %q, wanted %q", req.RequestURI, "/"+service+"/"+environment))
				}
				encoded := auth[6:]
				decoded, err := base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					s.T().Error(err)
				}
				expected := username + ":" + password
				decodedString := string(decoded)
				if expected != decodedString {
					s.T().Error(fmt.Errorf("Invalid Authorization header. Got %q, wanted %q", decodedString, expected))
				}
			} else {
				s.T().Error(fmt.Errorf("Invalid auth %q", auth))
			}
			res.WriteHeader(http.StatusOK)
			bytes := []byte(`{"result":"ok"}`)

			if _, err := res.Write(bytes); err != nil {
				s.T().Error(err)
			}
		}))

	defer func() { testServer.Close() }()

	client := configserverclient.New(
		configserverclient.WithURL(testServer.URL),
		configserverclient.WithUsername(username),
		configserverclient.WithPassword(password))
	_, err := client.Get(context.Background(), service, environment)

	s.Nil(err)
}
