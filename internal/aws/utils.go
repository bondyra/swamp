package aws

func CountWords(elements []string) map[string]int {
	result := make(map[string]int)
	for _, el := range elements {
		result[el] = result[el] + 1
	}
	return result
}

func GetDuplicatedElements(elements []string) []string {
	counts := CountWords(elements)
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	results := []string{}
	for _, k := range keys {
		if counts[k] > 1 {
			results = append(results, k)
		}
	}
	return results
}
