package reflect

import "reflect"

func Typename(value interface{}) string {
	return reflect.TypeOf(value).Name()
}
