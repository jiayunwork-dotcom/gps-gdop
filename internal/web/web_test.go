package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"gps-gdop/internal/wgs84"
)

func testServer() http.Handler {
	return NewServer(Assets{
		WebFS: nil,
		Examples: map[string][]byte{
			"four-good": []byte(`{"receiver_ecef":{"x":-2178657.082724949,"y":4388876.23355147,"z":4069505.7479817173},"satellites":[{"id":"G01","x":-2148953.7,"y":4388876.2,"z":30351114.8},{"id":"G02","x":-3398115.4,"y":24924566.8,"z":3416720.4},{"id":"G03","x":-10441238.1,"y":-9806172.5,"z":3416720.4},{"id":"G04","x":-1133260.1,"y":-30125314.3,"z":3416720.4}],"elevation_mask_deg":5}`),
		},
	})
}

func enuToJSON(t *testing.T, latDeg, lonDeg float64, dirs []wgs84.ENU) DopRequest {
	t.Helper()
	e := wgs84.NewWGS84()
	g, err := wgs84.NewLatLonHeightDeg(latDeg, lonDeg, 0)
	if err != nil {
		t.Fatal(err)
	}
	rxPos := g.ToECEF(e)
	basis, err := wgs84.BasisAtECEF(rxPos, e)
	if err != nil {
		t.Fatal(err)
	}
	req := DopRequest{ReceiverECEF: ECEFPoint{X: rxPos.X, Y: rxPos.Y, Z: rxPos.Z}}
	mask := 5.0
	req.ElevationMaskDeg = &mask
	for i, dir := range dirs {
		satPos := rxPos.Add(basis.ToECEF(dir).Scale(26560000.0))
		req.Satellites = append(req.Satellites, SatelliteDTO{
			ID: fmt.Sprintf("G%02d", i+1),
			X:  satPos.X,
			Y:  satPos.Y,
			Z:  satPos.Z,
		})
	}
	return req
}

func postJSON(t *testing.T, h http.Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	r.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

func TestAPIDopSuccess(t *testing.T) {
	h := testServer()
	req := enuToJSON(t, 39.9, 116.4, []wgs84.ENU{
		{East: 0, North: 0, Up: 1},
		{East: 0, North: math.Cos(19.47 * math.Pi / 180), Up: math.Sin(19.47 * math.Pi / 180)},
		{East: math.Cos(19.47 * math.Pi / 180) * math.Sin(120 * math.Pi / 180), North: math.Cos(19.47 * math.Pi / 180) * math.Cos(120 * math.Pi / 180), Up: math.Sin(19.47 * math.Pi / 180)},
		{East: math.Cos(19.47 * math.Pi / 180) * math.Sin(240 * math.Pi / 180), North: math.Cos(19.47 * math.Pi / 180) * math.Cos(240 * math.Pi / 180), Up: math.Sin(19.47 * math.Pi / 180)},
	})
	rec := postJSON(t, h, "/api/dop", req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out DopResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.GDOP < 2.0 || out.GDOP > 3.0 {
		t.Errorf("GDOP = %.6f, want in [2.0, 3.0] for the well-spread four-sat constellation", out.GDOP)
	}
	if out.SatellitesUsed != 4 {
		t.Errorf("satellites_used = %d, want 4", out.SatellitesUsed)
	}
	identity := math.Abs(out.PDOP*out.PDOP + out.TDOP*out.TDOP - out.GDOP*out.GDOP)
	if identity > 1e-9 {
		t.Errorf("|PDOP^2+TDOP^2-GDOP^2| = %.3g in API response, want <= 1e-9", identity)
	}
}

func TestAPIDopFewerThanFour(t *testing.T) {
	h := testServer()
	req := enuToJSON(t, 39.9, 116.4, []wgs84.ENU{
		{East: 0, North: 0, Up: 1},
		{East: 0, North: 1, Up: 0},
		{East: 1, North: 0, Up: 0},
	})
	rec := postJSON(t, h, "/api/dop", req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatal(err)
	}
	if errResp.Error == "" {
		t.Error("error body has empty error field")
	}
}

func TestAPIDopCoincidentSatellite(t *testing.T) {
	h := testServer()
	req := enuToJSON(t, 39.9, 116.4, []wgs84.ENU{
		{East: 0, North: 0, Up: 1},
		{East: 0, North: 1, Up: 0},
		{East: 1, North: 0, Up: 0},
		{East: 0, North: 0, Up: 1},
	})
	req.Satellites[3] = SatelliteDTO{ID: "G04", X: req.ReceiverECEF.X, Y: req.ReceiverECEF.Y, Z: req.ReceiverECEF.Z}
	rec := postJSON(t, h, "/api/dop", req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestAPIDopInvalidMask(t *testing.T) {
	h := testServer()
	req := enuToJSON(t, 39.9, 116.4, []wgs84.ENU{
		{East: 0, North: 0, Up: 1},
		{East: 0, North: 1, Up: 0},
		{East: 1, North: 0, Up: 0},
		{East: 0, North: 1, Up: 1},
	})
	mask := 90.0
	req.ElevationMaskDeg = &mask
	rec := postJSON(t, h, "/api/dop", req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestAPISkySuccess(t *testing.T) {
	h := testServer()
	req := enuToJSON(t, 0, 0, []wgs84.ENU{
		{East: 0, North: 0, Up: 1},
		{East: math.Cos(45 * math.Pi / 180), North: 0, Up: math.Sin(45 * math.Pi / 180)},
		{East: math.Cos(30 * math.Pi / 180), North: 0, Up: math.Sin(30 * math.Pi / 180)},
		{East: 0, North: -math.Cos(60 * math.Pi / 180), Up: math.Sin(60 * math.Pi / 180)},
	})
	rec := postJSON(t, h, "/api/sky", req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var out SkyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if len(out.Satellites) != 4 {
		t.Errorf("sky satellites = %d, want 4", len(out.Satellites))
	}
	// The east satellite must come back with azimuth 90 and elevation 45.
	var eastFound bool
	for _, p := range out.Satellites {
		if p.ID == "G02" {
			eastFound = true
			if math.Abs(p.AzimuthDeg-90) > 1e-9 {
				t.Errorf("azimuth of east satellite = %.9f, want 90", p.AzimuthDeg)
			}
			if math.Abs(p.ElevationDeg-45) > 1e-9 {
				t.Errorf("elevation of east satellite = %.9f, want 45", p.ElevationDeg)
			}
		}
	}
	if !eastFound {
		t.Error("east satellite G02 missing from sky response")
	}
}

func TestAPIMethodNotAllowed(t *testing.T) {
	h := testServer()
	r := httptest.NewRequest(http.MethodGet, "/api/dop", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestAPIBadJSON(t *testing.T) {
	h := testServer()
	r := httptest.NewRequest(http.MethodPost, "/api/dop", bytes.NewBufferString("{not json"))
	r.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}
