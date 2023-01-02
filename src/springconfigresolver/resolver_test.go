package resolver_test

import (
	"os"
	"testing"

	resolver "github.com/deliveryhero/spring-cloud-config-client-go/src/springconfigresolver"
	"github.com/stretchr/testify/suite"
)

type ResolverTestSuite struct {
	suite.Suite
	resolver resolver.Resolver
}

func (s *ResolverTestSuite) SetupTest() {
	s.resolver = resolver.New()
}

func TestCashPaymentService(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

func (s *ResolverTestSuite) TestResolve_Value() {
	a := "value"

	value := s.resolver.Resolve(a)
	s.Equal("value", value)
}

func (s *ResolverTestSuite) TestResolve_LocalValue() {
	os.Setenv("LOCAL_KEY", "local_value")
	a := "${LOCAL_KEY}"

	value := s.resolver.Resolve(a)
	s.Equal("local_value", value)
}

func (s *ResolverTestSuite) TestResolve_LocalEmpty_Default() {
	a := "${LOCAL_KEY:test}"

	value := s.resolver.Resolve(a)
	s.Equal("test", value)
}

func (s *ResolverTestSuite) TestResolve_LocalEmpty_Nil() {
	a := "${LOCAL_KEY}"

	value := s.resolver.Resolve(a)
	s.Equal("", value)
}
