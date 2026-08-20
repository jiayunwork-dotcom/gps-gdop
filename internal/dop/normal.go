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

// Mat4 is a dense 4x4 matrix stored row-major, used for the normal matrix
// N = H^T H and its inverse.
type Mat4 struct {
	Data [4][4]float64
}

// NormalMatrix computes N = H^T H by accumulating the outer product of
// every geometry row. The 4th diagonal entry is exactly the number of
// participating satellites because the clock column is all ones.
func NormalMatrix(h los.H) (Mat4, error) {
	if err := h.Validate(); err != nil {
		return Mat4{}, err
	}
	var n Mat4
	for i := 0; i < h.Rows; i++ {
		row := h.Data[i]
		for c := 0; c < 4; c++ {
			for d := 0; d < 4; d++ {
				n.Data[c][d] += row[c] * row[d]
			}
		}
	}
	return n, nil
}

// Trace returns the sum of the diagonal entries.
func (m Mat4) Trace() float64 {
	return m.Data[0][0] + m.Data[1][1] + m.Data[2][2] + m.Data[3][3]
}

// MaxAbs returns the largest absolute entry, used as the scale for the
// singularity threshold during inversion.
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
