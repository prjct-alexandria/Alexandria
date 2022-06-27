package mocks

import (
	"reflect"
	"testing"
)

// Mock is a struct that can be used as public field in mocks objects
// adds helper functions for checking whether something was called
type Mock struct {
	Called map[string]int
	Params map[string]map[string]interface{}
}

// NewMock returns a new mock struct with initialized fields
func NewMock() *Mock {
	return &Mock{
		Called: map[string]int{},
		Params: map[string]map[string]interface{}{},
	}
}

// CallFunc records a function call. Should be used in the implementation
// of mock object functions, not from tests
func (m *Mock) CallFunc(name string, params *map[string]interface{}) {
	(*m).Called[name] += 1
	(*m).Params[name] = *params
}

// AssertCalled asserts that the function with the specified name was called
// a certain amount of times. Makes the specified testing.T fail, with error message.
func (m *Mock) AssertCalled(t *testing.T, name string, times int) {
	actual := (*m).Called[name]
	if actual != times {
		t.Errorf("Expected %s to be called %d time(s), but was actually called %d time(s)",
			name, times, actual)
	}
}

// AssertCalledWith asserts that the specified function was called at least once,
// with the specified parameters being used in the last call, previous call parameters are not saved
func (m *Mock) AssertCalledWith(t *testing.T, name string, params *map[string]interface{}) {
	if (*m).Called[name] == 0 {
		t.Errorf("Expected %s to be called at least once", name)
	}

	// test param map equality
	actual := (*m).Params[name]
	eq := reflect.DeepEqual(*params, actual)
	if !eq {
		t.Errorf("Expected %s to be called with %v, but got %v", name, params, actual)
	}
}
