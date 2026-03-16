package utils

import (
	"regexp"
)

func CheckPasswordStrength(password string) string {
	length := len(password)
	if length < 6 {
		return "weak"
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(password)
	score := 0
	if hasLower {
		score++
	}
	if hasUpper {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}
	switch score {
	case 1:
		return "weak"
	case 2:
		return "medium"
	case 3:
		return "strong"
	case 4:
		return "strong"
	default:
		return "weak"
	}
}
