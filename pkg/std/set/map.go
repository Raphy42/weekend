package set

import "reflect"

func Merge(a, b map[any]any) map[any]any {
	out := make(map[any]any, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[any]any); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[any]any); ok {
					out[k] = Merge(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func Contains[K comparable, V any](values map[K]V, key K) bool {
	_, ok := values[key]
	return ok
}

func Keys[K comparable, V any](a map[K]V) []K {
	keys := make([]K, 0)
	for k := range a {
		keys = append(keys, k)
	}
	return keys
}

func Values[K comparable, V any](a map[K]V) []V {
	items := make([]V, 0)
	for _, v := range a {
		items = append(items, v)
	}
	return items
}

func CollectMap[K comparable, T, V any](slice []T, fn func(item T) (K, V)) map[K]V {
	set := make(map[K]V)
	for _, item := range slice {
		k, v := fn(item)
		set[k] = v
	}
	return set
}

func Clone[K comparable, V any](values map[K]V) map[K]V {
	out := make(map[K]V)
	for k, v := range values {
		out[k] = v
	}
	return out
}

func AsMapStringInterface(in map[any]any) map[string]any {
	out := make(map[string]any)
	for k, v := range in {
		out[reflect.ValueOf(k).String()] = v
	}
	return out
}

func AsMapInterfaceInterface[K comparable](in map[K]any) map[any]any {
	out := make(map[any]any)
	for k, v := range in {
		out[reflect.ValueOf(k).Interface()] = v
	}
	return out
}
