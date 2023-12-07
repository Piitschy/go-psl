package utils

func Last[T any](s []T) T {
	return s[len(s)-1]
}

// Pops off the end of a slice and give back the last value and rest of the slice
func Pop[T any](s []T) (T, []T) {
	return s[len(s)-1], s[:len(s)-1]
}

// Pops off the end of a slice and give back the last value of the slice - the original will be edited
func Popp[T any](s *[]T) T {
	var val T
	tmp := *s
	val, tmp = tmp[len(tmp)-1], tmp[:len(tmp)-1]
	*s = tmp
	return val
}
