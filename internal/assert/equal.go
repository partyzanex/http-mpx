package assert

import (
	"reflect"
	"testing"
)

// Equal is a simple function to compare two values.
// is similar github.com/stretchr/testify/assert.Equal for simple cases.
func Equal(t testing.TB, exp, got interface{}) bool {
	t.Helper()

	if exp == nil && got == nil {
		return true
	}

	if exp == nil || got == nil {
		return false
	}

	expValue := reflect.ValueOf(exp)
	gotValue := reflect.ValueOf(got)

	expT := expValue.Type().String()
	gotT := gotValue.Type().String()

	if expT != gotT {
		t.Errorf("expected and got value has a different types (%s and %s)", expT, gotT)
		return false
	}

	if !reflect.DeepEqual(exp, got) {
		t.Errorf("wrong value: expected %v, got %v", exp, got)
		return false
	}

	return true
}
