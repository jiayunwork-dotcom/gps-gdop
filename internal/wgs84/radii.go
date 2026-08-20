package wgs84

import "math"

// MeridianRadius is the radius of curvature of the meridian ellipse at a
// given geodetic latitude, used when geodetic quantities must be scaled
// into metres along the north direction.
func (e Ellipsoid) MeridianRadius(latRad float64) float64 {
	sinLat := math.Sin(latRad)
	denom := 1.0 - e.EccentricitySq*sinLat*sinLat
	return e.SemiMajorAxis * (1.0 - e.EccentricitySq) / math.Pow(denom, 1.5)
}

// GeocentricRadius is the distance from the geocentre to a point on the
// ellipsoid surface at the given geodetic latitude.
func (e Ellipsoid) GeocentricRadius(latRad float64) float64 {
	cosLat := math.Cos(latRad)
	sinLat := math.Sin(latRad)
	eq := e.SemiMajorAxis * cosLat
	pol := e.SemiMinorAxis * sinLat
	return math.Hypot(eq, pol)
}

// RectifyingRadius is the mean radius of a circle whose circumference
// equals the ellipsoid meridian arc, a scalar used to express latitudinal
// spans in metres without integrating.
func (e Ellipsoid) RectifyingRadius() float64 {
	a := e.SemiMajorAxis
	b := e.SemiMinorAxis
	return (2.0*a*a + b*b) / (3.0 * a)
}

// DegreesToMetres scales a small latitude span on the meridian to metres.
func (e Ellipsoid) DegreesToMetres(latDeg float64, latRad float64) float64 {
	return latDeg * math.Pi / 180.0 * e.MeridianRadius(latRad)
}
