package los

import (
	"math"

	"gps-gdop/internal/wgs84"
)

func ElevationAngle(unit wgs84.ECEF, basis wgs84.ENUBasis) float64 {
	local := basis.ToENU(unit)
	return math.Atan2(local.Up, math.Hypot(local.East, local.North))
}

func AzimuthAngle(unit wgs84.ECEF, basis wgs84.ENUBasis) float64 {
	local := basis.ToENU(unit)
	az := math.Atan2(local.East, local.North)
	if az < 0 {
		az += 2 * math.Pi
	}
	return az
}

func HorizontalRange(unit wgs84.ECEF, basis wgs84.ENUBasis) float64 {
	local := basis.ToENU(unit)
	return math.Hypot(local.East, local.North)
}

func Degrees(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func Radians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func NormalizeAzimuth(rad float64) float64 {
	twoPi := 2 * math.Pi
	rad = math.Mod(rad, twoPi)
	if rad < 0 {
		rad += twoPi
	}
	return rad
}
