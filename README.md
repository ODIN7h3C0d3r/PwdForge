# PwdForge: Advanced Password Manager CLI

---

**PwdForge** is a modern, production-grade command-line password manager and generator written in Go. It empowers users, power users, and developers to generate, manage, and audit passwords and passphrases with maximum flexibility, security, and automation. PwdForge supports batch operations, config files, output in multiple formats, interactive and scripted use, and integrates with HaveIBeenPwned for breach checks.

---

## 📋 Menu

- [Features](#features)
- [Quick Start](#quick-start)
- [Normal User Guide](#normal-user-guide)
- [Power User Guide](#power-user-guide)
- [Developer Guide](#developer-guide)
- [Configuration](#configuration)
- [Batch & Automation](#batch--automation)
- [Interactive Mode](#interactive-mode)
- [Breach Checking](#breach-checking)
- [Clipboard Integration](#clipboard-integration)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

---

## 🚀 Features

- **Secure password generation** with customizable length, charset, and rules
- **Passphrase generation** (memorable, multi-word)
- **Batch operations** via input files (JSON/YAML per line)
- **Output formats**: plain, JSON, CSV, table
- **Config file support** (YAML/JSON)
- **Clipboard integration** (stub, see below)
- **Interactive CLI** for guided password and breach check workflows
- **Breach check** via HaveIBeenPwned API
- **Enforce-all**: require at least one of each selected character type
- **Custom charset** for advanced password policies
- **Verbose output** with strength/entropy analysis

---

## ⚡ Quick Start

Clone and build:

```sh

git clone https://github.com/ODIN7h3C0d3r/PwdForge.git
cd PwdForge
go build ./...

```

Generate a strong password:

```sh

go run main.go generate --length 16 --count 2 --format table

```

Check a password for breaches:

```sh

go run main.go checkpwn --password "MySecret123!"

```

---

## 👤 Normal User Guide

**Generate a password:**

```sh

go run main.go generate

```
(Default: 12 chars, all types)

**Generate a passphrase:**

```sh

go run main.go generate --passphrase

```

**Check if a password is breached:**

```sh

go run main.go checkpwn --password "yourpassword"

```

**Interactive mode:**

```sh

go run main.go interactive

```
(Menu-driven, no flags needed)

---

## ⚙️ Power User Guide

**Batch password generation:**

Create `test_batch.txt`:

```json
{"length": 10, "count": 1, "uppercase": true, "lowercase": true, "digits": true, "specials": false}
{"length": 14, "count": 2, "uppercase": true, "lowercase": true, "digits": false, "specials": true}
```

Run:

```sh

go run main.go generate --input test_batch.txt --format table

```

**Save passwords to a file:**

```sh

go run main.go generate --count 5 --output mypasswords.txt

```

**Use a config file (YAML/JSON):**

`config.yaml`:

```yaml
length: 20
count: 3
include_upper: true
include_lower: true
include_digits: true
include_specials: true
enforce_all: true
```

Run:

```sh

go run main.go generate --config config.yaml

```

**Custom charset:**

```sh

go run main.go generate --length 24 --custom-charset "abc123!@#"

```

**Enforce all character types:**

```sh

go run main.go generate --enforce-all

```

**Output as JSON/CSV:**

```sh

go run main.go generate --format json
go run main.go generate --format csv

```

---

## 🛠️ Developer Guide

- **Extend wordlist:** Edit `defaultWordlist` in `cmd/generate.go`.
- **Add new output formats:** Edit output section in `cmd/generate.go`.
- **Integrate clipboard:** Replace stub in `pkg/clipboard.go` and usage sites with a Go clipboard library (e.g., `github.com/atotto/clipboard`).
- **Add new breach sources:** Extend `internal/pwnchecker/pwncheck.go`.
- **Testing:**

```sh

go test ./...

```

- **Build binary:**

```sh

go build -o pwdforge main.go

```

---

## ⚙️ Configuration

**Config file (YAML/JSON):**

All CLI flags can be set in a config file. CLI flags override config.

Example YAML:

```yaml
length: 16
count: 2
include_upper: true
include_lower: true
include_digits: true
include_specials: true
passphrase: false
word_count: 4
enforce_all: false
custom_charset: ""
```

---

## 📦 Batch & Automation

- **Batch input:** Each line in the input file is a JSON or YAML object specifying password parameters.
- **Output file:** Use `--output` to save results.
- **Script integration:** Output in JSON/CSV for easy parsing.

---

## 🖥️ Interactive Mode

Run `go run main.go interactive` for a menu-driven experience:

1. Generate Password(s)
2. Check Password (pwned)
3. Exit

Prompts for all options, no flags needed.

---

## 🔎 Breach Checking

- Uses HaveIBeenPwned API (k-anonymity, privacy-safe)
- Single or batch mode supported
- Output in plain, table, or JSON

---

## 📋 Clipboard Integration

- Currently a stub (prints a warning)
- To enable, add a Go clipboard library and update `pkg/clipboard.go` and usage sites

---

## ❓ FAQ

- **Q: Is my password sent to the internet?**
  - A: Only for breach checks, and only a hash prefix (k-anonymity, never the full password).
- **Q: Can I use this in scripts?**
  - A: Yes! Use `--format json` or `--format csv` for easy parsing.
- **Q: How do I add new wordlists or charsets?**
  - A: Edit the Go source or pass via CLI/config.

---

## 🤝 Contributing

- PRs and issues welcome!
- Please lint and test before submitting.

---

## 📄 License

MIT

---

## 👤 Author

- [ODIN7h3C0d3r](https://github.com/ODIN7h3C0d3r)

---

## ⚠️ Disclaimer

This tool is for educational and personal use. Do not use for illegal or malicious purposes. Always use strong, unique passwords for every site and service.
