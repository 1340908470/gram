package base

// GetFollow 根据推导式的左部，得到其对应的Follow集
func GetFollow(left Tag) []Tag {
	_, err := GetProductionsByTag(GetProductions(), left)
	if err != nil {
		panic(err)
	}

	return nil
}
