package los

import (
	"errors"

	"gps-gdop/internal/wgs84"
)

var (
	ErrTooFewVisible = errors.New("fewer than four satellites clear the elevation mask")
)

// SkyPoint is one satellite's position on the sky dome, ready to be
// plotted on a zenith chart.
type SkyPoint struct {
	SatID string
	AzDeg float64
	ElDeg float64
	Used  bool
}

// Visibility tallies how many satellites entered the computation and how
// many were rejected by the elevation mask.
type Visibility struct {
	Total    int
	Used     int
	Rejected int
}

// SkyView is the complete result for the /api/sky endpoint: every
// satellite with its azimuth and elevation plus the visibility summary.
type SkyView struct {
	Basis    wgs84.ENUBasis
	Points   []SkyPoint
	Visible  []FilteredSight
	Rejected []string
}

// Counts returns the visibility summary derived from the view.
func (v SkyView) Counts() Visibility {
	return Visibility{
		Total:    len(v.Points),
		Used:     len(v.Visible),
		Rejected: len(v.Rejected),
	}
}

// BuildSkyView computes the topocentric direction of every satellite.
// It needs the list of raw sights and the receiver basis; the caller is
// responsible for building both from the shared WGS84 conversions.
func BuildSkyView(sights []LineOfSight, basis wgs84.ENUBasis, mask Mask) (SkyView, error) {
	visible, rejected := Filter(sights, basis, mask)
	if len(visible) < 4 {
		return SkyView{}, ErrTooFewVisible
	}
	points := make([]SkyPoint, 0, len(sights))
	usedIDs := make(map[string]bool, len(visible))
	for _, f := range visible {
		usedIDs[f.SatID] = true
	}
	for _, s := range sights {
		el := ElevationAngle(s.Unit, basis)
		az := AzimuthAngle(s.Unit, basis)
		points = append(points, SkyPoint{
			SatID: s.SatID,
			AzDeg: Degrees(az),
			ElDeg: Degrees(el),
			Used:  usedIDs[s.SatID],
		})
	}
	return SkyView{
		Basis:    basis,
		Points:   points,
		Visible:  visible,
		Rejected: rejected,
	}, nil
}
