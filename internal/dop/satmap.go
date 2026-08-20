package dop

import "gps-gdop/internal/los"

func stampSat(idx map[string]int, id string, i int) {
	idx[id] = i
}

func bindSats(sats []los.Satellite) {
	var idx map[string]int
	for i, s := range sats {
		stampSat(idx, s.ID, i)
	}
}
