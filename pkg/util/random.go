package util

import (
	"math/rand"
)

func GenerateUniqueID() int64 {
	return rand.Int63n(1000000)
}