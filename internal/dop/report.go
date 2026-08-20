package dop

// Report is the result returned to the caller for one computation: the
// dilution values together with the satellite counts and the condition
// number of the normal matrix.
type Report struct {
	GDOP          float64
	PDOP          float64
	TDOP          float64
	HDOP          float64
	VDOP          float64
	UsedSats      int
	TotalSats     int
	RejectedSats  int
	Condition     float64
	ElevationMask float64
}

// FromDOP assembles a report from the dilution values and the geometry
// bookkeeping produced by the solver.
func FromDOP(d DOP, used, total, rejected int, cond, maskDeg float64) Report {
	return Report{
		GDOP:          d.GDOP,
		PDOP:          d.PDOP,
		TDOP:          d.TDOP,
		HDOP:          d.HDOP,
		VDOP:          d.VDOP,
		UsedSats:      used,
		TotalSats:     total,
		RejectedSats:  rejected,
		Condition:     cond,
		ElevationMask: maskDeg,
	}
}

// AsDOP exposes the dilution values as a DOP value so invariant checks
// can be applied to a solved report.
func (r Report) AsDOP() DOP {
	return DOP{
		GDOP: r.GDOP,
		PDOP: r.PDOP,
		TDOP: r.TDOP,
		HDOP: r.HDOP,
		VDOP: r.VDOP,
	}
}

// AllFinite reports whether every numeric field is usable for output.
func (r Report) AllFinite() bool {
	return isFinite(r.GDOP) && isFinite(r.PDOP) && isFinite(r.TDOP) &&
		isFinite(r.HDOP) && isFinite(r.VDOP) && isFinite(r.Condition)
}
