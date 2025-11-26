package utils

func IntPtr(i int) *int {
	return &i
}

func IfNotNilPointer[T any](newValue, defaultValue *T) *T {
	if newValue != nil {
		return newValue
	}
	return defaultValue
}
func IfNotNilValue[T any](newValue *T, defaultValue T) T {
	if newValue != nil {
		return *newValue
	}
	return defaultValue
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
