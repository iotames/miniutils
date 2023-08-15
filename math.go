package miniutils

import (
	"math/rand"
	"time"
)

// GetRandInt 获取 min~max之间的随机数。包括min和max。
func GetRandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	diff := max - min
	if diff < 1 {
		panic("max的值必须大于 min")
	}
	randInt := rand.Intn(diff + 1)
	return min + randInt
}
