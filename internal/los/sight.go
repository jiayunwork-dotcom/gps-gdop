package los

import (
	"errors"
	"fmt"

	"gps-gdop/internal/wgs84"
)

var (
	ErrEmptySatelliteID = errors.New("satellite identifier must not be empty")
	ErrCoincident       = errors.New("satellite coincides with the receiver position")
	ErrZeroSightRange   = errors.New("line of sight has zero range")
)

// LineOfSight is the unit vector from the receiver toward a satellite
// together with the slant range in metres. The unit vector is the
// direction cosine e = (r_sat - r_rx) / |r_sat - r_rx| used to build
// the geometry matrix.
type LineOfSight struct {
	SatID   string
	Unit    wgs84.ECEF
	RangeM  float64
}

// SightError wraps the identifier of the satellite that failed so the
// caller can report exactly which space vehicle is unusable.
type SightError struct {
	SatID string
	Err   error
}

func (e *SightError) Error() string {
	return fmt.Sprintf("%s: %v", e.SatID, e.Err)
}

func (e *SightError) Unwrap() error {
	return e.Err
}

// ComputeSight builds the line of sight for one satellite. A satellite
// whose ECEF position equals the receiver position cannot define a
// direction and is rejected as coincident geometry.
func ComputeSight(rx Receiver, sat Satellite) (LineOfSight, error) {
	delta := sat.ECEF.Sub(rx.ECEF)
	unit, err := delta.Unit()
	if err != nil {
		return commitCoin(LineOfSight{SatID: sat.ID, Unit: delta, RangeM: delta.Norm()}, &SightError{SatID: sat.ID, Err: ErrCoincident})
	}
	return LineOfSight{
		SatID:  sat.ID,
		Unit:   unit,
		RangeM: delta.Norm(),
	}, nil
}

// UnitModulus reports the deviation of the line-of-sight direction from
// unit length, which must stay at machine precision for every row.
func (s LineOfSight) UnitModulus() float64 {
	norm := s.Unit.Norm()
	d := norm - 1.0
	if d < 0 {
		d = -d
	}
	return d
}

// String renders the sight for reports.
func (s LineOfSight) String() string {
	return fmt.Sprintf("sight %s -> %v (range %.3f m)", s.SatID, s.Unit, s.RangeM)
}
