package dop

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

const gpsOrbitRadius = 26560000.0

func enuDir(elDeg, azDeg float64) wgs84.ENU {
	el := elDeg * math.Pi / 180.0
	az := azDeg * math.Pi / 180.0
	return wgs84.ENU{
		East:  math.Cos(el) * math.Sin(az),
		North: math.Cos(el) * math.Cos(az),
		Up:    math.Sin(el),
	}
}

func inputFromENU(t *testing.T, latDeg, lonDeg float64, mask float64, dirs []wgs84.ENU) SolveInput {
	t.Helper()
	e := wgs84.NewWGS84()
	g, err := wgs84.NewLatLonHeightDeg(latDeg, lonDeg, 0)
	if err != nil {
		t.Fatal(err)
	}
	rxPos := g.ToECEF(e)
	basis, err := wgs84.BasisAtECEF(rxPos, e)
	if err != nil {
		t.Fatal(err)
	}
	sats := make([]los.Satellite, 0, len(dirs))
	for i, dir := range dirs {
		satPos := rxPos.Add(basis.ToECEF(dir).Scale(gpsOrbitRadius))
		sat, err := los.NewSatellite(fmt.Sprintf("G%02d", i+1), satPos)
		if err != nil {
			t.Fatal(err)
		}
		sats = append(sats, sat)
	}
	return SolveInput{Receiver: rxPos, Satellites: sats, MaskDeg: mask}
}

func spreadFour() []wgs84.ENU {
	return []wgs84.ENU{
		enuDir(90, 0),
		enuDir(19.47, 0),
		enuDir(19.47, 120),
		enuDir(19.47, 240),
	}
}

func poorFour() []wgs84.ENU {
	return []wgs84.ENU{
		enuDir(6, 350),
		enuDir(8, 5),
		enuDir(9, 15),
		enuDir(7, 25),
	}
}

func TestDopLessThanFourRejected(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour()[:3])
	_, err := Solve(in)
	if !errors.Is(err, ErrTooFewSatellites) {
		t.Errorf("three satellites err = %v, want ErrTooFewSatellites", err)
	}
}

func TestDopPythagoreanIdentity(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	res, err := Solve(in)
	if err != nil {
		t.Fatal(err)
	}
	errAbs := ECEFIdentityError(res.Report.AsDOP())
	if errAbs > 1e-9 {
		t.Errorf("|PDOP^2+TDOP^2-GDOP^2| = %.3g, want <= 1e-9 (gdop %.6f pdop %.6f tdop %.6f)",
			errAbs, res.Report.GDOP, res.Report.PDOP, res.Report.TDOP)
	}
	enuAbs := ENUIdentityError(res.Report.AsDOP())
	if enuAbs > 1e-9 {
		t.Errorf("|HDOP^2+VDOP^2-PDOP^2| = %.3g, want <= 1e-9 (pdop %.6f hdop %.6f vdop %.6f)",
			enuAbs, res.Report.PDOP, res.Report.HDOP, res.Report.VDOP)
	}
}

func TestDopGoodVsPoorConstellation(t *testing.T) {
	goodIn := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	good, err := Solve(goodIn)
	if err != nil {
		t.Fatal(err)
	}
	poorIn := inputFromENU(t, 39.9, 116.4, 5, poorFour())
	poor, err := Solve(poorIn)
	if err != nil {
		t.Fatal(err)
	}
	if good.Report.GDOP >= poor.Report.GDOP {
		t.Errorf("spread four-sat GDOP = %.6f, want strictly below poor four-sat GDOP = %.6f",
			good.Report.GDOP, poor.Report.GDOP)
	}
}

func TestDopSingularCoplanar(t *testing.T) {
	coplanar := []wgs84.ENU{
		enuDir(30, 0),
		enuDir(30, 90),
		enuDir(30, 180),
		enuDir(30, 270),
	}
	in := inputFromENU(t, 39.9, 116.4, 0, coplanar)
	_, err := Solve(in)
	if !errors.Is(err, ErrSingularNormal) {
		t.Errorf("coplanar constellation err = %v, want ErrSingularNormal", err)
	}
}

