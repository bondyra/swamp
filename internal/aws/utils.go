package aws

func Intersect(inputs ...[]string) []string {
	keys := make(map[string]bool)
	for _, input := range inputs {
		for _, val := range input {
			keys[val] = true
		}
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
}
