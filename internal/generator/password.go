package generator

import (
    "crypto/rand"
    "encoding/binary"
    "io"
    "math"
    "os"
)

type PasswordConfig struct {
    Length          int
    Count           int
    IncludeUpper    bool
    IncludeLower    bool
    IncludeDigits   bool
    IncludeSpecials bool
    ExcludeSimilar  bool
}

const (
    lowerChars = "abcdefghijklmnopqrstuvwxyz"
    upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    digitChars = "0123456789"
    specialChars = "!@#$%^&*()-_=+[]{}|;:,.<>/?"
)

var reader io.Reader = rand.Reader

func GeneratePasswords(config PasswordConfig) []string {
    var charset string

    if config.IncludeLower {
        charset += lowerChars
    }
    if config.IncludeUpper {
        charset += upperChars
    }
    if config.IncludeDigits {
        charset += digitChars
    }
    if config.IncludeSpecials {
        charset += specialChars
    }

    if config.ExcludeSimilar {
        charset = removeSimilar(charset)
    }

    if len(charset) == 0 {
        panic("no character set selected for password generation")
    }

    rand.Seed(time.Now().UnixNano())

    var passwords []string
    for i := 0; i < config.Count; i++ {
        password := make([]byte, config.Length)
        for j := 0; j < config.Length; j++ {
            num, _ := randomUint8()
            password[j] = charset[num%uint8(len(charset))]
        }
        passwords = append(passwords, string(password))
    }
    return passwords
}

func randomUint8() (uint8, error) {
    b := make([]byte, 1)
    _, err := reader.Read(b)
    return b[0], err
}

func removeSimilar(s string) string {
    similars := "iIlL1oO0"
    for _, ch := range similars {
        s = removeChar(s, byte(ch))
    }
    return s
}

func removeChar(s string, c byte) string {
    var result []byte
    for i := 0; i < len(s); i++ {
        if s[i] != c {
            result = append(result, s[i])
        }
    }
    return string(result)
}

func SavePasswordsToFile(passwords []string, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, pwd := range passwords {
        _, err := file.WriteString(pwd + "\n")
        if err != nil {
            return err
        }
    }
    return nil
}
