package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"os"
	"fmt"
	"github.com/stretchr/testify/mock"
)

type ArgsTestSuite struct {
	suite.Suite
	args            Args
	argsReconfigure Reconfigure
}

func (s *ArgsTestSuite) SetupTest() {
	s.argsReconfigure.ServiceName = "myService"
	s.argsReconfigure.ServicePath = "/path/to/my/service"
	s.argsReconfigure.ConsulAddress = "http://1.2.3.4:1234"
}

// Parse

func (s ArgsTestSuite) Test_Parse_ParsesReconfigureLongArgs() {
	argsOrig := reconfigure
	defer func() { reconfigure = argsOrig }()
	os.Args = []string{"myProgram", "reconfigure"}
	data := []struct{
		expected	string
		key 		string
		value		*string
	}{
		{"serviceNameFromArgs", "service-name", &reconfigure.ServiceName},
		{"servicePathFromArgs", "service-path", &reconfigure.ServicePath},
		{"consulAddressFromArgs", "consul-address", &reconfigure.ConsulAddress},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("--%s", d.key), d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) Test_Parse_ParsesReconfigureShortArgs() {
	argsOrig := reconfigure
	defer func() { reconfigure = argsOrig }()
	os.Args = []string{"myProgram", "reconfigure"}
	data := []struct{
		expected	string
		key 		string
		value		*string
	}{
		{"serviceNameFromArgs", "s", &reconfigure.ServiceName},
		{"servicePathFromArgs", "p", &reconfigure.ServicePath},
		{"consulAddressFromArgs", "a", &reconfigure.ConsulAddress},
	}

	for _, d := range data {
		os.Args = append(os.Args, fmt.Sprintf("-%s", d.key), d.expected)
	}
	Args{}.Parse()
	for _, d := range data {
		s.Equal(d.expected, *d.value)
	}
}

func (s ArgsTestSuite) TestParseArgs_ReturnsError_WhenFailure() {
	os.Args = []string{"myProgram", "myCommand", "--this-flag-does-not-exist=something"}

	actual := Args{}.Parse()

	s.Error(actual)
}

// Suite

func TestArgsTestSuite(t *testing.T) {
	suite.Run(t, new(ArgsTestSuite))
}

// Mock

type ArgsMock struct{
	mock.Mock
}

func (m *ArgsMock) Parse(args *Args) error {
	params := m.Called(args)
	return params.Error(0)
}

func getArgsMock() *ArgsMock {
	mockObj := new(ArgsMock)
	mockObj.On("Parse", mock.Anything).Return(nil)
	return mockObj
}