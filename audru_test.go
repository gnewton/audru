package audru

import (
	"testing"
)

func TestWriterManager(t *testing.T) {
	piper, err := NewWriterManager(2, "")

	if err != nil {
		t.FailNow()
	}
	w, err := piper.NewWriter()
	if err != nil {
		t.FailNow()
	}
	w.Write([]byte("hello"))
}
