# Mantrid: Your Command-Line Productivity Companion

Mantrid is a powerful, user-friendly command-line tool designed to streamline your workflow by managing aliases and dotfiles with ease. Built with Go, Mantrid offers a robust solution for developers, system administrators, and power users who want to optimize their command-line experience across multiple devices.

## 🚀 Features

- **Intuitive Alias Management**: Create, list, edit, and remove aliases effortlessly.
- **Dotfile Syncing**: Keep your configuration files in sync across multiple machines.
- **Cloud Synchronization**: Store and sync your aliases and dotfiles securely in the cloud.
- **Cross-Platform Compatibility**: Works seamlessly on Linux, macOS, and Windows.
- **Lightweight and Fast**: Written in Go for optimal performance.
- **Easy to Use**: Simple, intuitive commands for all operations.

## 🛠️ Installation

```bash
go install github.com/msaglietto/mantrid@latest
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

### Cloud Synchronization

5. Set up cloud synchronization:
   ```bash
   mantrid cloud setup
   ```

6. Sync your data to the cloud:
   ```bash
   mantrid cloud sync
   ```

## 🌟 Why Mantrid?

- **Boost Productivity**: Save time by creating shortcuts for your most-used commands.
- **Consistency Across Machines**: Sync your dotfiles and aliases across multiple computers with ease.
- **Cloud-Powered**: Keep your configurations backed up and accessible from anywhere.
- **Customizable**: Tailor your command-line environment to your specific needs.
- **Version Control**: Keep track of changes to your aliases and dotfiles over time.
- **Secure**: Your data is encrypted and securely stored in the cloud.

## ☁️ Cloud Synchronization

Mantrid offers seamless cloud synchronization to keep your aliases and dotfiles consistent across all your devices:

- **Automatic Backups**: Your configurations are always safe and up-to-date.
- **Multi-Device Sync**: Access your aliases and dotfiles from any machine.
- **Selective Sync**: Choose which items to sync and which to keep local.
- **Conflict Resolution**: Smart handling of changes made on different devices.
- **Provider Flexibility**: Choose your preferred cloud storage provider.

To get started with cloud sync:

1. Run `mantrid cloud setup` and follow the prompts to connect to your cloud account.
2. Use `mantrid cloud sync` to synchronize your data.
3. On a new machine, run `mantrid cloud restore` to retrieve your configurations.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## 📜 License

Mantrid is released under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.

## 🙏 Acknowledgements

Mantrid is built with love and the following amazing open-source projects:
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)

---

Mantrid: Simplify your command-line life, one alias at a time. Now with the power of the cloud! 🚀☁️✨

