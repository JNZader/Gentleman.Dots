# Javi.Dots

> 🍴 **Fork of [Gentleman.Dots](https://github.com/Gentleman-Programming/Gentleman.Dots)** by [Gentleman Programming](https://youtube.com/@GentlemanProgramming). All upstream updates are regularly merged.

> ℹ️ **Update (January 2026)**: OpenCode now supports Claude Max/Pro subscriptions via the `opencode-anthropic-auth` plugin (included in this config). Both **Claude Code** and **OpenCode** work with your Claude subscription. *Note: This workaround is stable for now, but Anthropic could block it in the future.*
📄 Read this in: **English** | [Español](README.es.md)

## Table of Contents

- [What is this?](#what-is-this)
- [Quick Start](#quick-start)
- [Supported Platforms](#supported-platforms)
- [AI Tools & Framework](#-ai-tools--framework)
- [Vim Mastery Trainer](#-vim-mastery-trainer)
- [Documentation](#documentation)
- [Tools Overview](#tools-overview)
- [Project Structure](#project-structure)
- [Support](#support)

---

## Preview

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

---

## What is this?

A complete development environment configuration including:

- **Neovim** with LSP, autocompletion, and AI assistants (Claude Code, Gemini, OpenCode)
- **Zed** editor with Vim mode and AI agent support
- **AI Tools**: Claude Code, OpenCode, Gemini CLI, GitHub Copilot, Codex CLI with configs, skills, and themes
- **AI Framework**: 198 modules (72 agents, 85 skills, 10 hooks, 20 commands, 9 MCP servers) + 28 curated skills, with preset or custom selection
- **Shells**: Fish, Zsh, Nushell
- **Terminal Multiplexers**: Tmux, Zellij
- **Terminal Emulators**: Alacritty, WezTerm, Kitty, Ghostty

## Quick Start

### Option 1: Homebrew (Recommended)

```bash
brew install Gentleman-Programming/tap/gentleman-dots
gentleman-dots
```

### Option 2: Direct Download

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-darwin-arm64 -o gentleman.dots

# macOS Intel
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-darwin-amd64 -o gentleman.dots

# Linux x86_64
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-linux-amd64 -o gentleman.dots

# Linux ARM64 (Raspberry Pi, etc.)
curl -fsSL https://github.com/JNZader/Javi.Dots/releases/latest/download/gentleman-installer-linux-arm64 -o gentleman.dots

# Then run
chmod +x gentleman.dots
./gentleman.dots
```

### Option 3: Termux (Android)

Termux requires building the installer locally (Go cross-compilation to Android has limitations).

```bash
# 1. Install dependencies
pkg update && pkg upgrade
pkg install git golang

# 2. Clone the repository
git clone https://github.com/JNZader/Javi.Dots.git
cd Javi.Dots/installer

# 3. Build and run
go build -o ~/gentleman-installer ./cmd/gentleman-installer
cd ~
./gentleman-installer
```

| Termux Support | Status |
|----------------|--------|
| Shells (Fish, Zsh, Nushell) | ✅ Available |
| Multiplexers (Tmux, Zellij) | ✅ Available |
| Neovim with full config | ✅ Available |
| Nerd Fonts | ✅ Auto-installed to `~/.termux/font.ttf` |
| Terminal emulators | ❌ Not applicable |
| Homebrew | ❌ Uses `pkg` instead |

> **Tip:** After installation, restart Termux to apply the font, then run `tmux` or `zellij` to start your configured environment.

The TUI guides you through selecting your preferred tools and handles all the configuration automatically.

> **Windows users:** You must set up WSL first. See the [Manual Installation Guide](docs/manual-installation.md#windows-wsl).

---

## Supported Platforms

| Platform | Architecture | Install Method | Package Manager |
|----------|--------------|----------------|-----------------|
| macOS | Apple Silicon (ARM64) | Homebrew, Direct Download | Homebrew |
| macOS | Intel (x86_64) | Homebrew, Direct Download | Homebrew |
| Linux (Ubuntu/Debian) | x86_64, ARM64 | Homebrew, Direct Download | Homebrew |
| Linux (Fedora/RHEL) | x86_64, ARM64 | Direct Download | dnf |
| Linux (Arch) | x86_64 | Homebrew, Direct Download | Homebrew |
| Windows | WSL | Direct Download (see docs) | Homebrew |
| Android | Termux (ARM64) | Build locally (see above) | pkg |

---

## 🤖 AI Tools & Framework

The installer includes a complete AI integration system (Steps 8-9):

### AI Tools (Step 8)

Multi-select from 5 AI coding tools:

| Tool | What Gets Installed |
|------|-------------------|
| **Claude Code** | Binary + CLAUDE.md + Gentleman persona + 10+ skills + Kanagawa theme |
| **OpenCode** | Binary + Gentleman agent + SDD orchestrator + theme |
| **Gemini CLI** | CLI via npm |
| **GitHub Copilot** | gh extension |
| **Codex CLI** | Binary via npm + AGENTS.md config |

### AI Framework (Step 9)

Choose a preset or customize from **198 modules** across 6 categories:

| Category | Modules | Examples |
|----------|--------:|---------|
| 🪝 Hooks | 10 | Secret Scanner, Commit Guard, Model Router |
| ⚡ Commands | 20 | Git Commit, PR Review, TDD, Refactoring |
| 🤖 Agents | 72 | React Pro, DevOps Engineer, Security Auditor |
| 🎯 Skills | 85 | FastAPI, Spring Boot 4, Kubernetes, PyTorch |
| 📐 SDD | 2 | OpenSpec, Agent Teams Lite |
| 🔌 MCP | 9 | Context7, Engram, Jira, Atlassian, Figma, Notion, Brave Search, Sentry, Cloudflare |

**Presets**: Minimal, Frontend, Backend, Fullstack, Data, Complete

**SDD Choice**: Install [OpenSpec](https://github.com/JNZader/project-starter-framework) (file-based SDD), [Agent Teams Lite](https://github.com/Gentleman-Programming/agent-teams-lite) (lightweight SDD with 9 sub-agents), or both.

**Viewport Scrolling**: Long lists (Skills: 85, Agents: 72) scroll within the terminal with `▲`/`▼` indicators.

---

## 🎮 Vim Mastery Trainer

Learn Vim the fun way! The installer includes an interactive RPG-style trainer with:

| Module | Keys Covered |
|--------|--------------|
| 🔤 Horizontal Movement | `w`, `e`, `b`, `f`, `t`, `0`, `$`, `^` |
| ↕️ Vertical Movement | `j`, `k`, `G`, `gg`, `{`, `}` |
| 📦 Text Objects | `iw`, `aw`, `i"`, `a(`, `it`, `at` |
| ✂️ Change & Repeat | `d`, `c`, `dd`, `cc`, `D`, `C`, `x` |
| 🔄 Substitution | `r`, `R`, `s`, `S`, `~`, `gu`, `gU`, `J` |
| 🎬 Macros & Registers | `qa`, `@a`, `@@`, `"ay`, `"+p` |
| 🔍 Regex/Search | `/`, `?`, `n`, `N`, `*`, `#`, `\v` |

Each module includes 15 progressive lessons, practice mode with intelligent exercise selection, boss fights, and XP tracking.

Launch it from the main menu: **Vim Mastery Trainer**

---

## 📦 Project Initialization

Bootstrap any project with AI framework support:

```bash
# Interactive
gentleman-dots  # → Main Menu → Initialize Project

# Non-interactive
gentleman-dots --non-interactive --init-project \
  --project-path=/path/to/project \
  --project-memory=obsidian-brain \
  --project-ci=github --project-engram
```

**Memory modules**: Obsidian Brain, VibeKanban, Engram, Simple, None
**CI providers**: GitHub Actions, GitLab CI, Woodpecker, None

---

## 🎯 Skill Manager

Browse, install, and remove AI agent skills from the Gentleman-Skills catalog:

```bash
# Interactive
gentleman-dots  # → Main Menu → Skill Manager

# Non-interactive
gentleman-dots --non-interactive --skill-install=react-19,typescript,tailwind-4
gentleman-dots --non-interactive --skill-remove=react-19
```

Skills are organized by category (curated, community, plugin) and symlinked to `~/.claude/skills/`.

---

## 🔀 Fork Support

Override the clone URL and directory to point to your own fork:

```bash
# Via environment variables
REPO_URL=https://github.com/YourUser/YourFork.git REPO_DIR=YourFork gentleman-dots

# Via CLI flags
gentleman-dots --repo-url=https://github.com/YourUser/YourFork.git --repo-dir=YourFork
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [TUI Installer Guide](docs/tui-installer.md) | Interactive installer features, navigation, backup/restore |
| [AI Tools & Framework](docs/ai-tools-integration.md) | AI tools selection, framework presets, category drill-down, CLI flags |
| [AI Framework Modules](docs/ai-framework-modules.md) | Complete reference of all 198 modules across 6 categories + 28 curated skills |
| [Agent Teams Lite](docs/agent-teams-lite.md) | Lightweight SDD framework with 9 sub-agents |
| [AI Configuration](docs/ai-configuration.md) | Claude Code, OpenCode, Copilot, and other AI assistants |
| [Manual Installation](docs/manual-installation.md) | Step-by-step manual setup for all platforms |
| [Neovim Keymaps](docs/neovim-keymaps.md) | Complete reference of all keybindings |
| [Vim Trainer Spec](docs/vim-trainer-spec.md) | Technical specification for the Vim Mastery Trainer |
| [Docker Testing](docs/docker-testing.md) | E2E testing with Docker containers |
| [Contributing](docs/contributing.md) | Development setup, skills system, E2E tests, release process |

---

## Tools Overview

### Terminal Emulators

| Tool | Description |
|------|-------------|
| **Ghostty** | GPU-accelerated, native, blazing fast |
| **Kitty** | Feature-rich, GPU-based rendering |
| **WezTerm** | Lua-configurable, cross-platform |
| **Alacritty** | Minimal, Rust-based, lightweight |

### Shells

| Tool | Description |
|------|-------------|
| **Nushell** | Structured data, modern syntax, pipelines |
| **Fish** | User-friendly, great defaults, no config needed |
| **Zsh** | Highly customizable, POSIX-compatible, Powerlevel10k |

### Multiplexers

| Tool | Description |
|------|-------------|
| **Tmux** | Battle-tested, widely used, lots of plugins |
| **Zellij** | Modern, WebAssembly plugins, floating panes |

### Editors

| Tool | Description |
|------|-------------|
| **Neovim** | LazyVim config with LSP, completions, AI |
| **Zed** | High-performance editor with Vim mode and AI agent support |

### Prompts

| Tool | Description |
|------|-------------|
| **Starship** | Cross-shell prompt with Git integration |

---

## Project Structure

```
Javi.Dots/
├── installer/               # Go TUI installer source
│   ├── cmd/                 # Entry point
│   ├── internal/            # TUI, system, and trainer packages
│   └── e2e/                 # Docker-based E2E tests
├── docs/                    # Documentation
├── openspec/                # Spec-Driven Development artifacts
├── skills/                  # AI agent skills (repo-specific)
│
├── GentlemanNvim/           # Neovim configuration (LazyVim)
├── GentlemanClaude/         # Claude Code config + user skills
│   └── skills/              # Installable skills (React, Next.js, etc.)
├── GentlemanOpenCode/       # OpenCode AI config
├── GentlemanZed/            # Zed editor config (Vim mode + AI)
│
├── GentlemanFish/           # Fish shell config
├── GentlemanZsh/            # Zsh + Oh-My-Zsh + Powerlevel10k
├── GentlemanNushell/        # Nushell config
├── GentlemanTmux/           # Tmux config
├── GentlemanZellij/         # Zellij config
│
├── GentlemanGhostty/        # Ghostty terminal config
├── GentlemanKitty/          # Kitty terminal config
├── alacritty.toml           # Alacritty config
├── .wezterm.lua             # WezTerm config
│
└── starship.toml            # Starship prompt config
```

---

## Support

- **Issues**: [GitHub Issues](https://github.com/JNZader/Javi.Dots/issues)
- **Upstream**: [Gentleman.Dots](https://github.com/Gentleman-Programming/Gentleman.Dots) by Gentleman Programming
- **Discord**: [Gentleman Programming Community](https://discord.gg/gentleman-programming)
- **YouTube**: [@GentlemanProgramming](https://youtube.com/@GentlemanProgramming)

---

## License

MIT License - feel free to use, modify, and share.

**Happy coding!** 🎩
