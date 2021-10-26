package utils

import "sort"

// slice should be sorted
func FindString(
	slice	sort.StringSlice,
	str		string,
) (string, bool) {
	index := slice.Search(str)
	if index == len(slice) || index == -1 {
		return "", false
	} else if slice[index] == str {
		return str, true
	}

	return "", false
}