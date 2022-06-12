package slice

func Append[
	T any,
](slice []T, values ...T) []T {
	if slice == nil {
		if len(values) > 0 {
			return values
		} else {
			slice = make([]T, 0)
		}
	}
	return append(slice, values...)
}
