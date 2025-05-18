package cmd

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/spf13/cobra"
    "pwdforge/internal/generator"
)

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate one or more secure passwords",
    Run: func(cmd *cobra.Command, args []string) {
        length, _ := cmd.Flags().GetInt("length")
        count, _ := cmd.Flags().GetInt("count")
        includeUpper, _ := cmd.Flags().GetBool("uppercase")
        includeLower, _ := cmd.Flags().GetBool("lowercase")
        includeDigits, _ := cmd.Flags().GetBool("digits")
        includeSpecials, _ := cmd.Flags().GetBool("specials")
        excludeSimilar, _ := cmd.Flags().GetBool("exclude-similar")
        outputFile, _ := cmd.Flags().GetString("output")
        verbose, _ := cmd.Flags().GetBool("verbose")

        passwords := generator.GeneratePasswords(generator.PasswordConfig{
            Length:          length,
            Count:           count,
            IncludeUpper:    includeUpper,
            IncludeLower:    includeLower,
            IncludeDigits:   includeDigits,
            IncludeSpecials: includeSpecials,
            ExcludeSimilar:  excludeSimilar,
        })

        for _, pwd := range passwords {
            if verbose {
                strength := generator.CheckPasswordStrength(pwd)
                fmt.Printf("[+] Password: %s\t| Strength: %s\n", pwd, strength)
            } else {
                fmt.Println(pwd)
            }
        }

        if outputFile != "" {
            err := generator.SavePasswordsToFile(passwords, outputFile)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Error saving passwords to file: %v\n", err)
                os.Exit(1)
            }
            fmt.Fprintf(os.Stdout, "[+] Saved %d passwords to %s\n", len(passwords), outputFile)
        }
    },
}

func init() {
    generateCmd.Flags().IntP("length", "l", 12, "Length of the password")
    generateCmd.Flags().IntP("count", "c", 1, "Number of passwords to generate")
    generateCmd.Flags().BoolP("uppercase", "u", true, "Include uppercase letters")
    generateCmd.Flags().BoolP("lowercase", "w", true, "Include lowercase letters")
    generateCmd.Flags().BoolP("digits", "d", true, "Include digits")
    generateCmd.Flags().BoolP("specials", "s", true, "Include special characters")
    generateCmd.Flags().Bool("exclude-similar", false, "Exclude similar/confusing characters (e.g., l, 1, O, 0)")
    generateCmd.Flags().StringP("output", "o", "", "Save passwords to a file")
    generateCmd.Flags().BoolP("verbose", "v", false, "Show detailed output (strength, etc.)")

    rootCmd.AddCommand(generateCmd)
}
