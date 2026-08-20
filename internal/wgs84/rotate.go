package wgs84

// Covariance3 is a symmetric 3x3 covariance matrix stored row-major.
type Covariance3 struct {
	M00, M01, M02 float64
	M10, M11, M12 float64
	M20, M21, M22 float64
}

// FromMat3 wraps a Mat3 as a covariance, keeping the symmetry contract
// explicit at the type level.
func FromMat3(m Mat3) Covariance3 {
	return Covariance3{
		M00: m.M00, M01: m.M01, M02: m.M02,
		M10: m.M10, M11: m.M11, M12: m.M12,
		M20: m.M20, M21: m.M21, M22: m.M22,
	}
}

// ToMat3 exposes the covariance as a plain matrix for arithmetic.
func (c Covariance3) ToMat3() Mat3 {
	return Mat3{
		M00: c.M00, M01: c.M01, M02: c.M02,
		M10: c.M10, M11: c.M11, M12: c.M12,
		M20: c.M20, M21: c.M21, M22: c.M22,
	}
}

// RotateCovariance performs the similarity transform C' = R C R^T.
// It converts an ECEF covariance into the local ENU frame of the given
// basis while preserving the total variance (trace invariance).
func RotateCovariance(c Covariance3, basis ENUBasis) Covariance3 {
	cm := c.ToMat3()
	rot := basis.Matrix()
	t := cm.Mul(rot.Transpose())
	result := rot.Mul(t)
	return FromMat3(result)
}

// Trace is the total variance of the distribution.
func (c Covariance3) Trace() float64 {
	return c.M00 + c.M11 + c.M22
}

// HorizontalVariance is the East plus North variance, whose square root
// is the horizontal dilution HDOP.
func (c Covariance3) HorizontalVariance() float64 {
	return c.M00 + c.M11
}

// VerticalVariance is the Up variance, whose square root is VDOP.
func (c Covariance3) VerticalVariance() float64 {
	return c.M22
}
