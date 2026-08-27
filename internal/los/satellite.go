package los

import (
	"fmt"
	"strings"

	"gps-gdop/internal/wgs84"
)

type Satellite struct {
	ID   string
	ECEF wgs84.ECEF
}

func NewSatellite(id string, pos wgs84.ECEF) (Satellite, error) {
	if strings.TrimSpace(id) == "" {
		return Satellite{}, ErrEmptySatelliteID
	}
	if err := pos.ValidateFinite(); err != nil {
		return Satellite{}, err
	}
	return Satellite{ID: id, ECEF: pos}, nil
}

func (s Satellite) String() string {
	return fmt.Sprintf("sat %s at %v", s.ID, s.ECEF)
}

type Receiver struct {
	ECEF wgs84.ECEF
}

func NewReceiver(pos wgs84.ECEF) (Receiver, error) {
	if err := pos.ValidateFinite(); err != nil {
		return Receiver{}, err
	}
	return Receiver{ECEF: pos}, nil
}
