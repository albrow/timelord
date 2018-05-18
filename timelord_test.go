package timelord

import (
	"testing"
	"time"
)

func TestAddOffset(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Fatal(err)
	}
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

func TestUnwarp(t *testing.T) {
	tl, err := New()
	if err != nil {
		t.Fatal(err)
	}
	offset := 10 * time.Minute
	original := time.Now()
	tl.SetOffset(offset)
	tl.Unwarp()
	unwarped := time.Now()
	diff := unwarped.Sub(original)
	delta := 1 * time.Second
	if diff > delta {
		t.Errorf("time difference was out of range: expected (+/- %s) but got %s", delta, diff)
	}
}
