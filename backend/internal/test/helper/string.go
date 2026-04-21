package helper

func StringOfLength(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}
