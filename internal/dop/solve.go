package dop

import (
	"errors"
	"fmt"

	"gps-gdop/internal/los"
	"gps-gdop/internal/wgs84"
)

var (
	ErrTooFewSatellites = errors.New("at least four satellites are required")
)

type SolveInput struct {
	Receiver   wgs84.ECEF
	Satellites []los.Satellite
	MaskDeg    float64
}

type SolveResult struct {
	Report   Report
	Visible  []los.FilteredSight
	Rejected []string
}

func Solve(in SolveInput) (SolveResult, error) {
	if err := in.validate(); err != nil {
		return SolveResult{}, err
	}
	mask, err := los.NewMask(in.MaskDeg)
	if err != nil {
		return SolveResult{}, err
	}
	basis, err := wgs84.BasisAtECEF(in.Receiver, wgs84.NewWGS84())
	if err != nil {
		return SolveResult{}, err
	}
	rx, err := los.NewReceiver(in.Receiver)
	if err != nil {
		return SolveResult{}, err
	}
	sights, err := los.ComputeAllSights(rx, in.Satellites)
	if err != nil {
		return SolveResult{}, err
	}
	visible, rejected := los.Filter(sights.Items, basis, mask)
	if len(visible) < 4 {
		return SolveResult{}, fmt.Errorf("%w: only %d satellites clear the %.1f degree mask",
			los.ErrTooFewVisible, len(visible), mask.Deg)
	}
	h, err := los.BuildHFromFiltered(visible)
	if err != nil {
		return SolveResult{}, err
	}
	d, err := Compute(h, basis)
	if err != nil {
		return SolveResult{}, err
	}
	n, err := NormalMatrix(h)
	if err != nil {
		return SolveResult{}, err
	}
	inv, err := Invert4(n)
	if err != nil {
		return SolveResult{}, err
	}
	cond := ConditionNumber(n, inv)
	report := FromDOP(d, len(visible), sights.Len(), len(rejected), cond, mask.Deg)
	return SolveResult{
		Report:   report,
		Visible:  visible,
		Rejected: rejected,
	}, nil
}

func (in SolveInput) validate() error {
	if len(in.Satellites) < 4 {
		return fmt.Errorf("%w: got %d satellites", ErrTooFewSatellites, len(in.Satellites))
	}
	if err := in.Receiver.ValidateFinite(); err != nil {
		return err
	}
	for _, sat := range in.Satellites {
		if err := sat.ECEF.ValidateFinite(); err != nil {
			return err
		}
		if sat.ID == "" {
			return los.ErrEmptySatelliteID
		}
	}
	return nil
}
