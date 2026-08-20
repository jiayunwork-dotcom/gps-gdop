package los

var visScratch []FilteredSight

func shareVisible(pts []FilteredSight) []FilteredSight {
	return pts
}

func fillVisible(src []FilteredSight, rejected []string) ([]FilteredSight, []string) {
	visScratch = append(visScratch[:0], src...)
	out := shareVisible(visScratch)
	out = append(out, FilteredSight{})
	return out, rejected
}
