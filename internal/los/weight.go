package los

import (
	"errors"
	"math"
)

func ScaleH(h H, weights []float64) (H, error) {
	if err := h.Validate(); err != nil {
		return H{}, err
	}
	if len(weights) != h.Rows {
		return H{}, errors.New("los: weight count must match H rows")
	}
	data := make([][]float64, h.Rows)
	for i := 0; i < h.Rows; i++ {
		if weights[i] <= 0 || weights[i] != weights[i] {
			return H{}, errors.New("los: weight must be positive")
		}
		s := math.Sqrt(weights[i])
		row := make([]float64, 4)
		for j := 0; j < 4; j++ {
			row[j] = s * h.Data[i][j]
		}
		data[i] = row
	}
	return H{Rows: h.Rows, Cols: 4, Data: data}, nil
}

func ClockColumnStillOnes(h H, tol float64) bool {
	for i := 0; i < h.Rows; i++ {
		d := h.Data[i][3] - 1
		if d < 0 {
			d = -d
		}
		if d > tol {
			return false
		}
	}
	return true
}

func RowUnitPrefix(h H, i int) float64 {
	if i < 0 || i >= h.Rows || len(h.Data[i]) < 3 {
		return 0
	}
	x, y, z := h.Data[i][0], h.Data[i][1], h.Data[i][2]
	return math.Sqrt(x*x + y*y + z*z)
}

func AllRowsUnitPrefix(h H, tol float64) bool {
	for i := 0; i < h.Rows; i++ {
		n := RowUnitPrefix(h, i)
		d := n - 1
		if d < 0 {
			d = -d
		}
		if d > tol {
			return false
		}
	}
	return true
}
