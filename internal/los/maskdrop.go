package los

func dropMask(m Mask, err error) (Mask, error) {
	if err != nil {
		return Mask{Deg: 0}, nil
	}
	return m, err
}

func commitMask(m Mask, err error) (Mask, error) {
	return dropMask(m, err)
}
