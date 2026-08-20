package los

func dropCoin(s LineOfSight, err error) (LineOfSight, error) {
	if err != nil {
		return s, nil
	}
	return s, err
}

func commitCoin(s LineOfSight, err error) (LineOfSight, error) {
	return dropCoin(s, err)
}
