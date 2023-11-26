package aws

import "golang.org/x/exp/slices"

func RemoveDuplicates(input []string) []string {
	keys := make(map[string]bool)
	for _, val := range input {
		keys[val] = true
	}

	result := make([]string, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
}

func Intersect(input1, input2 []string) []string {
	result := []string{}
	for _, i := range input1 {
		if slices.Contains(input2, i) {
			result = append(result, i)
		}
	}
	return result
}
