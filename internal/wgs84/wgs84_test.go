package wgs84

import (
	"math"
	"testing"
)

func TestGeodeticRoundTrip(t *testing.T) {
	e := NewWGS84()
	cases := []LatLonHeight{
		{LatRad: 0, LonRad: 0, HeightM: 0},
		{LatRad: 0.7, LonRad: 2.1, HeightM: 500},
		{LatRad: -0.5, LonRad: -1.2, HeightM: 12345},
		{LatRad: 1.55, LonRad: -3.0, HeightM: -30},
	}
	for _, g := range cases {
		p := g.ToECEF(e)
		got, err := FromECEF(p, e)
		if err != nil {
			t.Fatalf("FromECEF(%v): %v", p, err)
		}
		if math.Abs(got.LatRad-g.LatRad) > 1e-10 {
			t.Errorf("lat round trip = %.12f, want %.12f", got.LatRad, g.LatRad)
		}
		if math.Abs(got.LonRad-g.LonRad) > 1e-12 {
			t.Errorf("lon round trip = %.12f, want %.12f", got.LonRad, g.LonRad)
		}
		if math.Abs(got.HeightM-g.HeightM) > 1e-4 {
			t.Errorf("height round trip = %.6f m, want %.6f m", got.HeightM, g.HeightM)
		}
	}
}

func TestKnownECEFAtEquator(t *testing.T) {
	e := NewWGS84()
	g, _ := NewLatLonHeightDeg(0, 0, 0)
	p := g.ToECEF(e)
	if math.Abs(p.X-e.SemiMajorAxis) > 1e-6 {
		t.Errorf("X = %.6f, want semi-major axis %.6f", p.X, e.SemiMajorAxis)
	}
	if math.Abs(p.Y) > 1e-9 || math.Abs(p.Z) > 1e-9 {
		t.Errorf("Y, Z = %.3g, %.3g, want 0, 0", p.Y, p.Z)
	}
}

func TestENUBasisOrthonormal(t *testing.T) {
	for _, g := range []LatLonHeight{
		{LatRad: 0, LonRad: 0, HeightM: 0},
		{LatRad: 0.8, LonRad: 1.9, HeightM: 250},
		{LatRad: -1.2, LonRad: -0.4, HeightM: 1000},
	} {
		b := BasisAt(g)
		if !b.IsOrthonormal(1e-12) {
			t.Errorf("basis at %v is not orthonormal: %v", g, b)
		}
	}
}

func TestKnownTopocentricDirection(t *testing.T) {
	e := NewWGS84()
	g, _ := NewLatLonHeightDeg(0, 0, 0)
	p := g.ToECEF(e)
	b, err := BasisAtECEF(p, e)
	if err != nil {
		t.Fatal(err)
	}
	direction := b.ToECEF(ENU{East: math.Sqrt(0.5), North: 0, Up: math.Sqrt(0.5)})
	los, err := direction.Unit()
	if err != nil {
		t.Fatal(err)
	}
	back := b.ToENU(los)
	if math.Abs(back.East-math.Sqrt(0.5)) > 1e-12 {
		t.Errorf("east component = %.12f, want %.12f", back.East, math.Sqrt(0.5))
	}
	if math.Abs(back.Up-math.Sqrt(0.5)) > 1e-12 {
		t.Errorf("up component = %.12f, want %.12f", back.Up, math.Sqrt(0.5))
	}
	if math.Abs(back.North) > 1e-12 {
		t.Errorf("north component = %.12f, want 0", back.North)
	}
}

func TestCovarianceRotationPreservesTrace(t *testing.T) {
	e := NewWGS84()
	g, _ := NewLatLonHeightDeg(39.9, 116.4, 50)
	p := g.ToECEF(e)
	b, err := BasisAtECEF(p, e)
	if err != nil {
		t.Fatal(err)
	}
	cov := Covariance3{
		M00: 12, M01: 1, M02: -2,
		M10: 1, M11: 8, M12: 0.5,
		M20: -2, M21: 0.5, M22: 6,
	}
	rotated := RotateCovariance(cov, b)
	if math.Abs(rotated.Trace()-cov.Trace()) > 1e-9 {
		t.Errorf("trace after rotation = %.12f, want %.12f", rotated.Trace(), cov.Trace())
	}
}

func TestECEFGeodeticResidualAndBasis(t *testing.T) {
	e := NewWGS84()
	g, err := NewLatLonHeightDeg(39.9, 116.4, 50)
	if err != nil {
		t.Fatal(err)
	}
	p := g.ToECEF(e)
	res, err := RoundTripResidual(p, e)
	if err != nil {
		t.Fatal(err)
	}
	if res > 1e-6 {
		t.Fatalf("ECEF round-trip residual %v m", res)
	}
	ok, err := BasisAgrees(p, e, 1e-12)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("ENU basis from ECEF must match basis from geodetic")
	}
	ok, err = UpAgreesWithGeodetic(p, e, 1e-12)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("local up must match geodetic normal")
	}
}
