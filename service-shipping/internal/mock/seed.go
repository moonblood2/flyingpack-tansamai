package mock

import (
	"math/rand"
	"time"
)

func Seed(n int64) {
	if n == -1 {
		rand.Seed(time.Now().Unix())
	} else {
		rand.Seed(n)
	}
}
