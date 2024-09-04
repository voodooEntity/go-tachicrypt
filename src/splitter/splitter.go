package splitter

import (
	"strings"
)

func SplitStringWithPadding(str string, partCount int) ([]string, int) {
	partLength := len(str) / partCount
	remainder := len(str) % partCount

	parts := make([]string, partCount)
	padding := remainder

	for i := 0; i < partCount; i++ {
		start := i * partLength
		end := start + partLength
		if i == partCount-1 {
			end += remainder
		}
		part := str[start:end]
		if i == partCount-1 && remainder > 0 {
			part += strings.Repeat(" ", padding)
		}
		parts[i] = part
	}

	return parts, padding
}
