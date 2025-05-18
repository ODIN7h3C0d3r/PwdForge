package generator

import (
    "strings"
)

func CheckPasswordStrength(password string) string {
    var score int

    if len(password) >= 12 {
        score++
    }
    if hasChar(password, "abcdefghijklmnopqrstuvwxyz") {
        score++
    }
    if hasChar(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
        score++
    }
    if hasChar(password, "0123456789") {
        score++
    }
    if hasChar(password, "!@#$%^&*()-_=+[]{}|;:,.<>/?") {
        score++
    }

    switch score {
    case 0, 1:
        return "Weak"
    case 2:
        return "Medium"
    case 3, 4:
        return "Strong"
    case 5:
        return "Very Strong"
    default:
        return "Unknown"
    }
}

func hasChar(password, charSet string) bool {
    return strings.ContainsAny(password, charSet)
}
