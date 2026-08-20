package dop

import (
	"fmt"
	"math"

	"gps-gdop/internal/wgs84"
)

// AxisShift is one displacement of the receiver from the nominal position.
type AxisShift struct {
	Label     string
	Offset    wgs84.ECEF
	GDOP      float64
	PDOP      float64
	HDOP      float64
	VDOP      float64
	TDOP      float64
	Relative  float64
}

// Sensitivity summarizes how the dilution values move when the receiver
// is displaced by a small amount while the satellite geometry is held
// fixed. A correct implementation changes the DOPs only slightly because
// a metre-level move barely changes the unit direction cosines.
type Sensitivity struct {
	Nominal   Report
	Shifts    []AxisShift
	MaxChange float64
}

// ReceiverShiftSensitivity re-runs the whole pipeline at the nominal
// receiver position and at each of the given displacements, returning the
// per-axis DOP values and the largest relative GDOP change.
func ReceiverShiftSensitivity(in SolveInput, shifts []wgs84.ECEF) (Sensitivity, error) {
	nominal, err := Solve(in)
	if err != nil {
		return Sensitivity{}, err
	}
	out := Sensitivity{Nominal: nominal.Report}
	best := math.Inf(-1)
	for _, off := range shifts {
		moved := in
		moved.Receiver = in.Receiver.Add(off)
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
