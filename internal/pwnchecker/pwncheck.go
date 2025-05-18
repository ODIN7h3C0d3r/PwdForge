package pwnchecker

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// CheckPasswordPwned checks if the given password has been pwned using HIBP API.
func CheckPasswordPwned(password string) (bool, int, error) {
	hash := sha1.New()
	_, err := io.WriteString(hash, password)
	if err != nil {
		return false, 0, err
	}
	hashStr := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	prefix := hashStr[:5]
	suffix := hashStr[5:]

	// Remove the space after /range/ in the URL
	url := "https://api.pwnedpasswords.com/range/" + prefix

	resp, err := http.Get(url)
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, err
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.Contains(line, suffix) {
			countStr := strings.TrimSpace(strings.Split(line, ":")[1])
			count, _ := strconv.Atoi(countStr)
			return true, count, nil
		}
	}

	return false, 0, nil
}
