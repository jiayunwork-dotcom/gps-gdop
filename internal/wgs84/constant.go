package wgs84

import "math"

type Ellipsoid struct {
	SemiMajorAxis        float64
	Flattening           float64
	SemiMinorAxis        float64
	EccentricitySq       float64
	SecondEccentricitySq float64
	GM                   float64
	RotationRate         float64
}

func NewWGS84() Ellipsoid {
	a := 6378137.0
	f := 1.0 / 298.257223563
	b := a * (1.0 - f)
	e2 := f * (2.0 - f)
	ep2 := e2 / (1.0 - e2)
	return Ellipsoid{
		SemiMajorAxis:        a,
		Flattening:           f,
		SemiMinorAxis:        b,
		EccentricitySq:       e2,
		SecondEccentricitySq: ep2,
		GM:                   3.986004418e14,
		RotationRate:         7.2921150e-5,
	}
}

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

func (e Ellipsoid) Eccentricity() float64 {
	return math.Sqrt(e.EccentricitySq)
}

func (e Ellipsoid) PrimeVerticalRadius(latRad float64) float64 {
	sinLat := math.Sin(latRad)
	return e.SemiMajorAxis / math.Sqrt(1.0-e.EccentricitySq*sinLat*sinLat)
}

func (e Ellipsoid) MeanEarthRadius() float64 {
	return (2.0*e.SemiMajorAxis + e.SemiMinorAxis) / 3.0
}
