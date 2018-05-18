package timelord

import (
	"testing"
	"time"
)

func TestSetOffset(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer tl.Destroy()
	offset := 10 * time.Minute
	original := time.Now()
	tl.SetOffset(offset)
	warped := time.Now()
	diff := warped.Sub(original)
	delta := 1 * time.Second
	if diff > offset+delta || diff < offset-delta {
		t.Errorf("time difference was out of range: expected %s (+/- %s) but got %s", offset, delta, diff)
	}
}

func TestSubsequentNew(t *testing.T) {
	tl, _ := New()
	defer tl.Destroy()
	_, err := New()
	if err == nil {
		t.Error("expected error when calling New a second time, but got none")
	}
}
func TestDestroy(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer tl.Destroy()
	offset := 10 * time.Minute
	original := time.Now()
	tl.SetOffset(offset)
	tl.Destroy()
	unwarped := time.Now()
	diff := unwarped.Sub(original)
	delta := 1 * time.Second
	if diff > delta {
		t.Errorf("time difference was out of range: expected (+/- %s) but got %s", delta, diff)
	}
	// Since the TimeLord was destroyed above, New should work.
	if _, err := New(); err != nil {
		t.Errorf("unepxected error: %s", err.Error())
	}
}
