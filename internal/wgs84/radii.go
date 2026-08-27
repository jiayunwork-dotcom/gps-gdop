package wgs84

import "math"

func (e Ellipsoid) MeridianRadius(latRad float64) float64 {
	sinLat := math.Sin(latRad)
	denom := 1.0 - e.EccentricitySq*sinLat*sinLat
	return e.SemiMajorAxis * (1.0 - e.EccentricitySq) / math.Pow(denom, 1.5)
}

func (e Ellipsoid) GeocentricRadius(latRad float64) float64 {
	cosLat := math.Cos(latRad)
	sinLat := math.Sin(latRad)
	eq := e.SemiMajorAxis * cosLat
	pol := e.SemiMinorAxis * sinLat
	return math.Hypot(eq, pol)
}

func (e Ellipsoid) RectifyingRadius() float64 {
	a := e.SemiMajorAxis
	b := e.SemiMinorAxis
	return (2.0*a*a + b*b) / (3.0 * a)
}

func (e Ellipsoid) DegreesToMetres(latDeg float64, latRad float64) float64 {
	return latDeg * math.Pi / 180.0 * e.MeridianRadius(latRad)
}
