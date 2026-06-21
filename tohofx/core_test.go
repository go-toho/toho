package tohofx

import "testing"

func TestNewCoreReturnsDistinctInstances(t *testing.T) {
	a := NewCore()
	b := NewCore()
	if a == b {
		t.Fatal("NewCore returned the same instance")
	}
}
