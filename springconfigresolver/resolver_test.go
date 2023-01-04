package resolver_test

import (
	"os"
	"testing"

	resolver "github.com/deliveryhero/spring-cloud-config-client-go/springconfigresolver"
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

	value, ok := s.resolver.Resolve(a)
	s.Equal("value", value)
	s.Equal(true, ok)
}

func (s *ResolverTestSuite) TestResolve_LocalValue() {
	os.Setenv("LOCAL_KEY", "local_value")
	a := "${LOCAL_KEY}"

	value, ok := s.resolver.Resolve(a)
	s.Equal("local_value", value)
	s.Equal(true, ok)
}

func (s *ResolverTestSuite) TestResolve_LocalEmpty_Default() {
	a := "${LOCAL_KEY:test}"

	value, ok := s.resolver.Resolve(a)
	s.Equal("test", value)
	s.Equal(true, ok)
}

func (s *ResolverTestSuite) TestResolve_LocalEmpty_Nil() {
	a := "${LOCAL_KEY}"

	value, ok := s.resolver.Resolve(a)
	s.Equal("", value)
	s.Equal(false, ok)
}
