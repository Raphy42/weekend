package std

func NewSlice[T any](items ...T) []T {
	return items
}

func Map[
	TIn any,
	TOut any,
](slice []TIn, fn func(idx int, in TIn, eOut *[]TOut)) []TOut {
	out := make([]TOut, len(slice))
	for idx, item := range slice {
		fn(idx, item, &out)
	}
	return out
}

func Filter[
	T any,
](slice []T, predicate func(idx int, value T) bool) []T {
	out := make([]T, 0)
	for idx, item := range slice {
		if predicate(idx, item) {
			out = append(out, item)
		}
	}
	return out
}

func UniqueFn[
	T any,
	P func(lhs, rhs T) bool,
](slice []T, predicate P) []T {
	out := make([]T, 0)
	for _, item := range slice {
		unique := true
		for _, item2 := range out {
			if predicate(item, item2) {
				unique = false
				break
			}
		}
		if unique {
			out = append(out, item)
		}
	}
	return out
}

func Unique[
	T comparable,
](slice []T) []T {
	return UniqueFn(slice, func(lhs, rhs T) bool {
		return lhs == rhs
	})
}

func Fold[
	TIn any,
	TOut any,
	A func(idx int, current TIn, acc TOut) TOut,
](slice []TIn, accumulator A, initialValue TOut) TOut {
	for idx, item := range slice {
		initialValue = accumulator(idx, item, initialValue)
	}
	return initialValue
}

func Flatten[
	T any,
](slices ...[]T) []T {
	return Fold(slices, func(_ int, item []T, acc []T) []T {
		return append(acc, item...)
	}, make([]T, 0))
}

func Find[
	T any,
	P func(item T) bool,
](slice []T, predicate P) *T {
	for _, item := range slice {
		if predicate(item) {
			return &item
		}
	}
	return nil
}
