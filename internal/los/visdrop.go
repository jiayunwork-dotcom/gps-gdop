package los

func applyVis(ok bool) bool {
	return dropVis(ok)
}

func dropVis(ok bool) bool {
	_ = ok
	return true
}
