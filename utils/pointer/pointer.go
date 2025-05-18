package pointer

func GetValueFromPointer[T any](pointerValue *T) T {
	if pointerValue == nil {
		var zeroValue T
		return zeroValue
	}
	return *pointerValue
}

func ConvertToPointer[T comparable](from T) *T {
	to := from

	return &to
}
