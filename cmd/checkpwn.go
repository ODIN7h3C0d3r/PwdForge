package cmd

import (
    "fmt"
    "strings"

    "github.com/spf13/cobra"
    "pwdforge/internal/pwnchecker"
)

var checkpwnCmd = &cobra.Command{
    Use:   "checkpwn",
    Short: "Check if a password has been exposed in data breaches",
    Long:  "Uses the HaveIBeenPwned API to securely check if a password has been pwned.",
    Run: func(cmd *cobra.Command, args []string) {
        password, _ := cmd.Flags().GetString("password")

        if strings.TrimSpace(password) == "" {
            fmt.Fprintln(os.Stderr, "Error: Password cannot be empty.")
            os.Exit(1)
        }

        exposed, count, err := pwnchecker.CheckPasswordPwned(password)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error checking password: %v\n", err)
            os.Exit(1)
        }

        if exposed {
            fmt.Fprintf(os.Stdout, "[!] WARNING: This password has been found in breaches (%d times).\n", count)
        } else {
            fmt.Println("[+] Good news! This password hasn't been found in any known breaches.")
        }
    },
}

func init() {
    checkpwnCmd.Flags().StringP("password", "p", "", "Password to check against breaches")
    _ = checkpwnCmd.MarkFlagRequired("password")

    cmd.RootCmd.AddCommand(checkpwnCmd)
}
