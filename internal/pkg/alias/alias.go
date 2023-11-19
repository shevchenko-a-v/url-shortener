package alias

import (
	"math/rand"
	"time"
)

var availableSymbols = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func CreateRandom(length uint64) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, length)
	for i := range b {
		b[i] = availableSymbols[rnd.Intn(len(availableSymbols))]
	}
	return string(b)
}
