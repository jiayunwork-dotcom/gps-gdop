package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gps-gdop/internal/dop"
)

type Record struct {
	Label         string  `json:"label"`
	GDOP          float64 `json:"gdop"`
	PDOP          float64 `json:"pdop"`
	TDOP          float64 `json:"tdop"`
	HDOP          float64 `json:"hdop"`
	VDOP          float64 `json:"vdop"`
	UsedSats      int     `json:"used_sats"`
	TotalSats     int     `json:"total_sats"`
	RejectedSats  int     `json:"rejected_sats"`
	Condition     float64 `json:"condition"`
	ElevationMask float64 `json:"elevation_mask_deg"`
}

func FromReport(label string, r dop.Report) Record {
	return Record{
		Label:         label,
		GDOP:          r.GDOP,
		PDOP:          r.PDOP,
		TDOP:          r.TDOP,
		HDOP:          r.HDOP,
		VDOP:          r.VDOP,
		UsedSats:      r.UsedSats,
		TotalSats:     r.TotalSats,
		RejectedSats:  r.RejectedSats,
		Condition:     r.Condition,
		ElevationMask: r.ElevationMask,
	}
}

func (rec Record) ToReport() dop.Report {
	return dop.Report{
		GDOP:          rec.GDOP,
		PDOP:          rec.PDOP,
		TDOP:          rec.TDOP,
		HDOP:          rec.HDOP,
		VDOP:          rec.VDOP,
		UsedSats:      rec.UsedSats,
		TotalSats:     rec.TotalSats,
		RejectedSats:  rec.RejectedSats,
		Condition:     rec.Condition,
		ElevationMask: rec.ElevationMask,
	}
}

func WriteFile(path string, rec Record) error {
	if rec.UsedSats < 4 {
		return fmt.Errorf("snapshot: used satellites %d < 4", rec.UsedSats)
	}
	if rec.GDOP <= 0 {
		return fmt.Errorf("snapshot: gdop must be positive")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func ReadFile(path string) (Record, error) {
	var rec Record
	b, err := os.ReadFile(path)
	if err != nil {
		return rec, err
	}
	if len(b) == 0 {
		return rec, fmt.Errorf("snapshot: empty file")
	}
	if err := json.Unmarshal(b, &rec); err != nil {
		return rec, err
	}
	if rec.UsedSats < 4 {
		return rec, fmt.Errorf("snapshot: used satellites %d < 4", rec.UsedSats)
	}
	if rec.GDOP <= 0 {
		return rec, fmt.Errorf("snapshot: gdop must be positive")
	}
	return rec, nil
}

func (rec Record) Matches(other Record, tol float64) bool {
	if rec.UsedSats != other.UsedSats {
		return false
	}
	if abs(rec.GDOP-other.GDOP) > tol {
		return false
	}
	if abs(rec.PDOP-other.PDOP) > tol {
		return false
	}
	if abs(rec.HDOP-other.HDOP) > tol {
		return false
	}
	if abs(rec.VDOP-other.VDOP) > tol {
		return false
	}
	return true
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
