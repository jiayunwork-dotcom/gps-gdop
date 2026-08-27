package web

import (
	"net/http"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

func (s *Server) handleSky(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var req DopRequest
	if err := decodeBody(r, &req); err != nil {
		badRequest(w, err)
		return
	}
	solveIn := req.ToSolveInput()
	mask, err := los.NewMask(solveIn.MaskDeg)
	if err != nil {
		badRequest(w, err)
		return
	}
	basis, err := wgs84.BasisAtECEF(solveIn.Receiver, wgs84.NewWGS84())
	if err != nil {
		badRequest(w, err)
		return
	}
	sights, err := los.ComputeAllSights(los.Receiver{ECEF: solveIn.Receiver}, solveIn.Satellites)
	if err != nil {
		badRequest(w, err)
		return
	}
	if sights.Len() < 4 {
		badRequest(w, los.ErrTooFewSights)
		return
	}
	view, err := los.BuildSkyView(sights.Items, basis, mask)
	if err != nil {
		badRequest(w, err)
		return
	}
	points := make([]SkyPointDTO, 0, len(view.Points))
	idx := view.ByID()
	for _, p := range view.Points {
		idx[p.SatID] = p
		points = append(points, SkyPointDTO{
			ID:           p.SatID,
			AzimuthDeg:   p.AzDeg,
			ElevationDeg: p.ElDeg,
			Used:         p.Used,
		})
	}
	counts := view.Counts()
	writeJSON(w, http.StatusOK, SkyResponse{
		ElevationMaskDeg: mask.Deg,
		Satellites:       points,
		UsedCount:        counts.Used,
		RejectedCount:    counts.Rejected,
	})
}
