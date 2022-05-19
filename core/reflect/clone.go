package reflect

func Clone(value interface{}) interface{} {
	return &value
}
