package wgs84

import (
	"fmt"
	"math"
)

// LatLonHeight is a geodetic position: latitude and longitude in radians,
// ellipsoidal height in metres. Latitude is signed north positive and
// longitude is signed east of Greenwich.
type LatLonHeight struct {
	LatRad float64
	LonRad float64
	HeightM float64
}

// NewLatLonHeightDeg validates a position given in degrees.
func NewLatLonHeightDeg(latDeg, lonDeg, heightM float64) (LatLonHeight, error) {
	if err := RequireFinite(latDeg, lonDeg, heightM); err != nil {
		return LatLonHeight{}, err
	}
	return LatLonHeight{
		LatRad:  latDeg * math.Pi / 180.0,
		LonRad:  lonDeg * math.Pi / 180.0,
		HeightM: heightM,
	}, nil
}

// ToECEF converts the geodetic position to ECEF using the shared WGS84
// ellipsoid. The formula is the standard closed form with the prime
// vertical radius of curvature.
func (g LatLonHeight) ToECEF(e Ellipsoid) ECEF {
	sinLat := math.Sin(g.LatRad)
	cosLat := math.Cos(g.LatRad)
	sinLon := math.Sin(g.LonRad)
	cosLon := math.Cos(g.LonRad)
	n := e.PrimeVerticalRadius(g.LatRad)
	zPart := (n*(1.0-e.EccentricitySq) + g.HeightM) * sinLat
	return ECEF{
		X: (n + g.HeightM) * cosLat * cosLon,
		Y: (n + g.HeightM) * cosLat * sinLon,
		Z: zPart,
	}
}

// FromECEF converts an ECEF point to geodetic coordinates with Bowring's
// closed-form iteration, which converges to sub-millimetre accuracy on a
// single pass for all terrestrial positions.
func FromECEF(p ECEF, e Ellipsoid) (LatLonHeight, error) {
	if err := p.ValidateFinite(); err != nil {
		return LatLonHeight{}, err
	}
	lon := math.Atan2(p.Y, p.X)
	radius := math.Hypot(p.X, p.Y)
	lat, height := solveGeodetic(p, e, radius)
	return LatLonHeight{LatRad: lat, LonRad: lon, HeightM: height}, nil
}

// solveGeodetic computes latitude and height given the equatorial radius.
// It handles the polar degeneracy where the equatorial radius is zero.
func solveGeodetic(p ECEF, e Ellipsoid, radius float64) (float64, float64) {
	if radius < 1e-9 {
		lat := math.Pi / 2.0
		if p.Z < 0 {
			lat = -lat
		}
		height := math.Abs(p.Z) - e.SemiMinorAxis
		return lat, height
	}
	a := e.SemiMajorAxis
	b := e.SemiMinorAxis
	theta := math.Atan2(p.Z*a, radius*b)
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)
	ep2 := e.SecondEccentricitySq
	lat := math.Atan2(
		p.Z+ep2*b*sinTheta*sinTheta*sinTheta,
		radius-e.EccentricitySq*a*cosTheta*cosTheta*cosTheta,
	)
	n := e.PrimeVerticalRadius(lat)
	cosLat := math.Cos(lat)
	var height float64
	if math.Abs(cosLat) < 1e-12 {
		height = math.Abs(p.Z) - b
	} else {
		height = radius/cosLat - n
	}
	return lat, height
}

// String renders the position in degrees for readability.
func (g LatLonHeight) String() string {
	return fmt.Sprintf("lat %.8f deg, lon %.8f deg, h %.3f m",
		g.LatRad*180.0/math.Pi, g.LonRad*180.0/math.Pi, g.HeightM)
}

// ValidateFinite checks every component of the ECEF point.
func (p ECEF) ValidateFinite() error {
	return RequireFinite(p.X, p.Y, p.Z)
}
