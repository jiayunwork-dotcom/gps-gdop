package wgs84

import (
	"fmt"
	"math"
)

type ENU struct {
	East  float64
	North float64
	Up    float64
}

func NewENU(east, north, up float64) (ENU, error) {
	if err := RequireFinite(east, north, up); err != nil {
		return ENU{}, err
	}
	return ENU{East: east, North: north, Up: up}, nil
}

type ENUBasis struct {
	EastRow  ECEF
	NorthRow ECEF
	UpRow    ECEF
}

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

func BasisAtECEF(p ECEF, e Ellipsoid) (ENUBasis, error) {
	g, err := FromECEF(p, e)
	if err != nil {
		return ENUBasis{}, err
	}
	return BasisAt(g), nil
}

func (b ENUBasis) ToENU(v ECEF) ENU {
	return ENU{
		East:  b.EastRow.Dot(v),
		North: b.NorthRow.Dot(v),
		Up:    b.UpRow.Dot(v),
	}
}

func (b ENUBasis) ToECEF(v ENU) ECEF {
	e := b.EastRow.Scale(v.East)
	n := b.NorthRow.Scale(v.North)
	u := b.UpRow.Scale(v.Up)
	return e.Add(n).Add(u)
}

func (b ENUBasis) Matrix() Mat3 {
	return Mat3{
		M00: b.EastRow.X, M01: b.EastRow.Y, M02: b.EastRow.Z,
		M10: b.NorthRow.X, M11: b.NorthRow.Y, M12: b.NorthRow.Z,
		M20: b.UpRow.X, M21: b.UpRow.Y, M22: b.UpRow.Z,
	}
}

func (b ENUBasis) IsOrthonormal(tol float64) bool {
	m := b.Matrix()
	mt := m.Transpose()
	prod := m.Mul(mt)
	id := Identity3()
	return prod.EqualWithin(id, tol)
}

func (b ENUBasis) String() string {
	return fmt.Sprintf("ENU basis E=%v N=%v U=%v", b.EastRow, b.NorthRow, b.UpRow)
}
