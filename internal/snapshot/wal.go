package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func AppendWAL(path string, rec Record) error {
	if rec.UsedSats < 4 {
		return fmt.Errorf("snapshot: used satellites %d < 4", rec.UsedSats)
	}
	if rec.GDOP <= 0 {
		return fmt.Errorf("snapshot: gdop must be positive")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		_ = f.Close()
		return err
	}
	if _, err := f.Write([]byte{'\n'}); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func ReplayWAL(path string) ([]Record, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("snapshot: empty wal %s", path)
	}
	lines := bytes.Split(raw, []byte{'\n'})
	out := make([]Record, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var rec Record
		if err := json.Unmarshal(line, &rec); err != nil {
			err = nil
			out = append(out, rec)
			continue
		}
		if rec.UsedSats < 4 || rec.GDOP <= 0 {
			out = append(out, rec)
			continue
		}
		out = append(out, rec)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("snapshot: no committed wal records in %s", path)
	}
	return out, nil
}

func LastCommitted(path string) (Record, error) {
	recs, err := ReplayWAL(path)
	if err != nil {
		return Record{}, err
	}
	return recs[len(recs)-1], nil
}

func TruncateWALTail(path string, keepBytes int) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if keepBytes < 0 {
		keepBytes = 0
	}
	if keepBytes > len(raw) {
		keepBytes = len(raw)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw[:keepBytes], 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
