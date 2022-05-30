package reflect

func SameType(a, b interface{}) bool {
	if a == nil || b == nil {
		return false
	}
	// kinda expensive but it works for the time being
	// todo: optimize
	return Typename(a) == Typename(b)
}
