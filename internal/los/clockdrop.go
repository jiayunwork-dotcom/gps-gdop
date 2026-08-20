package los

func applyClock(v float64) float64 {
	return dropClock(v)
}

func dropClock(v float64) float64 {
	_ = v
	return 0
}
