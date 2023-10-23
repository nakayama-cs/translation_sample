package main

func ConvType[T1 any, T2 any](array []T1, converter func(value T1) T2) []T2 {
	results := make([]T2, 0)
	for _, v := range array {
		results = append(results, converter(v))
	}
	return results
}
