package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "1.0.0"

// CLI flags for non-interactive mode
type cliFlags struct {
	version        bool
	help           bool
	test           bool
	dryRun         bool
	nonInteractive bool
	terminal       string
	shell          string
	windowMgr      string
	nvim           bool
	zed            bool
	font           bool
	backup         bool
	aiTools        string
	aiFramework    bool
	aiPreset       string
	aiModules      string
	agentTeamsLite bool
	initProject    bool
	projectPath    string
	projectMemory  string
	projectCI      string
	projectEngram  bool
	skillInstall   string // comma-separated skill names to install
	skillRemove    string // comma-separated skill names to remove
	repoDir        string // override repo directory name
	repoURL        string // override repo git URL
}

func parseFlags() *cliFlags {
	flags := &cliFlags{}

	flag.BoolVar(&flags.version, "version", false, "Show version information")
	flag.BoolVar(&flags.version, "v", false, "Show version information (shorthand)")
	flag.BoolVar(&flags.help, "help", false, "Show help message")
	flag.BoolVar(&flags.help, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&flags.test, "test", false, "Run in test mode (uses temporary directory)")
	flag.BoolVar(&flags.test, "t", false, "Run in test mode (shorthand)")
	flag.BoolVar(&flags.dryRun, "dry-run", false, "Show what would be installed without doing it")
	flag.BoolVar(&flags.nonInteractive, "non-interactive", false, "Run without TUI, use CLI flags")
	flag.StringVar(&flags.terminal, "terminal", "", "Terminal: alacritty, wezterm, kitty, ghostty, none")
	flag.StringVar(&flags.shell, "shell", "", "Shell: fish, zsh, nushell")
	flag.StringVar(&flags.windowMgr, "wm", "", "Window manager: tmux, zellij, none")
	flag.BoolVar(&flags.nvim, "nvim", false, "Install Neovim configuration")
	flag.BoolVar(&flags.zed, "zed", false, "Install Zed editor with config")
	flag.BoolVar(&flags.font, "font", false, "Install Nerd Font")
	flag.BoolVar(&flags.backup, "backup", true, "Backup existing configs (default: true)")
	flag.StringVar(&flags.aiTools, "ai-tools", "", "AI tools: claude,opencode,gemini,copilot,codex,qwen (comma-separated)")
	flag.BoolVar(&flags.aiFramework, "ai-framework", false, "Install AI coding framework")
	flag.StringVar(&flags.aiPreset, "ai-preset", "", "Framework preset: minimal, frontend, backend, fullstack, data, complete")
	flag.StringVar(&flags.aiModules, "ai-modules", "", "Framework features: hooks,commands,skills,agents,sdd,mcp (comma-separated)")
	flag.BoolVar(&flags.agentTeamsLite, "agent-teams-lite", false, "Install Agent Teams Lite SDD framework")
	flag.BoolVar(&flags.initProject, "init-project", false, "Initialize a project with AI framework")
	flag.StringVar(&flags.projectPath, "project-path", "", "Project directory path (required with --init-project)")
	flag.StringVar(&flags.projectMemory, "project-memory", "simple", "Memory module: obsidian-brain, vibekanban, engram, simple, none")
	flag.StringVar(&flags.projectCI, "project-ci", "none", "CI provider: github, gitlab, woodpecker, none")
	flag.BoolVar(&flags.projectEngram, "project-engram", false, "Add Engram alongside Obsidian Brain")
	flag.StringVar(&flags.skillInstall, "skill-install", "", "Skills to install (comma-separated)")
	flag.StringVar(&flags.skillRemove, "skill-remove", "", "Skills to remove (comma-separated)")
	flag.StringVar(&flags.repoDir, "repo-dir", "", "Override repo directory name (default: Gentleman.Dots, env: REPO_DIR)")
	flag.StringVar(&flags.repoURL, "repo-url", "", "Override repo git URL (default: upstream Gentleman.Dots, env: REPO_URL)")

	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()

	if flags.version {
		fmt.Printf("gentleman.dots v%s\n", Version)
		os.Exit(0)
	}

	if flags.help {
		printHelp()
		os.Exit(0)
	}

	if flags.test {
		setupTestMode()
	}

	if flags.dryRun {
		os.Setenv("GENTLEMAN_DRY_RUN", "1")
		fmt.Println("🧪 Dry-run mode: No actual installations will be performed")
	}

	// Non-interactive mode: run installation directly with provided flags
	if flags.nonInteractive {
		if err := runNonInteractive(flags); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Interactive TUI mode
	model := tui.NewModel()

	// Override repo dir: flag > env > default
	if flags.repoDir != "" {
		model.RepoDir = flags.repoDir
	} else if env := os.Getenv("REPO_DIR"); env != "" {
		model.RepoDir = env
	}

	// Override repo URL: flag > env > default
	if flags.repoURL != "" {
		model.RepoURL = flags.repoURL
	} else if env := os.Getenv("REPO_URL"); env != "" {
		model.RepoURL = env
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	tui.SetGlobalProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running installer: %v\n", err)
		os.Exit(1)
	}
}

func runNonInteractive(flags *cliFlags) error {
	// Handle project init
	if flags.initProject {
		if flags.projectPath == "" {
			return fmt.Errorf("--project-path is required with --init-project")
		}
		path := tui.ExpandPath(flags.projectPath)
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("invalid project path: %w", err)
		}
		info, err := os.Stat(absPath)
		if err != nil {
			return fmt.Errorf("project path does not exist: %s", absPath)
		}
		if !info.IsDir() {
			return fmt.Errorf("project path is not a directory: %s", absPath)
		}

		// Validate memory
		validMemory := map[string]bool{"obsidian-brain": true, "vibekanban": true, "engram": true, "simple": true, "none": true}
		memory := strings.ToLower(flags.projectMemory)
		if !validMemory[memory] {
			return fmt.Errorf("invalid memory module: %s (valid: obsidian-brain, vibekanban, engram, simple, none)", memory)
		}

		// Validate CI
		validCI := map[string]bool{"github": true, "gitlab": true, "woodpecker": true, "none": true}
		ci := strings.ToLower(flags.projectCI)
		if !validCI[ci] {
			return fmt.Errorf("invalid CI provider: %s (valid: github, gitlab, woodpecker, none)", ci)
		}

		// Validate engram requires obsidian-brain
		if flags.projectEngram && memory != "obsidian-brain" {
			return fmt.Errorf("--project-engram requires --project-memory=obsidian-brain")
		}

		fmt.Println("📦 Initializing project...")
		fmt.Printf("  Path:    %s\n", absPath)
		fmt.Printf("  Memory:  %s\n", memory)
		fmt.Printf("  CI:      %s\n", ci)
		if flags.projectEngram {
			fmt.Printf("  Engram:  yes\n")
		}
		fmt.Println()

		tui.SetNonInteractiveMode(true)
		if err := tui.RunProjectInitScript(absPath, memory, ci, flags.projectEngram); err != nil {
			return fmt.Errorf("project initialization failed: %w", err)
		}
		fmt.Println("✅ Project initialized successfully!")
		return nil // Don't continue to environment installation
	}

	// Handle skill operations
	if flags.skillInstall != "" {
		names := strings.Split(flags.skillInstall, ",")
		for i := range names {
			names[i] = strings.TrimSpace(names[i])
		}
		fmt.Printf("📥 Installing %d skill(s)...\n", len(names))
		tui.SetNonInteractiveMode(true)

		// Fetch catalog to get SkillInfo for each requested name
		catalog, err := tui.FetchSkillCatalog()
		if err != nil {
			return fmt.Errorf("failed to fetch skill catalog: %w", err)
		}
		var toInstall []tui.SkillInfo
		nameSet := make(map[string]bool)
		for _, n := range names {
			nameSet[n] = true
		}
		for _, s := range catalog {
			if nameSet[s.Name] || nameSet[s.DirName] {
				toInstall = append(toInstall, s)
			}
		}
		if len(toInstall) == 0 {
			return fmt.Errorf("no matching skills found in catalog for: %s", strings.Join(names, ", "))
		}

		logLines, err := tui.InstallSkillSymlinks(toInstall)
		for _, line := range logLines {
			fmt.Println("  " + line)
		}
		if err != nil {
			return fmt.Errorf("skill installation: %w", err)
		}
		fmt.Println("✅ Skills installed!")
		if flags.shell == "" {
			return nil // Only skill operation, no env install
		}
	}

	if flags.skillRemove != "" {
		names := strings.Split(flags.skillRemove, ",")
		for i := range names {
			names[i] = strings.TrimSpace(names[i])
		}
		fmt.Printf("🗑️  Removing %d skill(s)...\n", len(names))
		tui.SetNonInteractiveMode(true)

		// Build SkillInfo from names (we only need the Name field for removal)
		var toRemove []tui.SkillInfo
		for _, n := range names {
			toRemove = append(toRemove, tui.SkillInfo{Name: n})
		}

		logLines, err := tui.RemoveSkillSymlinks(toRemove)
		for _, line := range logLines {
			fmt.Println("  " + line)
		}
		if err != nil {
			return fmt.Errorf("skill removal: %w", err)
		}
		fmt.Println("✅ Skills removed!")
		if flags.shell == "" {
			return nil // Only skill operation, no env install
		}
	}

	// Validate required flags for environment installation
	if flags.shell == "" {
		return fmt.Errorf("--shell is required (fish, zsh, nushell)")
	}

	// Normalize inputs
	terminal := strings.ToLower(flags.terminal)
	shell := strings.ToLower(flags.shell)
	wm := strings.ToLower(flags.windowMgr)

	// Validate values
	validTerminals := map[string]bool{"alacritty": true, "wezterm": true, "kitty": true, "ghostty": true, "none": true, "": true}
	validShells := map[string]bool{"fish": true, "zsh": true, "nushell": true}
	validWMs := map[string]bool{"tmux": true, "zellij": true, "none": true, "": true}

	if !validTerminals[terminal] {
		return fmt.Errorf("invalid terminal: %s (valid: alacritty, wezterm, kitty, ghostty, none)", terminal)
	}
	if !validShells[shell] {
		return fmt.Errorf("invalid shell: %s (valid: fish, zsh, nushell)", shell)
	}
	if !validWMs[wm] {
		return fmt.Errorf("invalid window manager: %s (valid: tmux, zellij, none)", wm)
	}

	// Default empty values to "none"
	if terminal == "" {
		terminal = "none"
	}
	if wm == "" {
		wm = "none"
	}

	// Parse AI tools
	var aiTools []string
	if flags.aiTools != "" {
		validAITools := map[string]bool{"claude": true, "opencode": true, "gemini": true, "copilot": true, "codex": true, "qwen": true}
		for _, tool := range strings.Split(flags.aiTools, ",") {
			tool = strings.TrimSpace(strings.ToLower(tool))
			if !validAITools[tool] {
				return fmt.Errorf("invalid AI tool: %s (valid: claude, opencode, gemini, copilot, codex, qwen)", tool)
			}
			aiTools = append(aiTools, tool)
		}
	}

	// Validate AI preset
	aiPreset := strings.ToLower(flags.aiPreset)
	validPresets := map[string]bool{"minimal": true, "frontend": true, "backend": true, "fullstack": true, "data": true, "complete": true, "": true}
	if !validPresets[aiPreset] {
		return fmt.Errorf("invalid AI preset: %s (valid: minimal, frontend, backend, fullstack, data, complete)", aiPreset)
	}

	// Parse AI features (modules flag accepts feature names: hooks,commands,skills,agents,sdd,mcp)
	var aiModules []string
	if flags.aiModules != "" {
		validFeatures := map[string]bool{"hooks": true, "commands": true, "skills": true, "agents": true, "sdd": true, "mcp": true}
		for _, mod := range strings.Split(flags.aiModules, ",") {
			mod = strings.TrimSpace(strings.ToLower(mod))
			if mod == "" {
				continue
			}
			if !validFeatures[mod] {
				return fmt.Errorf("invalid AI feature: %s (valid: hooks, commands, skills, agents, sdd, mcp)", mod)
			}
			aiModules = append(aiModules, mod)
		}
	}

	// Determine if framework should be installed
	installFramework := flags.aiFramework || aiPreset != "" || len(aiModules) > 0 || flags.agentTeamsLite

	// Create choices
	choices := tui.UserChoices{
		Terminal:              terminal,
		Shell:                 shell,
		WindowMgr:             wm,
		InstallNvim:           flags.nvim,
		InstallZed:            flags.zed,
		InstallFont:           flags.font,
		CreateBackup:          flags.backup,
		AITools:               aiTools,
		InstallAIFramework:    installFramework,
		AIFrameworkPreset:     aiPreset,
		AIFrameworkModules:    aiModules,
		InstallAgentTeamsLite: flags.agentTeamsLite,
	}

	fmt.Println("🚀 Javi.Dots Non-Interactive Installer")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Terminal:    %s\n", choices.Terminal)
	fmt.Printf("  Shell:       %s\n", choices.Shell)
	fmt.Printf("  Window Mgr:  %s\n", choices.WindowMgr)
	fmt.Printf("  Neovim:      %v\n", choices.InstallNvim)
	fmt.Printf("  Zed:         %v\n", choices.InstallZed)
	fmt.Printf("  Font:        %v\n", choices.InstallFont)
	fmt.Printf("  Backup:      %v\n", choices.CreateBackup)
	if len(choices.AITools) > 0 {
		fmt.Printf("  AI Tools:    %s\n", strings.Join(choices.AITools, ", "))
	}
	if choices.InstallAIFramework {
		if choices.AIFrameworkPreset != "" {
			fmt.Printf("  AI Framework: preset=%s\n", choices.AIFrameworkPreset)
		} else if len(choices.AIFrameworkModules) > 0 {
			fmt.Printf("  AI Framework: modules=%s\n", strings.Join(choices.AIFrameworkModules, ","))
		} else {
			fmt.Printf("  AI Framework: yes\n")
		}
		if choices.InstallAgentTeamsLite {
			fmt.Printf("  Agent Teams:  yes\n")
		}
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Resolve repo dir: flag > env > default
	repoDir := tui.DefaultRepoDir
	if flags.repoDir != "" {
		repoDir = flags.repoDir
	} else if env := os.Getenv("REPO_DIR"); env != "" {
		repoDir = env
	}

	// Resolve repo URL: flag > env > default
	repoURL := tui.DefaultRepoURL
	if flags.repoURL != "" {
		repoURL = flags.repoURL
	} else if env := os.Getenv("REPO_URL"); env != "" {
		repoURL = env
	}

	// Run the installation
	return tui.RunNonInteractive(choices, repoDir, repoURL)
}

func setupTestMode() {
	// Create a temporary test directory
	testDir := filepath.Join(os.TempDir(), "gentleman-dots-test")
	testHome := filepath.Join(testDir, "home")
	testConfig := filepath.Join(testHome, ".config")

	// Create directories
	os.MkdirAll(testConfig, 0755)

	// Override HOME to use test directory
	os.Setenv("HOME", testHome)
	os.Setenv("GENTLEMAN_TEST_MODE", "1")

	fmt.Printf("🧪 Test mode enabled!\n")
	fmt.Printf("   Test HOME: %s\n", testHome)
	fmt.Printf("   Your real configs are SAFE.\n")
	fmt.Printf("   Press Enter to continue...\n")

	// Wait for user to acknowledge
	fmt.Scanln()
}

func printHelp() {
	fmt.Println(`javi.dots - TUI installer for Javi.Dots terminal environment (fork of Gentleman.Dots)

Usage:
  gentleman.dots [flags]

Interactive Mode (default):
  Just run 'gentleman.dots' to start the TUI installer.

Non-Interactive Mode:
  gentleman.dots --non-interactive --shell=<shell> [options]

Flags:
  -h, --help           Show this help message
  -v, --version        Show version information
  -t, --test           Run in test mode (uses temporary directory)
  --dry-run            Show what would be installed without doing it
  --non-interactive    Run without TUI, use CLI flags instead

Non-Interactive Options:
  --repo-dir=<dir>     Override repo directory name (default: Gentleman.Dots, env: REPO_DIR)
  --repo-url=<url>     Override repo git URL (default: upstream Gentleman.Dots, env: REPO_URL)
  --shell=<shell>      Shell to install (required): fish, zsh, nushell
  --terminal=<term>    Terminal: alacritty, wezterm, kitty, ghostty, none
  --wm=<wm>            Window manager: tmux, zellij, none
  --nvim               Install Neovim configuration
  --zed                Install Zed editor with config
  --font               Install Nerd Font
  --backup=false       Disable config backup (default: true)

AI Options:
  --ai-tools=<tools>   AI tools (comma-separated): claude, opencode, gemini, copilot, codex, qwen
  --ai-framework       Install AI coding framework
  --ai-preset=<name>   Framework preset: minimal, frontend, backend, fullstack, data, complete
  --ai-modules=<feats> Framework features (comma-separated): hooks, commands, skills, agents, sdd, mcp
                       Each feature installs ALL items in that category (91 agents, 85 skills, etc.)
  --agent-teams-lite   Install Agent Teams Lite SDD framework (can combine with --ai-modules=sdd for both)

Project Init Options:
  --init-project       Initialize a project with AI framework
  --project-path=<dir> Project directory (required with --init-project)
  --project-memory=<m> Memory module: obsidian-brain, vibekanban, engram, simple, none (default: simple)
  --project-ci=<ci>    CI provider: github, gitlab, woodpecker, none (default: none)
  --project-engram     Add Engram alongside Obsidian Brain

Skill Manager Options:
  --skill-install=<s>  Skills to install (comma-separated names)
  --skill-remove=<s>   Skills to remove (comma-separated names)

Examples:
  # Interactive TUI
  gentleman.dots

  # Non-interactive with Fish + Zellij + Neovim
  gentleman.dots --non-interactive --shell=fish --wm=zellij --nvim

  # Full setup with AI tools and framework preset
  gentleman.dots --non-interactive --shell=fish --nvim --ai-tools=claude,opencode,gemini,copilot,codex --ai-preset=fullstack

  # Custom feature selection with Agent Teams Lite
  gentleman.dots --non-interactive --shell=zsh --ai-tools=claude --ai-framework \
    --ai-modules=hooks,skills --agent-teams-lite

  # Test mode with Zsh + Tmux (no terminal, no nvim)
  gentleman.dots --test --non-interactive --shell=zsh --wm=tmux

  # Initialize a project
  gentleman.dots --non-interactive --init-project --project-path=/path/to/project --project-memory=obsidian-brain --project-ci=github

  # Install skills
  gentleman.dots --non-interactive --skill-install=react-19,typescript,tailwind-4

  # Remove skills
  gentleman.dots --non-interactive --skill-remove=react-19

  # Verbose output (shows all command logs)
  GENTLEMAN_VERBOSE=1 gentleman.dots --non-interactive --shell=fish --nvim

Navigation (TUI mode):
  ↑/k, ↓/j        Navigate up/down
  Enter/Space     Select option
  Esc             Go back
  q               Quit
  d               Toggle details (during installation)

For more info: https://github.com/Gentleman-Programming/Gentleman.Dots`)
}
