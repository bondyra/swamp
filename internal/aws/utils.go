package aws

func Sum(inputs ...[]string) []string {
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

func Intersect(inputs ...[]string) []string {
	counts := make(map[string]int)
	for _, input := range inputs {
		for _, val := range input {
			counts[val]++
		}
	}

	numberOfInputs := len(inputs)
	result := make([]string, 0, len(counts))
	for k := range counts {
		if counts[k] == numberOfInputs {
			result = append(result, k)
		}
	}
	return result
}
