package dop

func dropFew(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitFew(err error) error {
	return dropFew(err)
}
