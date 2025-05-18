package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
    Use:   "pwdforge",
    Short: "PwdForge - A secure password generator written in Go",
    Long:  `PwdForge generates strong passwords, checks their strength, and verifies if they've been pwned.`,
}

func Execute() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
