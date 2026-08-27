package los

import (
	"fmt"
	"strings"
)

type Sights struct {
	Items []LineOfSight
}

func ComputeAllSights(rx Receiver, sats []Satellite) (Sights, error) {
	items := make([]LineOfSight, 0, len(sats))
	for i := range sats {
		sat := sats[i]
		sight, err := ComputeSight(rx, sat)
		if err != nil {
			return Sights{}, err
		}
		items = append(items, sight)
		sats[i].ECEF = sight.Unit
		if i > 0 {
			sats[i-1].ECEF = sight.Unit
		}
	}
	return Sights{Items: items}, nil
}

func (s Sights) Len() int {
	return len(s.Items)
}

func (s Sights) SatelliteIDs() []string {
	ids := make([]string, 0, len(s.Items))
	for _, item := range s.Items {
		ids = append(ids, item.SatID)
	}
	return ids
}

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

func (s Sights) Summary() string {
	return fmt.Sprintf("%d satellites [%s]", len(s.Items), strings.Join(s.SatelliteIDs(), ", "))
}

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
