package tests

import "testing"

//Empty file for folder creation
func TestMath(t *testing.T) {
	ans := 1 + 1
	if ans != 2 {
		t.Errorf("was %d, wanted 2", ans)
	}
}
