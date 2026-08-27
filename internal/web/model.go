package web

import (
	"gps-gdop/internal/dop"
	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

type DopRequest struct {
	ReceiverECEF     ECEFPoint      `json:"receiver_ecef"`
	Satellites       []SatelliteDTO `json:"satellites"`
	ElevationMaskDeg *float64       `json:"elevation_mask_deg,omitempty"`
}

type ECEFPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type SatelliteDTO struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Z  float64 `json:"z"`
}

func (r DopRequest) ToSolveInput() dop.SolveInput {
	return dop.SolveInput{
		Receiver:   wgs84.ECEF{X: r.ReceiverECEF.X, Y: r.ReceiverECEF.Y, Z: r.ReceiverECEF.Z},
		Satellites: r.satellites(),
		MaskDeg:    r.maskDeg(),
	}
}

func (r DopRequest) satellites() []los.Satellite {
	sats := make([]los.Satellite, 0, len(r.Satellites))
	for _, s := range r.Satellites {
		sats = append(sats, los.Satellite{
			ID:   s.ID,
			ECEF: wgs84.ECEF{X: s.X, Y: s.Y, Z: s.Z},
		})
	}
	return sats
}

func (r DopRequest) maskDeg() float64 {
	if r.ElevationMaskDeg != nil {
		return *r.ElevationMaskDeg
	}
	return 5.0
}

type DopResponse struct {
	GDOP               float64 `json:"gdop"`
	PDOP               float64 `json:"pdop"`
	TDOP               float64 `json:"tdop"`
	HDOP               float64 `json:"hdop"`
	VDOP               float64 `json:"vdop"`
	SatellitesUsed     int     `json:"satellites_used"`
	SatellitesTotal    int     `json:"satellites_total"`
	SatellitesRejected int     `json:"satellites_rejected"`
	ConditionNumber    float64 `json:"condition_number"`
	ElevationMaskDeg   float64 `json:"elevation_mask_deg"`
}

type SkyPointDTO struct {
	ID           string  `json:"id"`
	AzimuthDeg   float64 `json:"azimuth_deg"`
	ElevationDeg float64 `json:"elevation_deg"`
	Used         bool    `json:"used"`
}

type SkyResponse struct {
	ElevationMaskDeg float64       `json:"elevation_mask_deg"`
	Satellites       []SkyPointDTO `json:"satellites"`
	UsedCount        int           `json:"used_count"`
	RejectedCount    int           `json:"rejected_count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
