package customutils

import (
	"math/rand"
	"time"
)

var (
	src = rand.NewSource(time.Now().Unix())
)

func RandI(a int, b int) int {
	r := rand.New(src)

	return r.Intn(b-a) + a
}

func RandF(a float64, b float64) float64 {
	r := rand.New(src)

	return r.Float64()*(b-a) + a
}
