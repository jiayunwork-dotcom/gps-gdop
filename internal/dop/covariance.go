package dop

import "gps-gdop/internal/wgs84"

// PositionBlock extracts the 3x3 upper-left submatrix of the covariance
// matrix, which contains the position uncertainty in the ECEF frame.
func PositionBlock(c Mat4) wgs84.Covariance3 {
	return wgs84.Covariance3{
		M00: c.Data[0][0], M01: c.Data[0][1], M02: c.Data[0][2],
		M10: c.Data[1][0], M11: c.Data[1][1], M12: c.Data[1][2],
		M20: c.Data[2][0], M21: c.Data[2][1], M22: c.Data[2][2],
	}
}

// PositionVariance sums the three diagonal entries of the position block.
func PositionVariance(c Mat4) float64 {
	return c.Data[0][0] + c.Data[1][1] + c.Data[2][2]
}

// ClockVariance returns the receiver clock entry of the covariance.
func ClockVariance(c Mat4) float64 {
	return c.Data[3][3]
}

// RotatePositionCovariance converts the ECEF position covariance into the
// local ENU frame using the shared WGS84 basis of the receiver. The total
// variance is preserved because the rotation is orthonormal.
func RotatePositionCovariance(c Mat4, basis wgs84.ENUBasis) wgs84.Covariance3 {
	block := PositionBlock(c)
	return wgs84.RotateCovariance(block, basis)
}
