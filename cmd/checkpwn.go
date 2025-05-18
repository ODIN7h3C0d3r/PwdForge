package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"pwdforge/internal/pwnchecker"

	"github.com/spf13/cobra"
)

var checkpwnCmd = &cobra.Command{
	Use:   "checkpwn",
	Short: "Check if a password has been exposed in data breaches",
	Long:  "Uses the HaveIBeenPwned API to securely check if a password has been pwned.",
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		inputFile, _ := cmd.Flags().GetString("input")
		format, _ := cmd.Flags().GetString("format")
		if strings.TrimSpace(password) == "" && strings.TrimSpace(inputFile) == "" {
			fmt.Fprintln(os.Stderr, "Error: Either --password or --input must be provided.")
			os.Exit(1)
		}
		if strings.TrimSpace(password) != "" && strings.TrimSpace(inputFile) != "" {
			fmt.Fprintln(os.Stderr, "Error: Only one of --password or --input should be provided.")
			os.Exit(1)
		}
		var results []map[string]interface{}
		if inputFile != "" {
			file, err := os.Open(inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				pw := strings.TrimSpace(scanner.Text())
				if pw == "" {
					continue
				}
				exposed, count, err := pwnchecker.CheckPasswordPwned(pw)
				result := map[string]interface{}{
					"password": pw,
					"exposed":  exposed,
					"count":    count,
					"error":    "",
				}
				if err != nil {
					result["error"] = err.Error()
				}
				results = append(results, result)
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
				os.Exit(1)
			}
		} else {
			exposed, count, err := pwnchecker.CheckPasswordPwned(password)
			result := map[string]interface{}{
				"password": password,
				"exposed":  exposed,
				"count":    count,
				"error":    "",
			}
			if err != nil {
				result["error"] = err.Error()
			}
			results = append(results, result)
		}
		// Output results in requested format
		if format == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(results)
		} else if format == "table" {
			fmt.Printf("%-20s %-8s %-8s %-s\n", "Password", "Exposed", "Count", "Error")
			for _, r := range results {
				fmt.Printf("%-20s %-8v %-8v %-s\n", r["password"], r["exposed"], r["count"], r["error"])
			}
		} else {
			for _, r := range results {
				if r["error"] != "" {
					fmt.Fprintf(os.Stderr, "Error checking password '%s': %s\n", r["password"], r["error"])
					continue
				}
				if r["exposed"].(bool) {
					fmt.Fprintf(os.Stdout, "[!] WARNING: '%s' found in breaches (%d times).\n", r["password"], r["count"])
				} else {
					fmt.Fprintf(os.Stdout, "[+] '%s' not found in any known breaches.\n", r["password"])
				}
			}
		}
	},
}

func init() {
	checkpwnCmd.Flags().StringP("password", "p", "", "Password to check against breaches")
	// Removed MarkFlagRequired("password")
	checkpwnCmd.Flags().String("format", "plain", "Output format: plain, json, table")
	checkpwnCmd.Flags().String("input", "", "Read passwords to check from a file (one per line)")
	checkpwnCmd.Flags().String("config", "", "Path to config file for default options")
	RootCmd.AddCommand(checkpwnCmd)
}
