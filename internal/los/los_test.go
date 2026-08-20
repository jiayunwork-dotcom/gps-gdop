package los

import (
	"errors"
	"math"
	"testing"

	"gps-gdop/internal/wgs84"
)

func TestLOSUnitModulus(t *testing.T) {
	rx := Receiver{ECEF: wgs84.ECEF{X: 2148744.0, Y: 4426641.0, Z: 4044660.0}}
	sats := []wgs84.ECEF{
		{X: 26600000, Y: 12000000, Z: 8000000},
		{X: 9000000, Y: 25000000, Z: 10000000},
		{X: -15000000, Y: 20000000, Z: 12000000},
		{X: 12000000, Y: -8000000, Z: 23000000},
	}
	for i, p := range sats {
		s, err := NewSatellite(string(rune('A'+i))+string(rune('0'+i)), p)
		if err != nil {
			t.Fatal(err)
		}
		sight, err := ComputeSight(rx, s)
		if err != nil {
			t.Fatalf("sight %d: %v", i, err)
		}
		if sight.UnitModulus() > 1e-12 {
			t.Errorf("unit sight %d modulus deviates by %.3g, want <= 1e-12", i, sight.UnitModulus())
		}
	}
}

func TestLOSCoincidentRejected(t *testing.T) {
	rx := Receiver{ECEF: wgs84.ECEF{X: 1000000, Y: 2000000, Z: 3000000}}
	sat, _ := NewSatellite("G01", rx.ECEF)
	_, err := ComputeSight(rx, sat)
	if err == nil {
		t.Error("coincident satellite produced a sight, want error")
	}
	var se *SightError
	if !errors.As(err, &se) {
		t.Errorf("coincident error = %T, want *SightError", err)
	}
}

func TestElevationAzimuthKnown(t *testing.T) {
	e := wgs84.NewWGS84()
	g, _ := wgs84.NewLatLonHeightDeg(0, 0, 0)
	rxPos := g.ToECEF(e)
	basis, err := wgs84.BasisAtECEF(rxPos, e)
	if err != nil {
		t.Fatal(err)
	}
	rangeM := 20000000.0
	dir := basis.ToECEF(wgs84.ENU{East: math.Sqrt(0.5), North: 0, Up: math.Sqrt(0.5)})
	satPos := rxPos.Add(dir.Scale(rangeM))
	sat, _ := NewSatellite("G01", satPos)
	sight, err := ComputeSight(Receiver{ECEF: rxPos}, sat)
	if err != nil {
		t.Fatal(err)
	}
	el := Degrees(ElevationAngle(sight.Unit, basis))
	az := Degrees(AzimuthAngle(sight.Unit, basis))
	if math.Abs(el-45) > 1e-9 {
		t.Errorf("elevation = %.9f deg, want 45", el)
	}
	if math.Abs(az-90) > 1e-9 {
		t.Errorf("azimuth = %.9f deg, want 90", az)
	}
}

func TestMaskInvalidRejected(t *testing.T) {
	if _, err := NewMask(90); !errors.Is(err, ErrMaskTooHigh) {
		t.Errorf("mask 90 err = %v, want ErrMaskTooHigh", err)
	}
	if _, err := NewMask(120); !errors.Is(err, ErrMaskTooHigh) {
		t.Errorf("mask 120 err = %v, want ErrMaskTooHigh", err)
	}
	if _, err := NewMask(-1); !errors.Is(err, ErrMaskNegative) {
		t.Errorf("mask -1 err = %v, want ErrMaskNegative", err)
	}
	if _, err := NewMask(math.NaN()); !errors.Is(err, ErrMaskTooHigh) {
		t.Errorf("mask NaN err = %v, want ErrMaskTooHigh", err)
	}
}

func TestMaskFiltersLowElevation(t *testing.T) {
	e := wgs84.NewWGS84()
	g, _ := wgs84.NewLatLonHeightDeg(0, 0, 0)
	rxPos := g.ToECEF(e)
	basis, err := wgs84.BasisAtECEF(rxPos, e)
	if err != nil {
		t.Fatal(err)
	}
	rx := Receiver{ECEF: rxPos}
	rangeM := 20000000.0
	build := func(name string, enu wgs84.ENU) (LineOfSight, error) {
		dir := basis.ToECEF(enu)
		sat, err := NewSatellite(name, rxPos.Add(dir.Scale(rangeM)))
		if err != nil {
			return LineOfSight{}, err
		}
		return ComputeSight(rx, sat)
	}
	sights := []LineOfSight{}
	// 45 deg elevation, azimuth 0 (north)
	s1, _ := build("G01", wgs84.ENU{East: 0, North: math.Cos(Radians(45)), Up: math.Sin(Radians(45))})
	// 3 deg elevation, azimuth 90 (east) -> below the 5 deg mask
	s2, _ := build("G02", wgs84.ENU{East: math.Cos(Radians(3)), North: 0, Up: math.Sin(Radians(3))})
	// 60 deg elevation, azimuth 180 (south)
	s3, _ := build("G03", wgs84.ENU{East: 0, North: -math.Cos(Radians(60)), Up: math.Sin(Radians(60))})
	// 30 deg elevation, azimuth 270 (west)
	s4, _ := build("G04", wgs84.ENU{East: -math.Cos(Radians(30)), North: 0, Up: math.Sin(Radians(30))})
	sights = append(sights, s1, s2, s3, s4)

	mask, _ := NewMask(5)
	visible, rejected := Filter(sights, basis, mask)
	if len(visible) != 3 {
		t.Errorf("visible count = %d, want 3", len(visible))
	}
	if len(rejected) != 1 || rejected[0] != "G02" {
		t.Errorf("rejected = %v, want [G02]", rejected)
	}
}

