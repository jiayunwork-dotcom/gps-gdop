package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleGood() Record {
	return Record{
		Label: "four-good", GDOP: 2.35, PDOP: 2.0, TDOP: 1.2,
		HDOP: 1.4, VDOP: 1.4, UsedSats: 4, TotalSats: 4, Condition: 12, ElevationMask: 5,
	}
}

func samplePoor() Record {
	return Record{
		Label: "four-poor", GDOP: 111, PDOP: 90, TDOP: 40,
		HDOP: 70, VDOP: 55, UsedSats: 4, TotalSats: 4, Condition: 800, ElevationMask: 5,
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	rec := sampleGood()
	dir := t.TempDir()
	path := filepath.Join(dir, "dop.json")
	if err := WriteFile(path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Matches(rec, 1e-9) {
		t.Fatalf("mismatch %+v vs %+v", got, rec)
	}
}

func TestRejectsTooFewSats(t *testing.T) {
	rec := sampleGood()
	rec.UsedSats = 3
	dir := t.TempDir()
	if err := WriteFile(filepath.Join(dir, "bad.json"), rec); err == nil {
		t.Fatal("expected error")
	}
}

func TestEmptyFileRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path); err == nil {
		t.Fatal("expected error")
	}
}

func TestRankByGDOP(t *testing.T) {
	ranked, err := RankByGDOP([]Record{samplePoor(), sampleGood()})
	if err != nil {
		t.Fatal(err)
	}
	if ranked.Best.Label != "four-good" || ranked.Worst.Label != "four-poor" {
		t.Fatalf("rank %+v", ranked)
	}
	if ranked.Delta <= 0 {
		t.Fatal("delta must be positive")
	}
}

func TestAddingSatelliteDoesNotWorsen(t *testing.T) {
	before := sampleGood()
	after := sampleGood()
	after.UsedSats = 5
	after.GDOP = 2.1
	if err := AddingSatelliteDoesNotWorsen(before, after, 1e-9); err != nil {
		t.Fatal(err)
	}
}

func TestIdentityHoldsOnPythagorean(t *testing.T) {
	rec := Record{
		Label: "id", GDOP: 5, PDOP: 4, TDOP: 3,
		HDOP: 3.2, VDOP: 2.4, UsedSats: 4, TotalSats: 4,
	}
	if !IdentityHolds(rec, 1e-9) {
		t.Fatal("4-3-5 ECEF identity should hold; ENU may need H^2+V^2=P^2")
	}
}

func TestPoorExceedsGood(t *testing.T) {
	if err := PoorExceedsGood(sampleGood(), samplePoor(), 10); err != nil {
		t.Fatal(err)
	}
}

func TestRaiseMaskWorsens(t *testing.T) {
	low := sampleGood()
	high := sampleGood()
	high.ElevationMask = 20
	high.UsedSats = 4
	high.GDOP = 3.1
	if err := RaiseMaskWorsens(low, high, 1e-9); err != nil {
		t.Fatal(err)
	}
}

func TestWALReplayAfterTruncatedTail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dop.wal")
	if err := AppendWAL(path, sampleGood()); err != nil {
		t.Fatal(err)
	}
	if err := AppendWAL(path, samplePoor()); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	cut := len(raw) - 12
	if cut < 1 {
		t.Fatal("wal too small")
	}
	if err := TruncateWALTail(path, cut); err != nil {
		t.Fatal(err)
	}
	last, err := LastCommitted(path)
	if err != nil {
		t.Fatal(err)
	}
	if !last.Matches(sampleGood(), 1e-9) {
		t.Fatalf("truncated tail polluted prefix: %+v", last)
	}
}
