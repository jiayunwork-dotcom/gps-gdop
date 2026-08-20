package dop

import (
	"math"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

// DOP carries the five dilution-of-precision values derived from the same
// normal matrix. GDOP, PDOP and TDOP come straight from the ECEF
// covariance diagonal; HDOP and VDOP come from that covariance rotated
// into the receiver's local ENU frame.
type DOP struct {
	GDOP float64
	PDOP float64
	TDOP float64
	HDOP float64
	VDOP float64
}

// Compute builds the normal matrix from the geometry rows, inverts it and
// extracts every dilution value. A singular normal matrix (for example a
// coplanar constellation that leaves the clock unobservable) returns an
// error instead of a meaningless number.
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
	return fillDOP(DOP{
		GDOP: gdop,
		PDOP: pdop,
		TDOP: tdop,
		HDOP: hdop,
		VDOP: vdop,
	}), nil
}

// ECEFIdentityError reports how far PDOP^2 + TDOP^2 deviates from GDOP^2.
// Both sides derive from the trace of the same ECEF covariance, so the
// deviation must sit at machine precision.
func ECEFIdentityError(d DOP) float64 {
	return math.Abs(d.PDOP*d.PDOP + d.TDOP*d.TDOP - d.GDOP*d.GDOP)
}

// ENUIdentityError reports how far HDOP^2 + VDOP^2 deviates from PDOP^2.
// The ENU rotation preserves the total position variance.
func ENUIdentityError(d DOP) float64 {
	return math.Abs(d.HDOP*d.HDOP + d.VDOP*d.VDOP - d.PDOP*d.PDOP)
}
