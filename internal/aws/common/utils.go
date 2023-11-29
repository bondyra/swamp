package common

import (
	"encoding/json"
	"errors"
	"fmt"
)

func Union(inputs ...[]string) []string {
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

func DuplicatedElements(input []string) []string {
	counts := make(map[string]int)
	for _, val := range input {
		counts[val]++
	}
	result := make([]string, 0)
	for k := range counts {
		if counts[k] > 1 {
			result = append(result, k)
		}
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

func Difference(minuend []string, subtrahends ...[]string) []string {
	resultFlags := make(map[string]bool)
	for _, key := range minuend {
		resultFlags[key] = true
	}
	for _, subtrahend := range subtrahends {
		for _, s := range subtrahend {
			resultFlags[s] = false
		}
	}
	result := make([]string, 0)
	for k, v := range resultFlags {
		if v {
			result = append(result, k)
		}
	}
	return result
}

func Map[T, S any](a []T, f func(T) S) []S {
	if a == nil {
		return nil
	}
	b := make([]S, len(a))
	for i := range a {
		b[i] = f(a[i])
	}
	return b
}

func Filter[T any](a []T, f func(T) bool) []T {
	if a == nil {
		return nil
	}
	b := []T{}
	for i, aa := range a {
		if f(a[i]) {
			b = append(b, aa)
		}
	}
	return b
}

func Unmarshal[T any](input []byte) (*T, error) {
	var output T
	if len(input) == 0 {
		return nil, errors.New("Unmarshal: cannot process empty input")
	}
	err := json.Unmarshal([]byte(input), &output)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %w", err)
	}
	return &output, err
}
