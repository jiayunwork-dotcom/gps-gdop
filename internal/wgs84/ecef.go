package wgs84

import (
	"fmt"
	"math"
)

// ECEF is a point or vector in the Earth-Centred Earth-Fixed frame.
// Coordinates are metres along the WGS84 axes: X through the Greenwich
// meridian in the equatorial plane, Z along the rotation axis.
type ECEF struct {
	X float64
	Y float64
	Z float64
}

// NewECEF validates the three components and returns a fresh point.
func NewECEF(x, y, z float64) (ECEF, error) {
	if err := RequireFinite(x, y, z); err != nil {
		return ECEF{}, err
	}
	return ECEF{X: x, Y: y, Z: z}, nil
}

// Origin is the geocentric origin of the ECEF frame.
func Origin() ECEF {
	return ECEF{}
}

// Add returns the vector sum of two points.
func (p ECEF) Add(q ECEF) ECEF {
	return ECEF{X: p.X + q.X, Y: p.Y + q.Y, Z: p.Z + q.Z}
}

// Sub returns the displacement from q to p (p - q).
func (p ECEF) Sub(q ECEF) ECEF {
	return ECEF{X: p.X - q.X, Y: p.Y - q.Y, Z: p.Z - q.Z}
}

// Scale multiplies every component by s.
func (p ECEF) Scale(s float64) ECEF {
	return ECEF{X: p.X * s, Y: p.Y * s, Z: p.Z * s}
}

// Dot is the Euclidean inner product of two vectors.
func (p ECEF) Dot(q ECEF) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

// Cross is the right-handed vector product p x q.
func (p ECEF) Cross(q ECEF) ECEF {
	return ECEF{
		X: p.Y*q.Z - p.Z*q.Y,
		Y: p.Z*q.X - p.X*q.Z,
		Z: p.X*q.Y - p.Y*q.X,
	}
}

// Norm is the Euclidean length of the vector.
func (p ECEF) Norm() float64 {
	return math.Sqrt(p.Dot(p))
}

// NormSquared is the squared Euclidean length, avoiding a square root.
func (p ECEF) NormSquared() float64 {
	return p.Dot(p)
}

// Unit returns the vector divided by its length. It returns an error when
// the length is zero so callers can reject degenerate geometry instead of
// producing a NaN direction.
func (p ECEF) Unit() (ECEF, error) {
	n := p.Norm()
	if n == 0 || !IsFinite(n) {
		return ECEF{}, ErrZeroNorm
	}
	return p.Scale(1.0 / n), nil
}

// Distance returns the straight-line separation between two points in metres.
func (p ECEF) Distance(q ECEF) float64 {
	return p.Sub(q).Norm()
}

// DistanceKm is the same separation expressed in kilometres.
func (p ECEF) DistanceKm(q ECEF) float64 {
	return p.Distance(q) / 1000.0
}

// String renders the vector for reports and error messages.
func (p ECEF) String() string {
	return fmt.Sprintf("ECEF(%.6f, %.6f, %.6f)", p.X, p.Y, p.Z)
}

// EqualWithin reports whether two points agree to within tol metres in
// every component.
func (p ECEF) EqualWithin(q ECEF, tol float64) bool {
	return NearlyEqual(p.X, q.X, tol) &&
		NearlyEqual(p.Y, q.Y, tol) &&
		NearlyEqual(p.Z, q.Z, tol)
}
