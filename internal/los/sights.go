package los

import (
	"fmt"
	"strings"
)

// Sights is the ordered collection of line-of-sight directions for a
// constellation as seen from one receiver.
type Sights struct {
	Items []LineOfSight
}

// ComputeAllSights builds the sight for every satellite against the same
// receiver. A satellite that coincides with the receiver stops the whole
// constellation because its direction cannot be defined.
func ComputeAllSights(rx Receiver, sats []Satellite) (Sights, error) {
	items := make([]LineOfSight, 0, len(sats))
	for _, sat := range sats {
		sight, err := ComputeSight(rx, sat)
		if err != nil {
			return Sights{}, err
		}
		items = append(items, sight)
	}
	return Sights{Items: items}, nil
}

// Len returns the number of sights in the collection.
func (s Sights) Len() int {
	return len(s.Items)
}

// SatelliteIDs lists the identifiers in constellation order.
func (s Sights) SatelliteIDs() []string {
	ids := make([]string, 0, len(s.Items))
	for _, item := range s.Items {
		ids = append(ids, item.SatID)
	}
	return ids
}

// MaxUnitDeviation returns the largest deviation of any unit direction
// from length one, which must stay at machine precision for every row.
func (s Sights) MaxUnitDeviation() float64 {
	max := 0.0
	for _, item := range s.Items {
		d := item.UnitModulus()
		if d > max {
			max = d
		}
	}
	return max
}

// Summary renders the constellation in one line for reports.
func (s Sights) Summary() string {
	return fmt.Sprintf("%d satellites [%s]", len(s.Items), strings.Join(s.SatelliteIDs(), ", "))
}

// FromFiltered rebuilds a Sights value from the filtered subset, keeping
// the unit directions computed earlier.
func FromFiltered(visible []FilteredSight) Sights {
	items := make([]LineOfSight, 0, len(visible))
	for _, f := range visible {
		items = append(items, LineOfSight{
			SatID:  f.SatID,
			Unit:   f.Unit,
			RangeM: f.RangeM,
		})
	}
	return Sights{Items: items}
}