func TestBuildHClockColumnOnes(t *testing.T) {
	rx := Receiver{ECEF: wgs84.ECEF{X: 2148744.0, Y: 4426641.0, Z: 4044660.0}}
	positions := []wgs84.ECEF{
		{X: 26600000, Y: 12000000, Z: 8000000},
		{X: 9000000, Y: 25000000, Z: 10000000},
		{X: -15000000, Y: 20000000, Z: 12000000},
		{X: 12000000, Y: -8000000, Z: 23000000},
		{X: 30000000, Y: -5000000, Z: -6000000},
	}
	var sights []LineOfSight
	for i, p := range positions {
		id := "G" + string(rune('0'+i+1))
		sat, _ := NewSatellite(id, p)
		s, err := ComputeSight(rx, sat)
		if err != nil {
			t.Fatalf("sight %s: %v", id, err)
		}
		sights = append(sights, s)
	}
	h, err := BuildH(sights)
	if err != nil {
		t.Fatal(err)
	}
	if h.Cols != 4 {
		t.Errorf("geometry matrix width = %d, want 4 (position + clock)", h.Cols)
	}
	if h.Rows != len(sights) {
		t.Errorf("geometry matrix height = %d, want %d", h.Rows, len(sights))
	}
	for i, v := range h.ClockColumn() {
		if v != 1.0 {
			t.Errorf("clock column entry %d = %v, want exactly 1.0", i, v)
		}
	}
}

func TestBuildHTooFewSights(t *testing.T) {
	rx := Receiver{ECEF: wgs84.ECEF{X: 0, Y: 0, Z: 0}}
	sats := []wgs84.ECEF{
		{X: 20000000, Y: 0, Z: 0},
		{X: 0, Y: 20000000, Z: 0},
		{X: 0, Y: 0, Z: 20000000},
	}
	var sights []LineOfSight
	for i, p := range sats {
		sat, _ := NewSatellite(string(rune('A'+i)), p)
		s, _ := ComputeSight(rx, sat)
		sights = append(sights, s)
	}
	if _, err := BuildH(sights); !errors.Is(err, ErrTooFewSights) {
		t.Errorf("BuildH with 3 sights err = %v, want ErrTooFewSights", err)
	}
}

func TestSkyViewRejectsWhenMaskLeavesFewerThanFour(t *testing.T) {
	e := wgs84.NewWGS84()
	g, _ := wgs84.NewLatLonHeightDeg(0, 0, 0)
	rxPos := g.ToECEF(e)
	basis, _ := wgs84.BasisAtECEF(rxPos, e)
	rx := Receiver{ECEF: rxPos}
	rangeM := 20000000.0
	build := func(name string, enu wgs84.ENU) LineOfSight {
		dir := basis.ToECEF(enu)
		sat, _ := NewSatellite(name, rxPos.Add(dir.Scale(rangeM)))
		s, _ := ComputeSight(rx, sat)
		return s
	}
	sights := []LineOfSight{
		build("G01", wgs84.ENU{East: 0, North: math.Cos(Radians(45)), Up: math.Sin(Radians(45))}),
		build("G02", wgs84.ENU{East: math.Cos(Radians(2)), North: 0, Up: math.Sin(Radians(2))}),
		build("G03", wgs84.ENU{East: 0, North: -math.Cos(Radians(60)), Up: math.Sin(Radians(60))}),
		build("G04", wgs84.ENU{East: -math.Cos(Radians(30)), North: 0, Up: math.Sin(Radians(30))}),
	}
	mask, _ := NewMask(5)
	_, err := BuildSkyView(sights, basis, mask)
	if !errors.Is(err, ErrTooFewVisible) {
		t.Errorf("sky view with 3 visible err = %v, want ErrTooFewVisible", err)
	}
}
