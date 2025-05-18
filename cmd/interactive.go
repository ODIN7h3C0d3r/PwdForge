package cmd

import (
	"bufio"
	"fmt"
	"os"
	"pwdforge/internal/generator"
	"pwdforge/internal/pwnchecker"
	"strings"

	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive mode for password generation and checking",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Welcome to PwdForge Interactive Mode!")
		for {
			fmt.Println("\nChoose an option:")
			fmt.Println("1. Generate Password(s)")
			fmt.Println("2. Check Password (pwned)")
			fmt.Println("3. Exit")
			fmt.Print("> ")
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			switch choice {
			case "1":
				interactiveGenerate(reader)
			case "2":
				interactiveCheck(reader)
			case "3":
				fmt.Println("Goodbye!")
				return
			default:
				fmt.Println("Invalid option.")
			}
		}
	},
}

func interactiveGenerate(reader *bufio.Reader) {
	fmt.Print("Password length (default 12): ")
	lengthStr, _ := reader.ReadString('\n')
	lengthStr = strings.TrimSpace(lengthStr)
	length := 12
	if lengthStr != "" {
		fmt.Sscanf(lengthStr, "%d", &length)
	}
	fmt.Print("How many passwords? (default 1): ")
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	count := 1
	if countStr != "" {
		fmt.Sscanf(countStr, "%d", &count)
	}
	fmt.Print("Include uppercase? (Y/n): ")
	upperStr, _ := reader.ReadString('\n')
	upper := !strings.HasPrefix(strings.ToLower(strings.TrimSpace(upperStr)), "n")
	fmt.Print("Include lowercase? (Y/n): ")
	lowerStr, _ := reader.ReadString('\n')
	lower := !strings.HasPrefix(strings.ToLower(strings.TrimSpace(lowerStr)), "n")
	fmt.Print("Include digits? (Y/n): ")
	digitStr, _ := reader.ReadString('\n')
	digit := !strings.HasPrefix(strings.ToLower(strings.TrimSpace(digitStr)), "n")
	fmt.Print("Include specials? (Y/n): ")
	specialStr, _ := reader.ReadString('\n')
	special := !strings.HasPrefix(strings.ToLower(strings.TrimSpace(specialStr)), "n")
	fmt.Print("Exclude similar/confusing characters? (y/N): ")
	similarStr, _ := reader.ReadString('\n')
	excludeSimilar := strings.HasPrefix(strings.ToLower(strings.TrimSpace(similarStr)), "y")
	fmt.Print("Copy to clipboard? (y/N): ")
	clipStr, _ := reader.ReadString('\n')
	copyClip := strings.HasPrefix(strings.ToLower(strings.TrimSpace(clipStr)), "y")

	passwords := generator.GeneratePasswords(generator.PasswordConfig{
		Length:          length,
		Count:           count,
		IncludeUpper:    upper,
		IncludeLower:    lower,
		IncludeDigits:   digit,
		IncludeSpecials: special,
		ExcludeSimilar:  excludeSimilar,
	})
	for _, pwd := range passwords {
		strength, entropy, suggestions := generator.CheckPasswordStrength(pwd)
		fmt.Printf("[+] Password: %s\t| Strength: %s | Entropy: %.2f\n", pwd, strength, entropy)
		if len(suggestions) > 0 {
			fmt.Println("  Suggestions:")
			for _, s := range suggestions {
				fmt.Printf("    - %s\n", s)
			}
		}
		if copyClip {
			fmt.Println("[Clipboard integration is currently unavailable due to Go import issues]")
		}
	}
	// If copyClip is true, copy the first password to clipboard and notify
	if copyClip && len(passwords) > 0 {
		fmt.Println("[Clipboard integration is currently unavailable due to Go import issues]")
	}
}

func interactiveCheck(reader *bufio.Reader) {
	fmt.Print("Enter password to check: ")
	pwd, _ := reader.ReadString('\n')
	pwd = strings.TrimSpace(pwd)
	if pwd == "" {
		fmt.Println("Password cannot be empty.")
		return
	}
	exposed, count, err := pwnchecker.CheckPasswordPwned(pwd)
	if err != nil {
		fmt.Printf("Error checking password: %v\n", err)
		return
	}
	if exposed {
		fmt.Printf("[!] WARNING: This password has been found in breaches (%d times).\n", count)
	} else {
		fmt.Println("[+] Good news! This password hasn't been found in any known breaches.")
	}
}

func init() {
	RootCmd.AddCommand(interactiveCmd)
}
