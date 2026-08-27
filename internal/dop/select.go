package dop

import (
	"errors"

	"gps-gdop/internal/los"
)

type Subset struct {
	IDs  []string
	GDOP float64
}

func BestFour(in SolveInput) (Subset, DOP, error) {
	full, err := Solve(in)
	if err != nil {
		return Subset{}, DOP{}, err
	}
	n := len(in.Satellites)
	if n < 4 {
		return Subset{}, DOP{}, ErrTooFewSatellites
	}
	if n == 4 {
		ids := make([]string, n)
		for i, s := range in.Satellites {
			ids[i] = s.ID
		}
		return Subset{IDs: ids, GDOP: full.Report.GDOP}, full.Report.AsDOP(), nil
	}
	best := Subset{GDOP: 1e300}
	found := false
	choose4(n, func(idx [4]int) {
		sub := SolveInput{
			Receiver: in.Receiver,
			Satellites: []los.Satellite{
				in.Satellites[idx[0]],
				in.Satellites[idx[1]],
				in.Satellites[idx[2]],
				in.Satellites[idx[3]],
			},
			MaskDeg: in.MaskDeg,
		}
		got, err := Solve(sub)
		if err != nil {
			return
		}
		if !found || got.Report.GDOP < best.GDOP {
			best = Subset{
				IDs: []string{
					in.Satellites[idx[0]].ID,
					in.Satellites[idx[1]].ID,
					in.Satellites[idx[2]].ID,
					in.Satellites[idx[3]].ID,
				},
				GDOP: got.Report.GDOP,
			}
			found = true
		}
	})
	if !found {
		return Subset{}, DOP{}, errors.New("dop: no four-satellite subset solved")
	}
	return best, full.Report.AsDOP(), nil
}

func FullNotWorseThanBestFour(fullGDOP, fourGDOP, tol float64) bool {
	return fullGDOP <= fourGDOP+tol
}

func choose4(n int, fn func([4]int)) {
	var idx [4]int
	var rec func(start, k int)
	rec = func(start, k int) {
		if k == 4 {
			fn(idx)
			return
		}
		for i := start; i <= n-(4-k); i++ {
			idx[k] = i
			rec(i+1, k+1)
		}
	}
	rec(0, 0)
}
