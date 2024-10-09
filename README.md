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

1. Add a new alias:
   ```
   mantrid alias add myalias "echo Hello, World!"
   ```

2. List all aliases:
   ```
   mantrid alias list
   ```

3. Edit an existing alias:
   ```
   mantrid alias edit myalias "echo Hello, Universe!"
   ```

4. Remove an alias:
   ```
   mantrid alias remove myalias
   ```

5. Set up cloud synchronization:
   ```
   mantrid cloud setup
   ```

6. Sync your data to the cloud:
   ```
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

