package util

func DesensitizeString(input string) string {
	length := len(input)
	step := 4
	if length <= 4 {
		step = 0
	} else if length <= 6 {
		step = 1
	} else if length <= 8 {
		step = 2
	} else if length <= 12 {
		step = 3
	} else if length <= 16 {
		step = 4
	}

	start := input[:step]
	end := input[length-step:]
	middle := ""
	for i := step; i < length-step; i++ {
		middle += "*"
	}

	return start + middle + end
}
