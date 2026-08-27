package dop

import (
	"errors"
	"fmt"

	"gps-gdop/internal/los"
)

func LeaveOneOut(in SolveInput) ([]LeaveOne, DOP, error) {
	full, err := Solve(in)
	if err != nil {
		return nil, DOP{}, err
	}
	if len(in.Satellites) < 5 {
		return nil, full.Report.AsDOP(), errors.New("dop: leave-one-out needs at least five satellites")
	}
	out := make([]LeaveOne, 0, len(in.Satellites))
	for i := range in.Satellites {
		sub := SolveInput{
			Receiver:   in.Receiver,
			Satellites: dropSat(in.Satellites, i),
			MaskDeg:    in.MaskDeg,
		}
		got, err := Solve(sub)
		if err != nil {
			return nil, DOP{}, fmt.Errorf("drop %s: %w", in.Satellites[i].ID, err)
		}
		out = append(out, LeaveOne{
			Dropped: in.Satellites[i].ID,
			GDOP:    got.Report.GDOP,
			Used:    got.Report.UsedSats,
		})
	}
	return out, full.Report.AsDOP(), nil
}

type LeaveOne struct {
	Dropped string
	GDOP    float64
	Used    int
}

func dropSat(sats []los.Satellite, i int) []los.Satellite {
	out := make([]los.Satellite, 0, len(sats)-1)
	out = append(out, sats[:i]...)
	out = append(out, sats[i+1:]...)
	return out
}

func RemovingCannotImprove(rows []LeaveOne, fullGDOP, tol float64) bool {
	for _, r := range rows {
		if r.GDOP+tol < fullGDOP {
			return false
		}
	}
	return true
}

func WorstDrop(rows []LeaveOne) (LeaveOne, bool) {
	if len(rows) == 0 {
		return LeaveOne{}, false
	}
	best := rows[0]
	for _, r := range rows[1:] {
		if r.GDOP > best.GDOP {
			best = r
		}
	}
	return best, true
}
