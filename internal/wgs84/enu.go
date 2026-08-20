package wgs84

import (
	"fmt"
	"math"
)

// ENU is a topocentric vector with components East, North and Up.
// Up is the outward ellipsoid normal, North points toward the rotation
// axis, East completes the right-handed frame.
type ENU struct {
	East  float64
	North float64
	Up    float64
}

// NewENU validates the three components.
func NewENU(east, north, up float64) (ENU, error) {
	if err := RequireFinite(east, north, up); err != nil {
		return ENU{}, err
	}
	return ENU{East: east, North: north, Up: up}, nil
}

// ENUBasis is the orthonormal rotation that maps ECEF vectors to the
// local topocentric frame at a geodetic position. Its rows are the East,
// North and Up unit vectors expressed in ECEF coordinates.
type ENUBasis struct {
	EastRow  ECEF
	NorthRow ECEF
	UpRow    ECEF
}

// BasisAt builds the topocentric frame for the receiver geodetic position
// on the shared WGS84 ellipsoid.
func BasisAt(g LatLonHeight) ENUBasis {
	sinLat := math.Sin(g.LatRad)
	cosLat := math.Cos(g.LatRad)
	sinLon := math.Sin(g.LonRad)
	cosLon := math.Cos(g.LonRad)
	return ENUBasis{
		EastRow:  ECEF{X: -sinLon, Y: cosLon, Z: 0},
		NorthRow: ECEF{X: -sinLat * cosLon, Y: -sinLat * sinLon, Z: cosLat},
		UpRow:    ECEF{X: cosLat * cosLon, Y: cosLat * sinLon, Z: sinLat},
	}
}

// BasisAtECEF converts the receiver point to geodetic and builds the frame.
// The geodetic conversion is exact for points on the ellipsoid, so the
// resulting basis is consistent with every other WGS84 transform here.
func BasisAtECEF(p ECEF, e Ellipsoid) (ENUBasis, error) {
	g, err := FromECEF(p, e)
	if err != nil {
		return ENUBasis{}, err
	}
	return BasisAt(g), nil
}

// ToENU rotates an ECEF vector into the local frame.
func (b ENUBasis) ToENU(v ECEF) ENU {
	return ENU{
		East:  b.EastRow.Dot(v),
		North: b.NorthRow.Dot(v),
		Up:    b.UpRow.Dot(v),
	}
}

// ToECEF rotates a local vector back into ECEF using the transpose of
// the orthonormal rotation.
func (b ENUBasis) ToECEF(v ENU) ECEF {
	e := b.EastRow.Scale(v.East)
	n := b.NorthRow.Scale(v.North)
	u := b.UpRow.Scale(v.Up)
	return e.Add(n).Add(u)
}

// Matrix returns the 3x3 rotation whose rows are the basis unit vectors.
func (b ENUBasis) Matrix() Mat3 {
	return Mat3{
		M00: b.EastRow.X, M01: b.EastRow.Y, M02: b.EastRow.Z,
		M10: b.NorthRow.X, M11: b.NorthRow.Y, M12: b.NorthRow.Z,
		M20: b.UpRow.X, M21: b.UpRow.Y, M22: b.UpRow.Z,
	}
}

// IsOrthonormal reports whether the basis really is a rotation, used by
// tests to pin the shared WGS84 construction.
func (b ENUBasis) IsOrthonormal(tol float64) bool {
	m := b.Matrix()
	mt := m.Transpose()
	prod := m.Mul(mt)
	id := Identity3()
	return prod.EqualWithin(id, tol)
}

// String renders the three basis rows.
func (b ENUBasis) String() string {
	return fmt.Sprintf("ENU basis E=%v N=%v U=%v", b.EastRow, b.NorthRow, b.UpRow)
}
