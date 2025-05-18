package generator

import (
	"math"
	"strings"
)

// CheckPasswordStrength returns entropy and suggestions in addition to strength
func CheckPasswordStrength(password string) (string, float64, []string) {
	var score int
	var suggestions []string
	entropy := 0.0

	if len(password) >= 12 {
		score++
	} else {
		suggestions = append(suggestions, "Use at least 12 characters.")
	}
	if hasChar(password, "abcdefghijklmnopqrstuvwxyz") {
		score++
	} else {
		suggestions = append(suggestions, "Add lowercase letters.")
	}
	if hasChar(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		score++
	} else {
		suggestions = append(suggestions, "Add uppercase letters.")
	}
	if hasChar(password, "0123456789") {
		score++
	} else {
		suggestions = append(suggestions, "Add digits.")
	}
	if hasChar(password, "!@#$%^&*()-_=+[]{}|;:,.<>/?") {
		score++
	} else {
		suggestions = append(suggestions, "Add special characters.")
	}

	// Entropy estimation: log2(pool size^length)
	pool := 0
	if hasChar(password, "abcdefghijklmnopqrstuvwxyz") {
		pool += 26
	}
	if hasChar(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		pool += 26
	}
	if hasChar(password, "0123456789") {
		pool += 10
	}
	if hasChar(password, "!@#$%^&*()-_=+[]{}|;:,.<>/?") {
		pool += 32
	}
	if pool > 0 {
		entropy = float64(len(password)) * math.Log2(float64(pool))
	}

	strength := "Unknown"
	switch score {
	case 0, 1:
		strength = "Weak"
	case 2:
		strength = "Medium"
	case 3, 4:
		strength = "Strong"
	case 5:
		strength = "Very Strong"
	}
	return strength, entropy, suggestions
}

func hasChar(password, charSet string) bool {
	return strings.ContainsAny(password, charSet)
}
