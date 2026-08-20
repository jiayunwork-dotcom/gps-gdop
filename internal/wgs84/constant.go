package wgs84

import "math"

// Ellipsoid carries the geodetic reference figure shared by every
// coordinate transform in the project: geodetic<->ECEF conversion,
// the ENU topocentric basis and the covariance rotations that turn
// ECEF dilution values into local East/North/Up components.
type Ellipsoid struct {
	SemiMajorAxis float64
	Flattening    float64
	SemiMinorAxis float64
	EccentricitySq float64
	SecondEccentricitySq float64
	GM            float64
	RotationRate  float64
}

// NewWGS84 returns the WGS84 reference ellipsoid. The constants are
// fixed by convention and are the single source of truth for the whole
// repository; the README pins them next to the formulas that consume them.
func NewWGS84() Ellipsoid {
	a := 6378137.0
	f := 1.0 / 298.257223563
	b := a * (1.0 - f)
	e2 := f * (2.0 - f)
	ep2 := e2 / (1.0 - e2)
	return Ellipsoid{
		SemiMajorAxis:       a,
		Flattening:          f,
		SemiMinorAxis:       b,
		EccentricitySq:      e2,
		SecondEccentricitySq: ep2,
		GM:                   3.986004418e14,
		RotationRate:         7.2921150e-5,
	}
}

// Validate rejects ellipsoids whose defining parameters are not finite
// or do not satisfy the geometric identities every formula relies on.
func (e Ellipsoid) Validate() error {
	if !IsFinite(e.SemiMajorAxis) || !IsFinite(e.SemiMinorAxis) {
		return ErrEllipsoidNotFinite
	}
	if e.SemiMajorAxis <= 0 || e.SemiMinorAxis <= 0 {
		return ErrEllipsoidDegenerate
	}
	if e.SemiMinorAxis > e.SemiMajorAxis {
		return ErrEllipsoidMinorAxisLonger
	}
	if e.Flattening <= 0 || e.Flattening >= 1 {
		return ErrEllipsoidFlattening
	}
	return nil
}

// Eccentricity derives the first eccentricity from the squared value.
func (e Ellipsoid) Eccentricity() float64 {
	return math.Sqrt(e.EccentricitySq)
}

// PrimeVerticalRadius is the radius of curvature of the prime vertical
// at the given geodetic latitude, used by both conversion directions.
func (e Ellipsoid) PrimeVerticalRadius(latRad float64) float64 {
	sinLat := math.Sin(latRad)
	return e.SemiMajorAxis / math.Sqrt(1.0-e.EccentricitySq*sinLat*sinLat)
}

// MeanEarthRadius is a convenience figure for reports that want a single
// scalar describing the size of the reference ellipsoid.
func (e Ellipsoid) MeanEarthRadius() float64 {
	return (2.0*e.SemiMajorAxis + e.SemiMinorAxis) / 3.0
}
