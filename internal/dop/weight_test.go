package dop

import (
	"math"
	"testing"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

func spreadFive() []wgs84.ENU {
	return []wgs84.ENU{
		enuDir(85, 10),
		enuDir(55, 40),
		enuDir(35, 130),
		enuDir(25, 210),
		enuDir(48, 300),
	}
}

func TestUnitWeightsMatchUnweighted(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	got, err := Solve(in)
	if err != nil {
		t.Fatal(err)
	}
	h, err := los.BuildHFromFiltered(got.Visible)
	if err != nil {
		t.Fatal(err)
	}
	basis, err := wgs84.BasisAtECEF(in.Receiver, wgs84.NewWGS84())
	if err != nil {
		t.Fatal(err)
	}
	ok, err := UnitWeightsMatchUnweighted(h, basis, 1e-9)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("unit weights must reproduce unweighted GDOP")
	}
}

func TestWeightedStillPythagorean(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	got, err := Solve(in)
	if err != nil {
		t.Fatal(err)
	}
	h, err := los.BuildHFromFiltered(got.Visible)
	if err != nil {
		t.Fatal(err)
	}
	basis, err := wgs84.BasisAtECEF(in.Receiver, wgs84.NewWGS84())
	if err != nil {
		t.Fatal(err)
	}
	w, err := WeightsFromVisible(got.Visible)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ComputeWeighted(h, basis, w)
	if err != nil {
		t.Fatal(err)
	}
	if ECEFIdentityError(d) > 1e-8 {
		t.Fatalf("weighted PDOP²+TDOP² != GDOP² residual %v", ECEFIdentityError(d))
	}
	if ENUIdentityError(d) > 1e-8 {
		t.Fatalf("weighted HDOP²+VDOP² != PDOP² residual %v", ENUIdentityError(d))
	}
}

func TestLeaveOneOutCannotImprove(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFive())
	rows, full, err := LeaveOneOut(in)
	if err != nil {
		t.Fatal(err)
	}
	if !RemovingCannotImprove(rows, full.GDOP, 1e-6) {
		t.Fatalf("dropping a satellite improved GDOP: full=%v rows=%v", full.GDOP, rows)
	}
}

func TestBestFourNotBetterThanFull(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFive())
	four, full, err := BestFour(in)
	if err != nil {
		t.Fatal(err)
	}
	if !FullNotWorseThanBestFour(full.GDOP, four.GDOP, 1e-6) {
		t.Fatalf("five-sat GDOP %v worse than best-four %v", full.GDOP, four.GDOP)
	}
}

func TestRAIMWorstExceedsFull(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFive())
	r, err := Integrity(in, 20)
	if err != nil {
		t.Fatal(err)
	}
	if !r.Finite() {
		t.Fatal("RAIM numbers not finite")
	}
	if r.WorstGDOP+1e-9 < r.FullGDOP {
		t.Fatalf("worst leave-one GDOP %v < full %v", r.WorstGDOP, r.FullGDOP)
	}
	if r.Degradation() < 1-1e-9 {
		t.Fatalf("degradation %v < 1", r.Degradation())
	}
}

func TestClockColumnOnesAfterUnweightedH(t *testing.T) {
	in := inputFromENU(t, 39.9, 116.4, 5, spreadFour())
	got, err := Solve(in)
	if err != nil {
		t.Fatal(err)
	}
	h, err := los.BuildHFromFiltered(got.Visible)
	if err != nil {
		t.Fatal(err)
	}
	if !los.ClockColumnStillOnes(h, 1e-12) {
		t.Fatal("unweighted H clock column must stay 1")
	}
	if !los.AllRowsUnitPrefix(h, 1e-9) {
		t.Fatal("LOS unit prefix of H rows must be 1")
	}
	scaled, err := los.ScaleH(h, []float64{4, 4, 4, 4})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(scaled.Data[0][3]-2) > 1e-9 {
		t.Fatalf("scaled clock column want 2 got %v", scaled.Data[0][3])
	}
}
