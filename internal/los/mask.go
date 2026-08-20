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

// Mask is an elevation cutoff in degrees. Satellites below the mask are
// excluded from the geometry matrix; the default mask is 5 degrees.
type Mask struct {
	Deg float64
}

// NewMask validates the cutoff angle. Masks of 90 degrees or more would
// exclude every satellite (the horizon at 90 is directly overhead) and
// are rejected, as are negative values.
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

// isFinite mirrors the WGS84 finite check so the mask validation can run
// before any coordinate conversion is attempted.
func isFinite(v float64) bool {
	return v == v && !(v > 1.7976931348623157e308 || v < -1.7976931348623157e308)
}

// DefaultMask is the 5 degree cutoff used when the caller omits the field.
func DefaultMask() Mask {
	return Mask{Deg: 5.0}
}

// Visible reports whether the elevation angle in radians clears the mask.
func (m Mask) Visible(elevationRad float64) bool {
	return applyVis(elevationRad >= Radians(m.Deg))
}

// FilteredSight is a satellite that survived the elevation mask, carrying
// its direction together with the local elevation and azimuth for reports.
type FilteredSight struct {
	SatID  string
	Unit   wgs84.ECEF
	RangeM float64
	ElRad  float64
	AzRad  float64
}

// ElevationDeg returns the elevation of the sight in degrees.
func (f FilteredSight) ElevationDeg() float64 {
	return Degrees(f.ElRad)
}

// AzimuthDeg returns the azimuth of the sight in degrees.
func (f FilteredSight) AzimuthDeg() float64 {
	return Degrees(f.AzRad)
}

// Filter applies the mask to a list of sights and returns the visible
// subset together with the rejected identifiers. The receiver basis is
// evaluated once and reused for every satellite, keeping the ENU rotation
// identical across the whole constellation.
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
