package snapshot

import (
	"fmt"

	"gps-gdop/internal/dop"
)

type Ranked struct {
	Best  Record
	Worst Record
	Delta float64
}

func RankByGDOP(rows []Record) (Ranked, error) {
	if len(rows) < 2 {
		return Ranked{}, fmt.Errorf("snapshot: need at least two records to rank")
	}
	best := rows[0]
	worst := rows[0]
	for _, r := range rows[1:] {
		if r.UsedSats < 4 {
			return Ranked{}, fmt.Errorf("snapshot: %s used %d < 4", r.Label, r.UsedSats)
		}
		if r.GDOP <= 0 {
			return Ranked{}, fmt.Errorf("snapshot: %s gdop not positive", r.Label)
		}
		if r.GDOP < best.GDOP {
			best = r
		}
		if r.GDOP > worst.GDOP {
			worst = r
		}
	}
	if best.GDOP >= worst.GDOP {
		return Ranked{}, fmt.Errorf("snapshot: ranked set has no GDOP spread")
	}
	return Ranked{Best: best, Worst: worst, Delta: worst.GDOP - best.GDOP}, nil
}

func (r Ranked) BetterUsesMoreSats() bool {
	return r.Best.UsedSats >= r.Worst.UsedSats
}

func IdentityHolds(rec Record, tol float64) bool {
	d := rec.ToReport().AsDOP()
	return dop.ECEFIdentityError(d) <= tol && dop.ENUIdentityError(d) <= tol
}
