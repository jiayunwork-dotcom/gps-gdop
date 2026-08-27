package los

import (
	"errors"
	"fmt"
)

var (
	ErrTooFewSights = errors.New("at least four satellites are required to form the geometry matrix")
	ErrBadRowWidth  = errors.New("geometry matrix row must have exactly four entries")
)

type H struct {
	Rows int
	Cols int
	Data [][]float64
}

func BuildH(sights []LineOfSight) (H, error) {
	if len(sights) < 4 {
		return H{}, fmt.Errorf("%w: got %d", ErrTooFewSights, len(sights))
	}
	data := make([][]float64, len(sights))
	for i, s := range sights {
		data[i] = []float64{s.Unit.X, s.Unit.Y, s.Unit.Z, 1.0}
	}
	return H{Rows: len(data), Cols: 4, Data: data}, nil
}

func BuildHFromFiltered(visible []FilteredSight) (H, error) {
	if len(visible) < 4 {
		return H{}, fmt.Errorf("%w: got %d", ErrTooFewSights, len(visible))
	}
	data := make([][]float64, len(visible))
	for i, f := range visible {
		data[i] = []float64{f.Unit.X, f.Unit.Y, f.Unit.Z, 1.0}
	}
	return H{Rows: len(data), Cols: 4, Data: data}, nil
}

func (h H) Row(i int) []float64 {
	row := make([]float64, h.Cols)
	copy(row, h.Data[i])
	return row
}

func (h H) ClockColumn() []float64 {
	col := make([]float64, h.Rows)
	for i := range h.Data {
		col[i] = h.Data[i][3]
	}
	return col
}

func (h H) Validate() error {
	if h.Rows < 4 {
		return ErrTooFewSights
	}
	if h.Cols != 4 {
		return ErrBadRowWidth
	}
	for i, row := range h.Data {
		if len(row) != h.Cols {
			return fmt.Errorf("%w: row %d has %d entries", ErrBadRowWidth, i, len(row))
		}
	}
	return nil
}

func (h H) String() string {
	return fmt.Sprintf("H %dx%d", h.Rows, h.Cols)
}
