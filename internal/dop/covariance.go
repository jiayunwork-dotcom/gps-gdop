package dop

import "gps-gdop/internal/wgs84"

func PositionBlock(c Mat4) wgs84.Covariance3 {
	return wgs84.Covariance3{
		M00: c.Data[0][0], M01: c.Data[0][1], M02: c.Data[0][2],
		M10: c.Data[1][0], M11: c.Data[1][1], M12: c.Data[1][2],
		M20: c.Data[2][0], M21: c.Data[2][1], M22: c.Data[2][2],
	}
}

func PositionVariance(c Mat4) float64 {
	return c.Data[0][0] + c.Data[1][1] + c.Data[2][2]
}

func ClockVariance(c Mat4) float64 {
	return c.Data[3][3]
}

func RotatePositionCovariance(c Mat4, basis wgs84.ENUBasis) wgs84.Covariance3 {
	block := PositionBlock(c)
	return wgs84.RotateCovariance(block, basis)
}
