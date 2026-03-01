# Javi.Dots

> 🍴 **Fork of [Gentleman.Dots](https://github.com/Gentleman-Programming/Gentleman.Dots)** by [Gentleman Programming](https://youtube.com/@GentlemanProgramming). All upstream updates are regularly merged.

```
                  ██╗███╗   ██╗███████╗
                  ██║████╗  ██║╚══███╔╝
                  ██║██╔██╗ ██║  ███╔╝
             ██   ██║██║╚██╗██║ ███╔╝
             ╚█████╔╝██║ ╚████║███████╗
              ╚════╝ ╚═╝  ╚═══╝╚══════╝

    ██╗ █████╗ ██╗   ██╗██╗    ██████╗  ██████╗ ████████╗███████╗
    ██║██╔══██╗██║   ██║██║    ██╔══██╗██╔═══██╗╚══██╔══╝██╔════╝
    ██║███████║██║   ██║██║    ██║  ██║██║   ██║   ██║   ███████╗
██  ██║██╔══██║╚██╗ ██╔╝██║    ██║  ██║██║   ██║   ██║   ╚════██║
╚████╔╝██║  ██║ ╚████╔╝ ██║ ██╗██████╔╝╚██████╔╝   ██║   ███████║
 ╚═══╝ ╚═╝  ╚═╝  ╚═══╝  ╚═╝ ╚═╝╚═════╝  ╚═════╝    ╚═╝   ╚══════╝
```

Personal dotfiles & TUI installer for terminal, shell, Neovim, and AI tools.

---

## What is this?

A complete development environment configuration including:

- **Neovim** with LSP, autocompletion, and AI assistants
- **AI Tools**: Claude Code, OpenCode, Gemini CLI, GitHub Copilot, Codex CLI, Qwen Code
- **AI Framework**: 198 modules (agents, skills, hooks, commands, SDD, MCP) + 6 domain orchestrators
- **Shells**: Fish, Zsh, Nushell
- **Terminal Multiplexers**: Tmux, Zellij
- **Terminal Emulators**: Alacritty, WezTerm, Kitty, Ghostty

## Quick Start

```bash
git clone https://github.com/JNZader/Javi.Dots.git
cd Javi.Dots/installer
go build -o gentleman-installer ./cmd/gentleman-installer
./gentleman-installer
```

## Documentation

Use the **sidebar** to navigate all docs, or jump directly:

| Section | Description |
|---------|-------------|
| [TUI Installer](tui-installer.md) | Interactive installer features, navigation, backup/restore |
| [AI Tools & Framework](ai-tools-integration.md) | AI tools selection, framework presets, category drill-down |
| [AI Framework Modules](ai-framework-modules.md) | Complete reference of all 198 modules + 28 curated skills |
| [Agent Teams Lite](agent-teams-lite.md) | Lightweight SDD framework with 9 sub-agents |
| [AI Configuration](ai-configuration.md) | Claude Code, OpenCode, Copilot setup |
| [Manual Installation](manual-installation.md) | Step-by-step setup for all platforms |
| [Neovim Keymaps](neovim-keymaps.md) | Complete keybinding reference |
| [Vim Trainer](vim-trainer-spec.md) | RPG-style Vim training system |
| [Contributing](contributing.md) | Development setup and release process |
| [Docker Testing](docker-testing.md) | E2E testing with containers |

## Support

- **Issues**: [GitHub Issues](https://github.com/JNZader/Javi.Dots/issues)
- **Upstream**: [Gentleman.Dots](https://github.com/Gentleman-Programming/Gentleman.Dots)
- **Discord**: [Gentleman Programming](https://discord.gg/gentleman-programming)
