package dop

import (
	"errors"
	"math"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

var ErrBadWeight = errors.New("dop: satellite weight must be positive and finite")

func ElevationWeight(elRad float64) (float64, error) {
	s := math.Sin(elRad)
	if s <= 0 {
		return 0, ErrBadWeight
	}
	w := s * s
	if w != w || w <= 0 {
		return 0, ErrBadWeight
	}
	return w, nil
}

func WeightsFromVisible(visible []los.FilteredSight) ([]float64, error) {
	out := make([]float64, len(visible))
	for i, f := range visible {
		w, err := ElevationWeight(f.ElRad)
		if err != nil {
			return nil, err
		}
		out[i] = w
	}
	return out, nil
}

func WeightedNormal(h los.H, weights []float64) (Mat4, error) {
	if err := h.Validate(); err != nil {
		return Mat4{}, err
	}
	if len(weights) != h.Rows {
		return Mat4{}, errors.New("dop: weight count must match H rows")
	}
	var n Mat4
	for i := 0; i < h.Rows; i++ {
		w := weights[i]
		if w <= 0 || w != w {
			return Mat4{}, ErrBadWeight
		}
		row := h.Data[i]
		for c := 0; c < 4; c++ {
			for d := 0; d < 4; d++ {
				n.Data[c][d] += w * row[c] * row[d]
			}
		}
	}
	return n, nil
}

func ComputeWeighted(h los.H, basis wgs84.ENUBasis, weights []float64) (DOP, error) {
	n, err := WeightedNormal(h, weights)
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

func UniformWeights(n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = 1
	}
	return out
}

func UnitWeightsMatchUnweighted(h los.H, basis wgs84.ENUBasis, tol float64) (bool, error) {
	plain, err := Compute(h, basis)
	if err != nil {
		return false, err
	}
	w, err := ComputeWeighted(h, basis, UniformWeights(h.Rows))
	if err != nil {
		return false, err
	}
	return absd(plain.GDOP-w.GDOP) <= tol && absd(plain.HDOP-w.HDOP) <= tol, nil
}

func absd(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
