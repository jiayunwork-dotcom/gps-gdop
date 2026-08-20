package wgs84

import "errors"

var (
	ErrEllipsoidNotFinite   = errors.New("ellipsoid parameters must be finite")
	ErrEllipsoidDegenerate  = errors.New("ellipsoid semi-axes must be positive")
	ErrEllipsoidMinorAxisLonger = errors.New("semi-minor axis must not exceed semi-major axis")
	ErrEllipsoidFlattening  = errors.New("flattening must lie strictly between 0 and 1")
	ErrCoordinateNotFinite  = errors.New("coordinate value must be finite")
	ErrZeroNorm             = errors.New("cannot normalize a vector with zero length")
)

// IsFinite reports whether v is neither NaN nor infinite.
func IsFinite(v float64) bool {
	return !mathIsNaN(v) && !mathIsInf(v)
}

func mathIsNaN(v float64) bool {
	return v != v
}

func mathIsInf(v float64) bool {
	return v > maxFloat || v < -maxFloat
}

const maxFloat = 1.7976931348623157e308

// RequireFinite returns an error if any of the given values is not finite.
func RequireFinite(values ...float64) error {
	for _, v := range values {
		if !IsFinite(v) {
			return ErrCoordinateNotFinite
		}
	}
	return nil
}

// NearlyEqual compares two scalars with an absolute tolerance.
func NearlyEqual(a, b, tol float64) bool {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d <= tol
}
