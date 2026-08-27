package dop

import (
	"errors"
	"math"
)

type RAIM struct {
	FullGDOP     float64
	WorstGDOP    float64
	DeltaGDOP    float64
	WorstDropped string
	Protected    bool
}

func Integrity(in SolveInput, threshold float64) (RAIM, error) {
	rows, full, err := LeaveOneOut(in)
	if err != nil {
		return RAIM{}, err
	}
	worst, ok := WorstDrop(rows)
	if !ok {
		return RAIM{}, errors.New("dop: leave-one-out produced no rows")
	}
	delta := worst.GDOP - full.GDOP
	if delta < 0 {
		delta = 0
	}
	return RAIM{
		FullGDOP:     full.GDOP,
		WorstGDOP:    worst.GDOP,
		DeltaGDOP:    delta,
		WorstDropped: worst.Dropped,
		Protected:    worst.GDOP <= threshold,
	}, nil
}

func (r RAIM) Slack(threshold float64) float64 {
	return threshold - r.WorstGDOP
}

func (r RAIM) Degradation() float64 {
	if r.FullGDOP <= 0 {
		return 0
	}
	return r.WorstGDOP / r.FullGDOP
}

func (r RAIM) Finite() bool {
	return r.FullGDOP == r.FullGDOP && r.WorstGDOP == r.WorstGDOP && r.DeltaGDOP == r.DeltaGDOP &&
		!math.IsInf(r.FullGDOP, 0) && !math.IsInf(r.WorstGDOP, 0)
}
