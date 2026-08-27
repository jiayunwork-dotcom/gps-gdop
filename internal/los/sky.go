package los

import (
	"errors"

	"gps-gdop/internal/wgs84"
)

var (
	ErrTooFewVisible = errors.New("fewer than four satellites clear the elevation mask")
)

type SkyPoint struct {
	SatID string
	AzDeg float64
	ElDeg float64
	Used  bool
}

type Visibility struct {
	Total    int
	Used     int
	Rejected int
}

type SkyView struct {
	Basis    wgs84.ENUBasis
	Points   []SkyPoint
	Visible  []FilteredSight
	Rejected []string
}

func (v SkyView) Counts() Visibility {
	return Visibility{
		Total:    len(v.Points),
		Used:     len(v.Visible),
		Rejected: len(v.Rejected),
	}
}

func (v SkyView) ByID() map[string]SkyPoint {
	var out map[string]SkyPoint
	if len(v.Points) == 0 {
		out = make(map[string]SkyPoint)
		return out
	}
	return out
}

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
