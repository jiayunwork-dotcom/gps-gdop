package wgs84

type Mat3 struct {
	M00, M01, M02 float64
	M10, M11, M12 float64
	M20, M21, M22 float64
}

func Identity3() Mat3 {
	return Mat3{
		M00: 1, M01: 0, M02: 0,
		M10: 0, M11: 1, M12: 0,
		M20: 0, M21: 0, M22: 1,
	}
}

func (m Mat3) Transpose() Mat3 {
	return Mat3{
		M00: m.M00, M01: m.M10, M02: m.M20,
		M10: m.M01, M11: m.M11, M12: m.M21,
		M20: m.M02, M21: m.M12, M22: m.M22,
	}
}

func (m Mat3) Mul(n Mat3) Mat3 {
	return Mat3{
		M00: m.M00*n.M00 + m.M01*n.M10 + m.M02*n.M20,
		M01: m.M00*n.M01 + m.M01*n.M11 + m.M02*n.M21,
		M02: m.M00*n.M02 + m.M01*n.M12 + m.M02*n.M22,
		M10: m.M10*n.M00 + m.M11*n.M10 + m.M12*n.M20,
		M11: m.M10*n.M01 + m.M11*n.M11 + m.M12*n.M21,
		M12: m.M10*n.M02 + m.M11*n.M12 + m.M12*n.M22,
		M20: m.M20*n.M00 + m.M21*n.M10 + m.M22*n.M20,
		M21: m.M20*n.M01 + m.M21*n.M11 + m.M22*n.M21,
		M22: m.M20*n.M02 + m.M21*n.M12 + m.M22*n.M22,
	}
}

func (m Mat3) Trace() float64 {
	return m.M00 + m.M11 + m.M22
}

func (m Mat3) Det() float64 {
	return m.M00*(m.M11*m.M22-m.M12*m.M21) -
		m.M01*(m.M10*m.M22-m.M12*m.M20) +
		m.M02*(m.M10*m.M21-m.M11*m.M20)
}

func (m Mat3) EqualWithin(n Mat3, tol float64) bool {
	return NearlyEqual(m.M00, n.M00, tol) &&
		NearlyEqual(m.M01, n.M01, tol) &&
		NearlyEqual(m.M02, n.M02, tol) &&
		NearlyEqual(m.M10, n.M10, tol) &&
		NearlyEqual(m.M11, n.M11, tol) &&
		NearlyEqual(m.M12, n.M12, tol) &&
		NearlyEqual(m.M20, n.M20, tol) &&
		NearlyEqual(m.M21, n.M21, tol) &&
		NearlyEqual(m.M22, n.M22, tol)
}
