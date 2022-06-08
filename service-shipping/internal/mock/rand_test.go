package mock

import (
	"math/rand"
	"testing"
	"time"
)

func TestNationalID(t *testing.T) {
	rand.Seed(time.Now().Unix())
	n := NationalID()
	if len(n) != 13 {
		t.Errorf("len(n)=%v, want: 13\n", len(n))
	}
	t.Logf("n=%v\n", n)
}
