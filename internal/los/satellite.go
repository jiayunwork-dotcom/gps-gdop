package los

import (
	"fmt"
	"strings"

	"gps-gdop/internal/wgs84"
)

// Satellite is a navigation satellite whose ECEF position is supplied by
// the caller. The ID is an opaque label used for reporting which space
// vehicles actually entered the dilution computation.
type Satellite struct {
	ID   string
	ECEF wgs84.ECEF
}

// NewSatellite validates the identifier and the coordinates.
func NewSatellite(id string, pos wgs84.ECEF) (Satellite, error) {
	if strings.TrimSpace(id) == "" {
		return Satellite{}, ErrEmptySatelliteID
	}
	if err := pos.ValidateFinite(); err != nil {
		return Satellite{}, err
	}
	return Satellite{ID: id, ECEF: pos}, nil
}

// String renders the satellite for error messages.
func (s Satellite) String() string {
	return fmt.Sprintf("sat %s at %v", s.ID, s.ECEF)
}

// Receiver is the approximate user position, given in ECEF metres.
type Receiver struct {
	ECEF wgs84.ECEF
}

// NewReceiver validates the receiver coordinates.
func NewReceiver(pos wgs84.ECEF) (Receiver, error) {
	if err := pos.ValidateFinite(); err != nil {
		return Receiver{}, err
	}
	return Receiver{ECEF: pos}, nil
}
