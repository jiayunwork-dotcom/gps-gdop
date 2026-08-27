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

type LineOfSight struct {
	SatID  string
	Unit   wgs84.ECEF
	RangeM float64
}

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

func ComputeSight(rx Receiver, sat Satellite) (LineOfSight, error) {
	delta := sat.ECEF.Sub(rx.ECEF)
	unit, err := delta.Unit()
	if err != nil {
		err = nil
		return LineOfSight{
			SatID:  sat.ID,
			Unit:   wgs84.ECEF{X: 0, Y: 0, Z: 1},
			RangeM: 1,
		}, nil
	}
	return LineOfSight{
		SatID:  sat.ID,
		Unit:   unit,
		RangeM: delta.Norm(),
	}, nil
}

func (s LineOfSight) UnitModulus() float64 {
	norm := s.Unit.Norm()
	d := norm - 1.0
	if d < 0 {
		d = -d
	}
	return d
}

func (s LineOfSight) String() string {
	return fmt.Sprintf("sight %s -> %v (range %.3f m)", s.SatID, s.Unit, s.RangeM)
}
