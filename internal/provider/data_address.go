package provider

import (
	"strings"
)

const hexDigits = "0123456789abcdef"

func completeIPv4Address(input string) string {
	switch strings.Count(input, ".") {
	case 3:
		return input
	case 2:
		return input + ".0"
	case 1:
		return input + ".0.0"
	case 0:
		return input + ".0.0.0"
	default:
		return input
	}
}
