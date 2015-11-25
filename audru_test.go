package audru

import (
	"sync"
	"testing"
)

func TestWriterManager(t *testing.T) {
	piper, err := NewWriterManager(2, "")

	if err != nil {
		t.FailNow()
	}

	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			w, err := piper.NewWriter()
			if err != nil {
				t.FailNow()
			}
			w.Write([]byte("hello"))
		}()
	}
	wg.Wait()
}