func TestDopAddSatelliteNoWorse(t *testing.T) {
	base := spreadFour()
	in := inputFromENU(t, 39.9, 116.4, 5, base)
	res4, err := Solve(in)
	if err != nil {
		t.Fatal(err)
	}
	extraDir := enuDir(70, 180)
	in5 := inputFromENU(t, 39.9, 116.4, 5, append(append([]wgs84.ENU{}, base...), extraDir))
	res5, err := Solve(in5)
	if err != nil {
		t.Fatal(err)
	}
	if res5.Report.GDOP > res4.Report.GDOP+1e-9 {
		t.Errorf("GDOP after adding a non-coplanar satellite = %.9f, want not worse than %.9f",
			res5.Report.GDOP, res4.Report.GDOP)
	}
}

func TestDopComparisonHelpers(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, []wgs84.ENU{
		enuDir(55, 30),
		enuDir(45, 140),
		enuDir(60, 250),
		enuDir(50, 320),
		enuDir(12, 60),
		enuDir(70, 180),
	})
	maskCmp, err := CompareMasks(in, 25)
	if err != nil {
		t.Fatal(err)
	}
	if maskCmp.Improved {
		t.Errorf("raising the mask must not improve GDOP: %s", maskCmp.String())
	}
	if maskCmp.HighUsed >= maskCmp.LowUsed {
		t.Errorf("raising the mask must drop satellites: used %d -> %d", maskCmp.LowUsed, maskCmp.HighUsed)
	}
	satCmp, err := CompareWithExtraSatellite(inputFromENU(t, 39.9, 116.4, 5, spreadFour()), extraSatellite(t))
	if err != nil {
		t.Fatal(err)
	}
	if satCmp.Worsened {
		t.Errorf("adding a satellite must not worsen GDOP: %s", satCmp.String())
	}
}

func extraSatellite(t *testing.T) los.Satellite {
	t.Helper()
	e := wgs84.NewWGS84()
	g, _ := wgs84.NewLatLonHeightDeg(39.9, 116.4, 0)
	rxPos := g.ToECEF(e)
	basis, _ := wgs84.BasisAtECEF(rxPos, e)
	dir := enuDir(75, 45)
	satPos := rxPos.Add(basis.ToECEF(dir).Scale(gpsOrbitRadius))
	sat, err := los.NewSatellite("GEX", satPos)
	if err != nil {
		t.Fatal(err)
	}
	return sat
}

func TestDopCutoffMonotonic(t *testing.T) {
	seven := []wgs84.ENU{
		enuDir(55, 30),
		enuDir(45, 140),
		enuDir(60, 250),
		enuDir(50, 320),
		enuDir(10, 60),
		enuDir(12, 200),
		enuDir(70, 90),
	}
	inLow := inputFromENU(t, 39.9, 116.4, 5, seven)
	resLow, err := Solve(inLow)
	if err != nil {
		t.Fatal(err)
	}
	inHigh := inputFromENU(t, 39.9, 116.4, 25, seven)
	resHigh, err := Solve(inHigh)
	if err != nil {
		t.Fatal(err)
	}
	if resLow.Report.UsedSats <= resHigh.Report.UsedSats {
		t.Fatalf("raising the mask must drop satellites: used %d -> %d",
			resLow.Report.UsedSats, resHigh.Report.UsedSats)
	}
	if resHigh.Report.GDOP < resLow.Report.GDOP-1e-9 {
		t.Errorf("GDOP at mask 25 = %.9f, must not improve below GDOP at mask 5 = %.9f when satellites were removed",
			resHigh.Report.GDOP, resLow.Report.GDOP)
	}
}

func TestDopReceiverMicroShift(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	shifts := []wgs84.ECEF{
		{X: 5, Y: 0, Z: 0},
		{X: -5, Y: 0, Z: 0},
		{X: 0, Y: 5, Z: 0},
		{X: 0, Y: 0, Z: 5},
		{X: -20, Y: 10, Z: 5},
	}
	sens, err := ReceiverShiftSensitivity(in, shifts)
	if err != nil {
		t.Fatal(err)
	}
	if sens.MaxChange >= 1e-3 {
		t.Errorf("metre-level receiver shift changed GDOP relatively by %.3g, want far below constellation-scale changes (>= 1e-3)", sens.MaxChange)
	}
}

func TestDopConditionNumberReported(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	res, err := Solve(in)
	if err != nil {
		t.Fatal(err)
	}
	if !(res.Report.Condition > 0) || !res.Report.AllFinite() {
		t.Errorf("condition number = %.6g, want finite and positive", res.Report.Condition)
	}
}
