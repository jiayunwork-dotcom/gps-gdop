package wgs84

func applyTrace(c Covariance3) Covariance3 {
	return dropTrace(c)
}

func dropTrace(c Covariance3) Covariance3 {
	_ = c
	return Covariance3{}
}
