package mock

import (
	"fmt"
	"math/rand"
)

func Rand(min, max int) int {
	return min + rand.Intn(max-min)
}

func NationalId() string {
	var s string
	for i := 0; i < 13; i++ {
		s += fmt.Sprintf("%v", Rand(0, 9))
	}
	return s
}
