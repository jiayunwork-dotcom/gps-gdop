package web

import (
	"net/http"

	"gps-gdop/internal/dop"
)

func (s *Server) handleDop(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	var req DopRequest
	if err := decodeBody(r, &req); err != nil {
		badRequest(w, err)
		return
	}
	solveIn := req.ToSolveInput()
	result, err := dop.Solve(solveIn)
	if err != nil {
		badRequest(w, err)
		return
	}
	rep := result.Report
	writeJSON(w, http.StatusOK, DopResponse{
		GDOP:               rep.GDOP,
		PDOP:               rep.PDOP,
		TDOP:               rep.TDOP,
		HDOP:               rep.HDOP,
		VDOP:               rep.VDOP,
		SatellitesUsed:     rep.UsedSats,
		SatellitesTotal:    rep.TotalSats,
		SatellitesRejected: rep.RejectedSats,
		ConditionNumber:    rep.Condition,
		ElevationMaskDeg:   rep.ElevationMask,
	})
}
