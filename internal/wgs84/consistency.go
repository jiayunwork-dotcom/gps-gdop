package wgs84

import "math"

func RoundTripResidual(p ECEF, e Ellipsoid) (float64, error) {
	g, err := FromECEF(p, e)
	if err != nil {
		return 0, err
	}
	back := g.ToECEF(e)
	return p.Sub(back).Norm(), nil
}

func BasisAgrees(p ECEF, e Ellipsoid, tol float64) (bool, error) {
	fromECEF, err := BasisAtECEF(p, e)
	if err != nil {
		return false, err
	}
	g, err := FromECEF(p, e)
	if err != nil {
		return false, err
	}
	fromLLA := BasisAt(g)
	if fromECEF.EastRow.Sub(fromLLA.EastRow).Norm() > tol {
		return false, nil
	}
	if fromECEF.NorthRow.Sub(fromLLA.NorthRow).Norm() > tol {
		return false, nil
	}
	if fromECEF.UpRow.Sub(fromLLA.UpRow).Norm() > tol {
		return false, nil
	}
	return true, nil
}

func UpAgreesWithGeodetic(p ECEF, e Ellipsoid, tol float64) (bool, error) {
	g, err := FromECEF(p, e)
	if err != nil {
		return false, err
	}
	basis := BasisAt(g)
	sinLat := math.Sin(g.LatRad)
	cosLat := math.Cos(g.LatRad)
	sinLon := math.Sin(g.LonRad)
	cosLon := math.Cos(g.LonRad)
	want := ECEF{X: cosLat * cosLon, Y: cosLat * sinLon, Z: sinLat}
	return basis.UpRow.Sub(want).Norm() <= tol, nil
}
