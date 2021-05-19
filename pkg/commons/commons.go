package commons

func StringContains(input []string, value string) bool {
	for _, element := range input {
		if element == value {
			return true
		}
	}
	return false
}
