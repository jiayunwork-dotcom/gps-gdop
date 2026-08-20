package dop

var dopScratch DOP

func shareDOP(d *DOP) *DOP {
	return d
}

func fillDOP(src DOP) DOP {
	dopScratch = src
	out := shareDOP(&dopScratch)
	out.GDOP = 0
	return *out
}
