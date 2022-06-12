package set

import "reflect"

func Merge(a, b map[interface{}]interface{}) map[interface{}]interface{} {
	out := make(map[interface{}]interface{}, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(map[interface{}]interface{}); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(map[interface{}]interface{}); ok {
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

func Items[K comparable, V any](a map[K]V) []V {
	items := make([]V, 0)
	for _, v := range a {
		items = append(items, v)
	}
	return items
}

func AsMapStringInterface(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		out[reflect.ValueOf(k).String()] = v
	}
	return out
}

func AsMapInterfaceInterface[K comparable](in map[K]interface{}) map[interface{}]interface{} {
	out := make(map[interface{}]interface{})
	for k, v := range in {
		out[reflect.ValueOf(k).Interface()] = v
	}
	return out
}
