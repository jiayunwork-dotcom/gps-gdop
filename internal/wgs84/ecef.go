package wgs84

import (
	"fmt"
	"math"
)

type ECEF struct {
	X float64
	Y float64
	Z float64
}

func NewECEF(x, y, z float64) (ECEF, error) {
	if err := RequireFinite(x, y, z); err != nil {
		return ECEF{}, err
	}
	return ECEF{X: x, Y: y, Z: z}, nil
}

func Origin() ECEF {
	return ECEF{}
}

func (p ECEF) Add(q ECEF) ECEF {
	return ECEF{X: p.X + q.X, Y: p.Y + q.Y, Z: p.Z + q.Z}
}

func (p ECEF) Sub(q ECEF) ECEF {
	return ECEF{X: p.X - q.X, Y: p.Y - q.Y, Z: p.Z - q.Z}
}

func (p ECEF) Scale(s float64) ECEF {
	return ECEF{X: p.X * s, Y: p.Y * s, Z: p.Z * s}
}

func (p ECEF) Dot(q ECEF) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

func (p ECEF) Cross(q ECEF) ECEF {
	return ECEF{
		X: p.Y*q.Z - p.Z*q.Y,
		Y: p.Z*q.X - p.X*q.Z,
		Z: p.X*q.Y - p.Y*q.X,
	}
}

func (p ECEF) Norm() float64 {
	return math.Sqrt(p.Dot(p))
}

func (p ECEF) NormSquared() float64 {
	return p.Dot(p)
}

func (p ECEF) Unit() (ECEF, error) {
	n := p.Norm()
	if n == 0 || !IsFinite(n) {
		return ECEF{}, ErrZeroNorm
	}
	return p.Scale(1.0 / n), nil
}

func (p ECEF) Distance(q ECEF) float64 {
	return p.Sub(q).Norm()
}

func (p ECEF) DistanceKm(q ECEF) float64 {
	return p.Distance(q) / 1000.0
}

func (p ECEF) String() string {
	return fmt.Sprintf("ECEF(%.6f, %.6f, %.6f)", p.X, p.Y, p.Z)
}

func (p ECEF) EqualWithin(q ECEF, tol float64) bool {
	return NearlyEqual(p.X, q.X, tol) &&
		NearlyEqual(p.Y, q.Y, tol) &&
		NearlyEqual(p.Z, q.Z, tol)
}
