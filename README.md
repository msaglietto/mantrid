# Mantrid: Your Command-Line Productivity Companion

Mantrid is a user-friendly command-line tool that streamlines your workflow by storing and running command aliases. Built with Go, it's a fast, lightweight solution for developers, system administrators, and power users who want to optimize their command-line experience.

## 🚀 Features

- **Intuitive Alias Management**: Create, list, edit, and remove aliases effortlessly.
- **Parameter Substitution**: Use `$1`, `$@`, and friends, or auto-append arguments.
- **Cross-Platform Compatibility**: Works on Linux, macOS, and Windows.
- **Lightweight and Fast**: Written in Go for optimal performance.
- **Easy to Use**: Simple, intuitive commands for all operations.

## 🛠️ Installation

### Homebrew (macOS / Linux)

```bash
brew install msaglietto/tap/mantrid
```

### Binary Download

Download the latest binary for your platform from the [releases page](https://github.com/msaglietto/mantrid/releases/latest):

| Platform | Architecture | File |
|----------|-------------|------|
| macOS | Apple Silicon (arm64) | `mantrid_X.Y.Z_darwin_arm64.tar.gz` |
| macOS | Intel (x86_64) | `mantrid_X.Y.Z_darwin_x86_64.tar.gz` |
| Linux | x86_64 | `mantrid_X.Y.Z_linux_x86_64.tar.gz` |
| Linux | ARM64 | `mantrid_X.Y.Z_linux_aarch64.tar.gz` |
| Linux | ARMv6 | `mantrid_X.Y.Z_linux_armv6.tar.gz` |
| Windows | x86_64 | `mantrid_X.Y.Z_windows_x86_64.zip` |
| Windows | ARM64 | `mantrid_X.Y.Z_windows_arm64.zip` |

Extract the archive and move the binary to a directory in your `$PATH`.

### Linux Packages

**Debian / Ubuntu:**

Download the `.deb` file from the [releases page](https://github.com/msaglietto/mantrid/releases/latest) and install:

```bash
sudo dpkg -i mantrid_X.Y.Z_linux_amd64.deb
```

**Fedora / RHEL:**

Download the `.rpm` file from the [releases page](https://github.com/msaglietto/mantrid/releases/latest) and install:

```bash
sudo rpm -i mantrid_X.Y.Z_linux_x86_64.rpm
```

### Windows

**Scoop:**

```powershell
scoop bucket add msaglietto https://github.com/msaglietto/scoop-bucket
scoop install mantrid
```

### Go Install

> **Note:** Requires [Go 1.23+](https://go.dev/dl/)

```bash
go install github.com/msaglietto/mantrid@latest
```

Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `$PATH`.

### Verify Installation

```bash
mantrid --version
```

## 🏁 Quick Start

### Basic Alias Management

1. Add a new alias:
   ```bash
   mantrid alias add hello "echo Hello, World!"
   ```

2. Execute an alias:
   ```bash
   mantrid do hello
   ```

3. List all aliases:
   ```bash
   mantrid alias list
   ```

4. Edit an existing alias:
   ```bash
   mantrid alias edit hello "echo Hello, Universe!"
   ```

5. Remove an alias:
   ```bash
   mantrid alias remove hello
   ```

### Simple Aliases (Auto-Append)

For simple command aliases without placeholders, parameters are automatically appended:

```bash
# Create simple aliases
mantrid alias add ls "ls"
mantrid alias add dk "docker"
mantrid alias add k "kubectl"

# Parameters are automatically appended
mantrid do ls -- -la /tmp           # Executes: ls -la /tmp
mantrid do dk -- ps -a              # Executes: docker ps -a
mantrid do k -- get pods            # Executes: kubectl get pods
```

### Aliases with Parameter Substitution

For advanced control, use placeholders:

- **Positional parameters**: `$1`, `$2`, `$3`, etc.
- **All parameters**: `$@` or `$*`

**Examples:**

```bash
# Create alias with placeholders
mantrid alias add greet "echo Hello, $1!"
mantrid alias add deploy "kubectl apply -f $1 -n $2"

# Parameters are substituted
mantrid do greet World              # Executes: echo Hello, World!
mantrid do deploy app.yaml prod     # Executes: kubectl apply -f app.yaml -n prod

# Use all parameters with $@
mantrid alias add search "grep -r $@ ."
mantrid do search "TODO"            # Executes: grep -r TODO .
```

### Passing Flags to Aliases

When you need to pass flags (arguments starting with `-` or `--`) to your aliases, use the `--` separator to prevent Cobra from interpreting them as flags to the `do` command itself:

```bash
# Simple execution
mantrid do hello

# With parameters
mantrid do greet Alice

# Using -- separator (useful for flags)
mantrid do ls -- -la /tmp
mantrid do docker -- run --rm -it ubuntu bash
mantrid do grep -- -r "pattern" .

# The -- tells Mantrid to pass everything after it as parameters
```

The `--` separator is especially useful when your alias needs to receive flags that would otherwise conflict with Mantrid's own command-line parsing.

**Security Note:** Aliases execute commands directly in your system shell. Only create aliases for commands you trust. Parameter substitution does not perform shell escaping - use with caution.

## ⚙️ Configuration

Mantrid reads an optional config file from `~/.mantrid/config.yaml`. All values
have sensible defaults, so a config file is not required.

```yaml
# Storage configuration
alias_file: "~/.mantrid/aliases.json"
storage_type: "json"

# Logging configuration
# Levels: debug, info, warn, error. Default is "warn" (quiet during normal use).
log_level: "warn"
log_format: "json"
```

Settings can also be overridden with `MANTRID_`-prefixed environment variables
(e.g. `MANTRID_LOG_LEVEL=debug`).

## 🌟 Why Mantrid?

- **Boost Productivity**: Save time by creating shortcuts for your most-used commands.
- **Customizable**: Tailor your command-line environment to your specific needs.
- **Cross-Platform**: One tool that behaves the same on Linux, macOS, and Windows.

## 🗺️ Roadmap

Planned for future releases (not yet available):

- **Dotfile management**: Track and sync configuration files.
- **Cloud synchronization**: Back up and sync aliases across machines.

## 📜 License

Mantrid is released under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.

## 🙏 Acknowledgements

Mantrid is built with the following open-source projects:
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)

---

Mantrid: Simplify your command-line life, one alias at a time. 🚀

