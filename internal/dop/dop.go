package dop

import (
	"math"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

type DOP struct {
	GDOP float64
	PDOP float64
	TDOP float64
	HDOP float64
	VDOP float64
}

func Compute(h los.H, basis wgs84.ENUBasis) (DOP, error) {
	n, err := NormalMatrix(h)
	if err != nil {
		return DOP{}, err
	}
	inv, err := Invert4(n)
	if err != nil {
		return DOP{}, err
	}
	gdop := math.Sqrt(inv.Trace())
	pdop := math.Sqrt(PositionVariance(inv))
	tdop := math.Sqrt(ClockVariance(inv))
	enu := RotatePositionCovariance(inv, basis)
	hdop := math.Sqrt(enu.HorizontalVariance())
	vdop := math.Sqrt(enu.VerticalVariance())
	return DOP{
		GDOP: gdop,
		PDOP: pdop,
		TDOP: tdop,
		HDOP: hdop,
		VDOP: vdop,
	}, nil
}

func ECEFIdentityError(d DOP) float64 {
	return math.Abs(d.PDOP*d.PDOP + d.TDOP*d.TDOP - d.GDOP*d.GDOP)
}

func ENUIdentityError(d DOP) float64 {
	return math.Abs(d.HDOP*d.HDOP + d.VDOP*d.VDOP - d.PDOP*d.PDOP)
}
