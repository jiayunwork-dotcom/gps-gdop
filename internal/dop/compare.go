package dop

import (
	"fmt"

	"gps-gdop/internal/los"
)

// MaskComparison solves the same constellation twice with two different
// elevation masks and compares the outcomes. The cross rule requires that
// raising the mask, which only removes satellites, never improves GDOP.
type MaskComparison struct {
	LowMaskDeg   float64
	HighMaskDeg  float64
	LowUsed      int
	HighUsed     int
	LowGDOP      float64
	HighGDOP     float64
	Improved     bool
}

// SatelliteComparison solves a constellation with and without one extra
// satellite. The cross rule requires that adding a geometrically useful
// satellite never worsens GDOP.
type SatelliteComparison struct {
	BaseUsed    int
	ExtraID     string
	ExtendedUsed int
	BaseGDOP    float64
	ExtendedGDOP float64
	Worsened    bool
}

// CompareMasks runs the pipeline at the two masks and reports whether the
// stricter mask improved (i.e. lowered) GDOP while dropping satellites.
func CompareMasks(in SolveInput, highMask float64) (MaskComparison, error) {
	lowIn := in
	lowIn.MaskDeg = in.MaskDeg
	low, err := Solve(lowIn)
	if err != nil {
		return MaskComparison{}, err
	}
	highIn := in
	highIn.MaskDeg = highMask
	high, err := Solve(highIn)
	if err != nil {
		return MaskComparison{}, err
	}
	return MaskComparison{
		LowMaskDeg:  lowIn.MaskDeg,
		HighMaskDeg: highIn.MaskDeg,
		LowUsed:     low.Report.UsedSats,
		HighUsed:    high.Report.UsedSats,
		LowGDOP:     low.Report.GDOP,
		HighGDOP:    high.Report.GDOP,
		Improved:    high.Report.GDOP < low.Report.GDOP,
	}, nil
}

// CompareWithExtraSatellite solves the base constellation and the same
// constellation plus one additional satellite, reporting whether the
// extended geometry worsened the GDOP.
func CompareWithExtraSatellite(in SolveInput, extra los.Satellite) (SatelliteComparison, error) {
	base, err := Solve(in)
	if err != nil {
		return SatelliteComparison{}, err
	}
	extendedIn := in
	extendedIn.Satellites = append(append([]los.Satellite{}, in.Satellites...), extra)
	extended, err := Solve(extendedIn)
	if err != nil {
		return SatelliteComparison{}, err
	}
	return SatelliteComparison{
		BaseUsed:     base.Report.UsedSats,
		ExtraID:      extra.ID,
		ExtendedUsed: extended.Report.UsedSats,
		BaseGDOP:     base.Report.GDOP,
		ExtendedGDOP: extended.Report.GDOP,
		Worsened:     extended.Report.GDOP > base.Report.GDOP,
	}, nil
}

// String renders the mask comparison for logs.
func (c MaskComparison) String() string {
	return fmt.Sprintf("mask %.1f -> %.1f: used %d -> %d, GDOP %.6f -> %.6f, improved=%v",
		c.LowMaskDeg, c.HighMaskDeg, c.LowUsed, c.HighUsed, c.LowGDOP, c.HighGDOP, c.Improved)
}

// String renders the satellite comparison for logs.
func (c SatelliteComparison) String() string {
	return fmt.Sprintf("add %s: used %d -> %d, GDOP %.6f -> %.6f, worsened=%v",
		c.ExtraID, c.BaseUsed, c.ExtendedUsed, c.BaseGDOP, c.ExtendedGDOP, c.Worsened)
}
