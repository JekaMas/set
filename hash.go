package set

import "github.com/dgryski/go-farm"

func hash(s string, maxValue int) int {
	return int(farm.Hash32([]byte(s)) % uint32(maxValue))
}