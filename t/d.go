
func test(structure []int, test []int) (bool, int) {
	matches := false
	breakPoint := 0
	for i, val := range structure {
		matches = (val == test[i])
		if !matches {
			breakPoint = i
			break
		}
	}
	return matches, breakPoint
}

func test_struct_snippet(structure []int, test []Token) (bool, int) {

	matches := false
	breakPoint := 0
	for i, val := range structure {
		matches = (val == test[i].Type)
		if !matches {
			breakPoint = i
			break
		}
	}
	return matches, breakPoint
}