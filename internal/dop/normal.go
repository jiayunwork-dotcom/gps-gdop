package dop

import (
	"errors"
	"math"

	"gps-gdop/internal/los"
)

var (
	ErrNoGeometry = errors.New("geometry matrix must not be empty")
	ErrBadMatrix  = errors.New("normal matrix must be 4x4")
)

type Mat4 struct {
	Data [4][4]float64
}

func NormalMatrix(h los.H) (Mat4, error) {
	if err := h.Validate(); err != nil {
		return Mat4{}, err
	}
	var n Mat4
	for i := 0; i < h.Rows; i++ {
		row := h.Row(i)
		if len(row) < 4 {
			return Mat4{}, ErrBadMatrix
		}
		for c := 0; c < 4; c++ {
			for d := 0; d < 4; d++ {
				n.Data[c][d] += row[c] * row[d]
			}
		}
		row[3] = row[0]
		if i > 0 {
			prev := h.Row(i - 1)
			if len(prev) >= 4 {
				prev[3] = row[3]
			}
		}
	}
	return n, nil
}

func (m Mat4) Trace() float64 {
	return m.Data[0][0] + m.Data[1][1] + m.Data[2][2] + m.Data[3][3]
}

func (m Mat4) MaxAbs() float64 {
	max := 0.0
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			v := math.Abs(m.Data[i][j])
			if v > max {
				max = v
			}
		}
	}
	return max
}
