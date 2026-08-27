package snapshot

import "fmt"

func RaiseMaskWorsens(lowMask, highMask Record, tol float64) error {
	if highMask.ElevationMask <= lowMask.ElevationMask {
		return fmt.Errorf("snapshot: higher mask record must have larger mask")
	}
	if highMask.UsedSats > lowMask.UsedSats {
		return fmt.Errorf("snapshot: raising mask must not increase used satellites")
	}
	if highMask.GDOP+tol < lowMask.GDOP {
		return fmt.Errorf("snapshot: raising mask must not improve GDOP: %g -> %g", lowMask.GDOP, highMask.GDOP)
	}
	return nil
}

func AddingSatelliteDoesNotWorsen(before, after Record, tol float64) error {
	if after.UsedSats <= before.UsedSats {
		return fmt.Errorf("snapshot: added satellite must increase used count")
	}
	if after.GDOP > before.GDOP+tol {
		return fmt.Errorf("snapshot: extra satellite worsened GDOP: %g -> %g", before.GDOP, after.GDOP)
	}
	return nil
}

func PoorExceedsGood(good, poor Record, minRatio float64) error {
	if good.GDOP <= 0 || poor.GDOP <= 0 {
		return fmt.Errorf("snapshot: both GDOP values must be positive")
	}
	if poor.GDOP < good.GDOP*minRatio {
		return fmt.Errorf("snapshot: poor GDOP %g is not at least %g times good %g", poor.GDOP, minRatio, good.GDOP)
	}
	return nil
}
