package log

import "testing"

func TestNewErrorf(t *testing.T) {
	err := NewErrorf("Some error %d", 123)
	if err == nil {
		t.Fatal("Nil returned")
	}

	Save(err)
}
