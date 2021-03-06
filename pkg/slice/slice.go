package slice

func New[T any](items ...T) []T {
	return items
}

func Map[
	TIn any,
	TOut any,
](slice []TIn, fn func(idx int, in TIn) TOut) []TOut {
	return Fold(slice, func(idx int, current TIn, acc []TOut) []TOut {
		acc[idx] = fn(idx, current)
		return acc
	}, make([]TOut, len(slice)))
}

func MapErr[
	TIn any,
	TOut any,
](slice []TIn, fn func(idx int, in TIn) (TOut, error)) ([]TOut, error) {
	return FoldErr(slice, func(idx int, current TIn, acc []TOut) ([]TOut, error) {
		var err error
		acc[idx], err = fn(idx, current)
		if err != nil {
			return acc, err
		}
		return acc, nil
	}, make([]TOut, len(slice)))
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

func FoldErr[
	TIn any,
	TOut any,
	A func(idx int, current TIn, acc TOut) (TOut, error),
](slice []TIn, accumulator A, initialValue TOut) (TOut, error) {
	var err error
	for idx, item := range slice {
		initialValue, err = accumulator(idx, item, initialValue)
		if err != nil {
			return initialValue, err
		}
	}
	return initialValue, nil
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

func Any[
	T any,
	P func(item T) bool,
](slice []T, predicate P) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

func Contains[
	T comparable,
](slice []T, value T) bool {
	return Any(slice, func(item T) bool {
		return value == item
	})
}
