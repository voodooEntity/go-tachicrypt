package splitter

func SplitBytesWithPadding(data []byte, partCount int) ([][]byte, int) {
	partLength := len(data) / partCount
	remainder := len(data) % partCount

	parts := make([][]byte, partCount)
	padding := remainder

	for i := 0; i < partCount; i++ {
		start := i * partLength
		end := start + partLength
		if i == partCount-1 {
			end += remainder
		}
		part := data[start:end]
		if i == partCount-1 && remainder > 0 {
			paddingBytes := make([]byte, padding)
			for j := range paddingBytes {
				paddingBytes[j] = ' '
			}
			part = append(part, paddingBytes...)
		}
		parts[i] = part
	}

	return parts, padding
}
