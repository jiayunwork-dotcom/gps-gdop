package los

import (
	"math"

	"gps-gdop/internal/wgs84"
)

// ElevationAngle converts a line-of-sight direction into the elevation
// angle above the local horizon. The ECEF unit vector is first rotated
// into the ENU frame of the receiver basis, then
// el = atan2(u, sqrt(e^2 + n^2)).
func ElevationAngle(unit wgs84.ECEF, basis wgs84.ENUBasis) float64 {
	local := basis.ToENU(unit)
	return math.Atan2(local.Up, math.Hypot(local.East, local.North))
}

// AzimuthAngle converts a line-of-sight direction into the azimuth angle
// measured clockwise from north: az = atan2(e, n), folded into [0, 360).
func AzimuthAngle(unit wgs84.ECEF, basis wgs84.ENUBasis) float64 {
	local := basis.ToENU(unit)
	az := math.Atan2(local.East, local.North)
	if az < 0 {
		az += 2 * math.Pi
	}
	return az
}

// HorizontalRange extracts the topocentric horizontal projection length.
func HorizontalRange(unit wgs84.ECEF, basis wgs84.ENUBasis) float64 {
	local := basis.ToENU(unit)
	return math.Hypot(local.East, local.North)
}

// Degrees converts radians to degrees.
func Degrees(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// Radians converts degrees to radians.
func Radians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// NormalizeAzimuth folds an angle in radians into [0, 2*pi).
func NormalizeAzimuth(rad float64) float64 {
	twoPi := 2 * math.Pi
	rad = math.Mod(rad, twoPi)
	if rad < 0 {
		rad += twoPi
	}
	return rad
}
