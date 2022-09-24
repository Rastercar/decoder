package arrays

func ConcatAppend[T any](slices [][]T) []T {
	var tmp []T
	for _, s := range slices {
		tmp = append(tmp, s...)
	}
	return tmp
}
