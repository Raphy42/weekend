package reflect

import "reflect"

func Typename(value any) string {
	return reflect.TypeOf(value).String()
}
