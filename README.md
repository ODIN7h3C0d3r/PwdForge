# PwdForge

**PwdForge** is an advanced, versatile, and interactive password manager CLI written in Go. It supports secure password and passphrase generation, batch operations, output in multiple formats, config file support, clipboard integration (stub), and interactive CLI mode. It also integrates with HaveIBeenPwned for breach checks.

---

## Features

- **Password Generation**: Generate strong passwords with customizable length, character sets, and options.
- **Passphrase Generation**: Generate memorable passphrases using a built-in wordlist.
- **Batch Operations**: Generate passwords in bulk using batch input files (JSON/YAML per line).
- **Output Formats**: Output passwords in plain text, JSON, CSV, or table format.
- **Config File Support**: Use YAML or JSON config files to set default options.
- **Clipboard Integration**: (Stub) Option to copy the first password to clipboard (integration can be restored).
- **Interactive CLI**: Interactive mode for password generation and breach checking.
- **Breach Check**: Check passwords against known breaches using the HaveIBeenPwned API.
- **Enforce-All**: Ensure at least one of each selected character type in generated passwords.
- **Custom Charset**: Use a custom character set for password generation.

---

## Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/yourusername/PwdForge.git
   cd PwdForge
   ```

2. **Build the CLI:**

   ```sh
   go build ./...
   ```

---

## Usage

### Password Generation

Generate a single password (default 12 chars, all types):

```sh
 go run main.go generate
```

Generate 3 passwords, 16 chars, all types, table output:

```sh
 go run main.go generate --length 16 --count 3 --format table
```

Generate a passphrase (4 words):

```sh
 go run main.go generate --passphrase
```

Generate a passphrase with 6 words:

```sh
 go run main.go generate --passphrase --word-count 6
```

Generate using a custom charset:

```sh
 go run main.go generate --length 20 --custom-charset "abc123!@#"
```

### Batch Password Generation

Prepare a batch file (e.g., `test_batch.txt`):

```json
{"length": 10, "count": 1, "uppercase": true, "lowercase": true, "digits": true, "specials": false}
{"length": 14, "count": 2, "uppercase": true, "lowercase": true, "digits": false, "specials": true}
```

Run batch generation:

```sh
 go run main.go generate --input test_batch.txt --format table
```

### Config File Support

Create a config file (YAML or JSON):

```yaml
length: 16
count: 2
include_upper: true
include_lower: true
include_digits: true
include_specials: true
```

Use it:

```sh
 go run main.go generate --config config.yaml
```

### Breach Check

Check a password against known breaches:

```sh
 go run main.go checkpwn --password "yourpassword"
```

Batch breach check:

```sh
 go run main.go checkpwn --input test_pwn.txt --format table
```

### Interactive CLI

Start interactive mode:

```sh
 go run main.go interactive
```

---

## Flags (Password Generation)

- `--length, -l`         Password length (default: 12)
- `--count, -c`          Number of passwords to generate (default: 1)
- `--uppercase, -u`      Include uppercase letters (default: true)
- `--lowercase, -w`      Include lowercase letters (default: true)
- `--digits, -d`         Include digits (default: true)
- `--specials, -s`       Include special characters (default: true)
- `--exclude-similar`    Exclude similar/confusing characters
- `--output, -o`         Save passwords to a file
- `--verbose, -v`        Show detailed output (strength, etc.)
- `--format`             Output format: plain, json, csv, table
- `--custom-charset`     Custom character set for password generation
- `--passphrase`         Generate passphrase using wordlist
- `--word-count`         Number of words in passphrase (default: 4)
- `--enforce-all`        Enforce at least one of each selected character type
- `--input`              Batch input file (JSON/YAML per line)
- `--config`             Config file (YAML/JSON)
- `--clipboard`          Copy first password to clipboard (stub)

---

## Development & Contribution

- PRs and issues welcome!
- To restore clipboard integration, add a Go clipboard library and update `pkg/clipboard.go` and usage sites.

---

## License

MIT

---

## Author

- [Your Name or GitHub handle]

---

## Disclaimer

This tool is for educational and personal use. Do not use for illegal or malicious purposes. Always use strong, unique passwords for every site and service.
