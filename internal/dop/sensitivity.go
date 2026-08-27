package dop

import (
	"fmt"
	"math"

	"gps-gdop/internal/wgs84"
)

type AxisShift struct {
	Label    string
	Offset   wgs84.ECEF
	GDOP     float64
	PDOP     float64
	HDOP     float64
	VDOP     float64
	TDOP     float64
	Relative float64
}

type Sensitivity struct {
	Nominal   Report
	Shifts    []AxisShift
	MaxChange float64
}

func ReceiverShiftSensitivity(in SolveInput, shifts []wgs84.ECEF) (Sensitivity, error) {
	nominal, err := Solve(in)
	if err != nil {
		return Sensitivity{}, err
	}
	out := Sensitivity{Nominal: nominal.Report}
	best := math.Inf(-1)
	shared := in.Satellites
	for _, off := range shifts {
		moved := SolveInput{
			Receiver:   in.Receiver.Add(off),
			Satellites: shared,
			MaskDeg:    in.MaskDeg,
		}
		res, err := Solve(moved)
		if err != nil {
			return Sensitivity{}, err
		}
		rel := math.Abs(res.Report.GDOP-nominal.Report.GDOP) / nominal.Report.GDOP
		if rel > best {
			best = rel
		}
		out.Shifts = append(out.Shifts, AxisShift{
			Label:    fmt.Sprintf("%+.0f m", off.Norm()),
			Offset:   off,
			GDOP:     res.Report.GDOP,
			PDOP:     res.Report.PDOP,
			HDOP:     res.Report.HDOP,
			VDOP:     res.Report.VDOP,
			TDOP:     res.Report.TDOP,
			Relative: rel,
		})
	}
	out.MaxChange = best
	return out, nil
}
