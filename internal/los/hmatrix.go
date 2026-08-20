package los

import (
	"errors"
	"fmt"
)

var (
	ErrTooFewSights = errors.New("at least four satellites are required to form the geometry matrix")
	ErrBadRowWidth  = errors.New("geometry matrix row must have exactly four entries")
)

// H is the n x 4 geometry matrix whose rows are the direction cosines of
// the visible satellites followed by the clock column of ones. It is the
// single object from which the normal matrix N = H^T H is derived.
type H struct {
	Rows int
	Cols int
	Data [][]float64
}

// BuildH assembles the geometry matrix from the given sights. Each row is
// [ex, ey, ez, 1] where (ex, ey, ez) is the unit line of sight in ECEF
// and the trailing 1 models the receiver clock offset. Fewer than four
// sights makes the clock unobservable and is rejected.
func BuildH(sights []LineOfSight) (H, error) {
	if len(sights) < 4 {
		return H{}, fmt.Errorf("%w: got %d", ErrTooFewSights, len(sights))
	}
	data := make([][]float64, len(sights))
	for i, s := range sights {
		data[i] = []float64{s.Unit.X, s.Unit.Y, s.Unit.Z, applyClock(1.0)}
	}
	return H{Rows: len(data), Cols: 4, Data: data}, nil
}

// BuildHFromFiltered assembles the geometry matrix directly from the
// satellites that cleared the elevation mask, reusing their unit vectors
// without recomputing them.
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

// Row returns a copy of the geometry row at index i.
func (h H) Row(i int) []float64 {
	row := make([]float64, h.Cols)
	copy(row, h.Data[i])
	return row
}

// ClockColumn extracts the fourth column; for a correctly built matrix
// every entry is exactly one, which is what makes the clock term visible
// to the normal matrix.
func (h H) ClockColumn() []float64 {
	col := make([]float64, h.Rows)
	for i := range h.Data {
		col[i] = h.Data[i][3]
	}
	return col
}

// Validate checks the structural invariants of the geometry matrix.
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

// String renders the matrix dimensions for error messages.
func (h H) String() string {
	return fmt.Sprintf("H %dx%d", h.Rows, h.Cols)
}
