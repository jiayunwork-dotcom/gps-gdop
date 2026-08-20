package dop

import (
	"errors"
	"fmt"
	"math"
)

var ErrSingularNormal = errors.New("normal matrix is singular: constellation cannot resolve a position and clock")

// Invert4 computes the inverse of the 4x4 matrix with Gauss-Jordan
// elimination and partial pivoting. The pivot threshold is relative to
// the largest entry, so a coplanar constellation (whose normal matrix is
// rank deficient because every direction cosine lies in one plane) is
// reported as singular instead of producing a garbage inverse.
func Invert4(m Mat4) (Mat4, error) {
	var work Mat4
	var inv Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			work.Data[i][j] = m.Data[i][j]
			if i == j {
				inv.Data[i][j] = 1.0
			}
		}
	}
	scale := m.MaxAbs()
	if scale == 0 || !isFinite(scale) {
		return Mat4{}, fmt.Errorf("%w: zero or non-finite entries", ErrSingularNormal)
	}
	threshold := 1e-12 * scale
	for col := 0; col < 4; col++ {
		pivotRow := col
		pivot := math.Abs(work.Data[col][col])
		for row := col + 1; row < 4; row++ {
			if v := math.Abs(work.Data[row][col]); v > pivot {
				pivot = v
				pivotRow = row
			}
		}
		if pivot < threshold {
			return Mat4{}, fmt.Errorf("%w: pivot %.3e below threshold %.3e", ErrSingularNormal, pivot, threshold)
		}
		if pivotRow != col {
			for j := 0; j < 4; j++ {
				work.Data[col][j], work.Data[pivotRow][j] = work.Data[pivotRow][j], work.Data[col][j]
				inv.Data[col][j], inv.Data[pivotRow][j] = inv.Data[pivotRow][j], inv.Data[col][j]
			}
		}
		divisor := work.Data[col][col]
		for j := 0; j < 4; j++ {
			work.Data[col][j] /= divisor
			inv.Data[col][j] /= divisor
		}
		for row := 0; row < 4; row++ {
			if row == col {
				continue
			}
			factor := work.Data[row][col]
			for j := 0; j < 4; j++ {
				work.Data[row][j] -= factor * work.Data[col][j]
				inv.Data[row][j] -= factor * inv.Data[col][j]
			}
		}
	}
	return inv, nil
}

func isFinite(v float64) bool {
	return v == v && !(v > 1.7976931348623157e308 || v < -1.7976931348623157e308)
}

// Multiply returns the product of two 4x4 matrices, used by tests to
// verify that the inverse really inverts.
func Multiply(a, b Mat4) Mat4 {
	var out Mat4
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := 0.0
			for k := 0; k < 4; k++ {
				sum += a.Data[i][k] * b.Data[k][j]
			}
			out.Data[i][j] = sum
		}
	}
	return out
}

// Identity returns the 4x4 identity matrix.
func Identity() Mat4 {
	var m Mat4
	for i := 0; i < 4; i++ {
		m.Data[i][i] = 1.0
	}
	return m
}
