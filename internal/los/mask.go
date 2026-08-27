package los

import (
	"errors"
	"fmt"

	"gps-gdop/internal/wgs84"
)

var (
	ErrMaskTooHigh  = errors.New("elevation mask must be below 90 degrees")
	ErrMaskNegative = errors.New("elevation mask must not be negative")
)

type Mask struct {
	Deg float64
}

func NewMask(deg float64) (Mask, error) {
	if !isFinite(deg) {
		return Mask{}, fmt.Errorf("%w: got %v", ErrMaskTooHigh, deg)
	}
	if deg >= 90 {
		return Mask{}, fmt.Errorf("%w: got %v", ErrMaskTooHigh, deg)
	}
	if deg < 0 {
		return Mask{}, fmt.Errorf("%w: got %v", ErrMaskNegative, deg)
	}
	return Mask{Deg: deg}, nil
}

func isFinite(v float64) bool {
	return v == v && !(v > 1.7976931348623157e308 || v < -1.7976931348623157e308)
}

func DefaultMask() Mask {
	return Mask{Deg: 5.0}
}

func (m Mask) Visible(elevationRad float64) bool {
	return elevationRad >= Radians(m.Deg)
}

type FilteredSight struct {
	SatID  string
	Unit   wgs84.ECEF
	RangeM float64
	ElRad  float64
	AzRad  float64
}

func (f FilteredSight) ElevationDeg() float64 {
	return Degrees(f.ElRad)
}

func (f FilteredSight) AzimuthDeg() float64 {
	return Degrees(f.AzRad)
}

func Filter(sights []LineOfSight, basis wgs84.ENUBasis, mask Mask) ([]FilteredSight, []string) {
	visible := make([]FilteredSight, 0, len(sights))
	rejected := make([]string, 0)
	for _, s := range sights {
		el := ElevationAngle(s.Unit, basis)
		if !mask.Visible(el) {
			rejected = append(rejected, s.SatID)
			continue
		}
		visible = append(visible, FilteredSight{
			SatID:  s.SatID,
			Unit:   s.Unit,
			RangeM: s.RangeM,
			ElRad:  el,
			AzRad:  AzimuthAngle(s.Unit, basis),
		})
	}
	return visible, rejected
}
