package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"math/rand"
	"pwdforge/internal/generator"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		format, _ := cmd.Flags().GetString("format")
		customCharset, _ := cmd.Flags().GetString("custom-charset")
		usePassphrase, _ := cmd.Flags().GetBool("passphrase")
		enforceAll, _ := cmd.Flags().GetBool("enforce-all")
		inputFile, _ := cmd.Flags().GetString("input")
		configFile, _ := cmd.Flags().GetString("config")
		copyClip, _ := cmd.Flags().GetBool("clipboard")
		wordCount, _ := cmd.Flags().GetInt("word-count")

		// Load config file if provided
		var cfg *GenerateConfig
		if configFile != "" {
			var err error
			cfg, err = loadGenerateConfig(configFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config file: %v\n", err)
				os.Exit(1)
			}
		}
		// Only use config values if CLI flag is not set (flags take precedence)
		if cfg != nil {
			if !cmd.Flags().Changed("length") && cfg.Length > 0 {
				length = cfg.Length
			}
			if !cmd.Flags().Changed("count") && cfg.Count > 0 {
				count = cfg.Count
			}
			if !cmd.Flags().Changed("uppercase") {
				includeUpper = cfg.IncludeUpper
			}
			if !cmd.Flags().Changed("lowercase") {
				includeLower = cfg.IncludeLower
			}
			if !cmd.Flags().Changed("digits") {
				includeDigits = cfg.IncludeDigits
			}
			if !cmd.Flags().Changed("specials") {
				includeSpecials = cfg.IncludeSpecials
			}
			if !cmd.Flags().Changed("exclude-similar") {
				excludeSimilar = cfg.ExcludeSimilar
			}
			if !cmd.Flags().Changed("custom-charset") && cfg.CustomCharset != "" {
				customCharset = cfg.CustomCharset
			}
			if !cmd.Flags().Changed("enforce-all") {
				enforceAll = cfg.EnforceAll
			}
			if !cmd.Flags().Changed("passphrase") {
				usePassphrase = cfg.Passphrase
			}
			if !cmd.Flags().Changed("word-count") && cfg.WordCount > 0 {
				wordCount = cfg.WordCount
			}
		}

		var passwords []string
		if inputFile != "" {
			file, err := os.Open(inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				// Try JSON first
				var params GenerateConfig
				if err := json.Unmarshal([]byte(line), &params); err != nil {
					// Try YAML if JSON fails
					err2 := yaml.Unmarshal([]byte(line), &params)
					if err2 != nil {
						fmt.Fprintf(os.Stderr, "Skipping invalid input line: %s\n", line)
						continue
					}
				}
				// Merge params with CLI/config defaults
				merged := GenerateConfig{
					Length:          params.Length,
					Count:           params.Count,
					IncludeUpper:    params.IncludeUpper,
					IncludeLower:    params.IncludeLower,
					IncludeDigits:   params.IncludeDigits,
					IncludeSpecials: params.IncludeSpecials,
					ExcludeSimilar:  params.ExcludeSimilar,
					CustomCharset:   params.CustomCharset,
					EnforceAll:      params.EnforceAll,
					Passphrase:      params.Passphrase,
					WordCount:       params.WordCount,
				}
				if merged.Length == 0 {
					merged.Length = length
				}
				if merged.Count == 0 {
					merged.Count = 1
				}
				if !params.IncludeUpper && !params.IncludeLower && !params.IncludeDigits && !params.IncludeSpecials && params.CustomCharset == "" {
					merged.IncludeUpper = includeUpper
					merged.IncludeLower = includeLower
					merged.IncludeDigits = includeDigits
					merged.IncludeSpecials = includeSpecials
				}
				if merged.CustomCharset == "" {
					merged.CustomCharset = customCharset
				}
				if !merged.EnforceAll {
					merged.EnforceAll = enforceAll
				}
				if !merged.Passphrase {
					merged.Passphrase = usePassphrase
				}
				if merged.WordCount == 0 {
					merged.WordCount = wordCount
				}
				// Use merged config to generate password(s)
				if merged.Passphrase {
					wc := merged.WordCount
					if wc <= 0 {
						wc = 4
					}
					passwords = append(passwords, GeneratePassphrase(wc, nil))
				} else {
					cfg := generator.PasswordConfig{
						Length:          merged.Length,
						Count:           merged.Count,
						IncludeUpper:    merged.IncludeUpper,
						IncludeLower:    merged.IncludeLower,
						IncludeDigits:   merged.IncludeDigits,
						IncludeSpecials: merged.IncludeSpecials,
						ExcludeSimilar:  merged.ExcludeSimilar,
					}
					if merged.CustomCharset != "" {
						pw := make([]byte, merged.Length)
						for j := 0; j < merged.Length; j++ {
							idx := RandomInt(len(merged.CustomCharset))
							pw[j] = merged.CustomCharset[idx]
						}
						passwords = append(passwords, string(pw))
					} else {
						passwords = append(passwords, generator.GeneratePasswords(cfg)...)
					}
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
				os.Exit(1)
			}
		} else {
			if usePassphrase {
				wc := wordCount
				if wc <= 0 {
					wc = 4
				}
				passwords = []string{}
				for i := 0; i < count; i++ {
					passwords = append(passwords, GeneratePassphrase(wc, nil))
				}
			} else {
				cfg := generator.PasswordConfig{
					Length:          length,
					Count:           count,
					IncludeUpper:    includeUpper,
					IncludeLower:    includeLower,
					IncludeDigits:   includeDigits,
					IncludeSpecials: includeSpecials,
					ExcludeSimilar:  excludeSimilar,
				}
				if customCharset != "" {
					passwords = []string{}
					for i := 0; i < count; i++ {
						pw := make([]byte, length)
						for j := 0; j < length; j++ {
							idx := RandomInt(len(customCharset))
							pw[j] = customCharset[idx]
						}
						passwords = append(passwords, string(pw))
					}
				} else {
					passwords = generator.GeneratePasswords(cfg)
				}
				if enforceAll && customCharset == "" {
					for i, pw := range passwords {
						for {
							valid := true
							if includeUpper && !HasChar(pw, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
								valid = false
							}
							if includeLower && !HasChar(pw, "abcdefghijklmnopqrstuvwxyz") {
								valid = false
							}
							if includeDigits && !HasChar(pw, "0123456789") {
								valid = false
							}
							if includeSpecials && !HasChar(pw, "!@#$%^&*()-_=+[]{}|;:,.<>/?") {
								valid = false
							}
							if valid {
								break
							}
							pw = generator.GeneratePasswords(cfg)[0]
							passwords[i] = pw
						}
					}
				}
			}
		}

		// Output results in requested format
		if format == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(passwords)
		} else if format == "table" {
			fmt.Printf("%-30s %-12s %-8s\n", "Password", "Strength", "Entropy")
			for _, pwd := range passwords {
				strength, entropy, _ := generator.CheckPasswordStrength(pwd)
				fmt.Printf("%-30s %-12s %-8.2f\n", pwd, strength, entropy)
			}
		} else if format == "csv" {
			fmt.Println("Password,Strength,Entropy")
			for _, pwd := range passwords {
				strength, entropy, _ := generator.CheckPasswordStrength(pwd)
				fmt.Printf("%s,%s,%.2f\n", pwd, strength, entropy)
			}
		} else {
			for _, pwd := range passwords {
				if verbose {
					strength, entropy, suggestions := generator.CheckPasswordStrength(pwd)
					fmt.Printf("[+] Password: %s\t| Strength: %s | Entropy: %.2f\n", pwd, strength, entropy)
					if len(suggestions) > 0 {
						fmt.Println("  Suggestions:")
						for _, s := range suggestions {
							fmt.Printf("    - %s\n", s)
						}
					}
				} else {
					fmt.Println(pwd)
				}
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

		if copyClip && len(passwords) > 0 {
			fmt.Println("[Clipboard integration is currently unavailable due to Go import issues]")
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
	generateCmd.Flags().String("format", "plain", "Output format: plain, json, csv, table")
	generateCmd.Flags().String("custom-charset", "", "Custom character set for password generation")
	generateCmd.Flags().Bool("passphrase", false, "Generate passphrase using wordlist")
	generateCmd.Flags().Bool("enforce-all", false, "Enforce at least one of each selected character type")
	generateCmd.Flags().String("input", "", "Read password generation parameters from a file (JSON/YAML)")
	generateCmd.Flags().String("config", "", "Path to config file for default options")
	generateCmd.Flags().Bool("clipboard", false, "Copy first password to clipboard")
	generateCmd.Flags().Int("word-count", 4, "Number of words in passphrase (for --passphrase)")
	RootCmd.AddCommand(generateCmd)
}

// Helper for random int
func RandomInt(n int) int {
	return rand.Intn(n)
}

// Helper for HasChar
func HasChar(password, charSet string) bool {
	return strings.ContainsAny(password, charSet)
}

// Config struct for YAML/JSON config
// Only fields relevant to password generation

type GenerateConfig struct {
	Length          int    `yaml:"length" json:"length"`
	Count           int    `yaml:"count" json:"count"`
	IncludeUpper    bool   `yaml:"include_upper" json:"include_upper"`
	IncludeLower    bool   `yaml:"include_lower" json:"include_lower"`
	IncludeDigits   bool   `yaml:"include_digits" json:"include_digits"`
	IncludeSpecials bool   `yaml:"include_specials" json:"include_specials"`
	ExcludeSimilar  bool   `yaml:"exclude_similar" json:"exclude_similar"`
	CustomCharset   string `yaml:"custom_charset" json:"custom_charset"`
	EnforceAll      bool   `yaml:"enforce_all" json:"enforce_all"`
	Passphrase      bool   `yaml:"passphrase" json:"passphrase"`
	WordCount       int    `yaml:"word_count" json:"word_count"`
}

// Helper to load config from YAML or JSON
func loadGenerateConfig(path string) (*GenerateConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := yaml.NewDecoder(f)
	var cfg GenerateConfig
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Built-in wordlist for passphrase generation (short, for demo)
var defaultWordlist = []string{
	"apple", "banana", "cat", "dog", "elephant", "fish", "grape", "hat", "ice", "jungle",
	"kite", "lemon", "monkey", "nest", "orange", "pear", "queen", "rose", "sun", "tree",
	"umbrella", "violet", "wolf", "xray", "yak", "zebra",
}

// GeneratePassphrase creates a passphrase of n words from the wordlist
func GeneratePassphrase(wordCount int, wordlist []string) string {
	if wordCount <= 0 {
		wordCount = 4
	}
	if len(wordlist) == 0 {
		wordlist = defaultWordlist
	}
	words := make([]string, wordCount)
	for i := 0; i < wordCount; i++ {
		idx := RandomInt(len(wordlist))
		words[i] = wordlist[idx]
	}
	return strings.Join(words, "-")
}
