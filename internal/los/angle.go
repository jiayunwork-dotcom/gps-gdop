package los

import (
	"math"

	"gps-gdop/internal/wgs84"
)

func AngleBetween(a, b wgs84.ECEF) float64 {
	na := a.Norm()
	nb := b.Norm()
	if na <= 0 || nb <= 0 {
		return 0
	}
	c := a.Dot(b) / (na * nb)
	if c > 1 {
		c = 1
	}
	if c < -1 {
		c = -1
	}
	return math.Acos(c)
}

func MinPairAngle(sights []LineOfSight) (float64, bool) {
	if len(sights) < 2 {
		return 0, false
	}
	best := math.Pi
	ok := false
	for i := 0; i < len(sights); i++ {
		for j := i + 1; j < len(sights); j++ {
			ang := AngleBetween(sights[i].Unit, sights[j].Unit)
			if !ok || ang < best {
				best = ang
				ok = true
			}
		}
	}
	return best, ok
}

func MaxPairAngle(sights []LineOfSight) (float64, bool) {
	if len(sights) < 2 {
		return 0, false
	}
	best := 0.0
	ok := false
	for i := 0; i < len(sights); i++ {
		for j := i + 1; j < len(sights); j++ {
			ang := AngleBetween(sights[i].Unit, sights[j].Unit)
			if !ok || ang > best {
				best = ang
				ok = true
			}
		}
	}
	return best, ok
}

func SpreadScore(sights []LineOfSight) float64 {
	minA, okMin := MinPairAngle(sights)
	maxA, okMax := MaxPairAngle(sights)
	if !okMin || !okMax || maxA <= 0 {
		return 0
	}
	return minA / maxA
}

func MinPairAngleFiltered(visible []FilteredSight) (float64, bool) {
	sights := make([]LineOfSight, len(visible))
	for i, f := range visible {
		sights[i] = LineOfSight{SatID: f.SatID, Unit: f.Unit, RangeM: f.RangeM}
	}
	return MinPairAngle(sights)
}

func SpreadOfVisible(visible []FilteredSight) float64 {
	sights := make([]LineOfSight, len(visible))
	for i, f := range visible {
		sights[i] = LineOfSight{SatID: f.SatID, Unit: f.Unit, RangeM: f.RangeM}
	}
	return SpreadScore(sights)
}

func ElevationSpan(visible []FilteredSight) (lo, hi float64, ok bool) {
	if len(visible) == 0 {
		return 0, 0, false
	}
	lo = visible[0].ElRad
	hi = visible[0].ElRad
	for _, f := range visible[1:] {
		if f.ElRad < lo {
			lo = f.ElRad
		}
		if f.ElRad > hi {
			hi = f.ElRad
		}
	}
	return lo, hi, true
}

func AzimuthWrapDelta(a, b float64) float64 {
	d := NormalizeAzimuth(a) - NormalizeAzimuth(b)
	if d > math.Pi {
		d -= 2 * math.Pi
	}
	if d < -math.Pi {
		d += 2 * math.Pi
	}
	if d < 0 {
		d = -d
	}
	return d
}
