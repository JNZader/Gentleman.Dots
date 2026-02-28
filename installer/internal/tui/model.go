package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui/trainer"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents the current screen being displayed
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenMainMenu
	ScreenOSSelect
	ScreenTerminalSelect
	ScreenFontSelect
	ScreenShellSelect
	ScreenWMSelect
	ScreenNvimSelect
	ScreenInstalling
	ScreenComplete
	ScreenError
	// Learn screens
	ScreenLearnTerminals
	ScreenLearnShells
	ScreenLearnWM
	ScreenLearnNvim
	// Keymaps screen
	ScreenKeymaps
	ScreenKeymapCategory
	// Tool Keymaps screens
	ScreenKeymapsMenu       // Menu to select which tool's keymaps to view
	ScreenKeymapsTmux       // Tmux keymaps
	ScreenKeymapsTmuxCat    // Tmux keymap category
	ScreenKeymapsZellij     // Zellij keymaps
	ScreenKeymapsZellijCat  // Zellij keymap category
	ScreenKeymapsGhostty    // Ghostty keymaps
	ScreenKeymapsGhosttyCat // Ghostty keymap category
	// LazyVim learn screens
	ScreenLearnLazyVim
	ScreenLazyVimTopic
	// Backup screens
	ScreenBackupConfirm
	ScreenRestoreBackup
	ScreenRestoreConfirm
	// AI Framework screens
	ScreenAIToolsSelect      // Select which AI coding tools to install
	ScreenAIFrameworkConfirm // Confirm AI framework installation
	ScreenAIFrameworkPreset  // Select framework preset (minimal, frontend, etc.)
	ScreenAIFrameworkCategories  // Select module category to drill into
	ScreenAIFrameworkCategoryItems // Select individual items within a category
	// Warning screens
	ScreenGhosttyWarning // Warning about Ghostty compatibility on Debian/Ubuntu
	// Vim Trainer screens
	ScreenTrainerMenu       // Module selection
	ScreenTrainerLesson     // Lesson mode
	ScreenTrainerPractice   // Practice mode
	ScreenTrainerBoss       // Boss fight
	ScreenTrainerResult     // Result after exercise
	ScreenTrainerBossResult // Result after boss fight
	// Project Init screens
	ScreenProjectPath       // Text input: project directory
	ScreenProjectStack      // Single-select: detected stack confirmation/override
	ScreenProjectMemory     // Single-select: memory module
	ScreenProjectEngram     // Yes/No: add Engram alongside Obsidian Brain
	ScreenProjectCI         // Single-select: CI provider
	ScreenProjectConfirm    // Summary before execution
	ScreenProjectInstalling // Progress log
	ScreenProjectResult     // Success/error
	// Skill Manager screens
	ScreenSkillMenu    // Browse / Install / Remove
	ScreenSkillBrowse  // Scrollable read-only list
	ScreenSkillInstall // Multi-select from available skills
	ScreenSkillRemove  // Multi-select from installed skills
	ScreenSkillResult  // Success/error output
)

// InstallStep represents a single installation step
type InstallStep struct {
	ID          string
	Name        string
	Description string
	Status      StepStatus
	Progress    float64
	Error       error
	Interactive bool // If true, this step needs terminal control (sudo, chsh, etc)
}

type StepStatus int

const (
	StatusPending StepStatus = iota
	StatusRunning
	StatusDone
	StatusFailed
	StatusSkipped
)

// UserChoices stores all user selections
type UserChoices struct {
	OS           string // "mac", "linux"
	Terminal     string // "alacritty", "wezterm", "kitty", "ghostty", "none"
	InstallFont  bool
	Shell        string // "fish", "zsh", "nushell"
	WindowMgr    string // "tmux", "zellij", "none"
	InstallNvim  bool
	CreateBackup bool // Whether to backup existing configs
	// AI Tools and Framework
	AITools              []string // Selected AI tools: "claude", "opencode"
	InstallAIFramework   bool     // Whether to install project-starter-framework
	AIFrameworkPreset    string   // Preset: "minimal", "frontend", "backend", "fullstack", "data", "complete"
	AIFrameworkModules   []string // Individual module names when preset is "custom"
	InstallAgentTeamsLite bool    // Whether to install agent-teams-lite SDD framework
	// Project init
	InitProject   bool
	ProjectPath   string
	ProjectStack  string
	ProjectMemory string
	ProjectCI     string
	ProjectEngram bool
	// Skill manager
	SkillInstall []string
	SkillRemove  []string
}

// Model is the main application state
type Model struct {
	Screen      Screen
	PrevScreen  Screen // For going back from learn/keymaps screens
	Width       int
	Height      int
	SystemInfo  *system.SystemInfo
	Choices     UserChoices
	Steps       []InstallStep
	CurrentStep int
	Cursor      int
	ErrorMsg    string
	ShowDetails bool
	LogLines    []string
	TotalTime   float64
	Quitting    bool
	// Program reference for sending messages during installation
	Program *tea.Program
	// Spinner animation
	SpinnerFrame int
	// Learn mode
	ViewingTool string // Current tool being viewed in learn mode
	// Keymaps mode
	KeymapCategories []KeymapCategory
	SelectedCategory int
	KeymapScroll     int // For scrolling through keymaps
	// Tool-specific keymaps
	TmuxKeymapCategories    []KeymapCategory
	TmuxSelectedCategory    int
	TmuxKeymapScroll        int
	ZellijKeymapCategories  []KeymapCategory
	ZellijSelectedCategory  int
	ZellijKeymapScroll      int
	GhosttyKeymapCategories []KeymapCategory
	GhosttySelectedCategory int
	GhosttyKeymapScroll     int
	// LazyVim mode
	LazyVimTopics        []LazyVimTopic
	SelectedLazyVimTopic int
	LazyVimScroll        int // For scrolling through topic content
	// Backup mode
	ExistingConfigs  []string            // Configs that will be overwritten
	AvailableBackups []system.BackupInfo // Available backups for restore
	SelectedBackup   int                 // Selected backup index
	BackupDir        string              // Last backup directory created
	// Vim Trainer mode
	TrainerStats       *trainer.UserStats   // User's training stats
	TrainerGameState   *trainer.GameState   // Current game session state
	TrainerModules     []trainer.ModuleInfo // Available modules
	TrainerCursor      int                  // Cursor for module selection
	TrainerInput       string               // User's input for current exercise
	TrainerLastCorrect bool                 // Was last answer correct
	TrainerMessage     string               // Feedback message to display
	// AI Tools multi-select toggle
	AIToolSelected []bool // Toggle state for each tool in ScreenAIToolsSelect
	// AI Framework category drill-down selection
	AICategorySelected     map[string][]bool // Toggle state per category: categoryID → []bool for items
	SelectedModuleCategory int               // Index into moduleCategories for current drill-down
	CategoryItemsScroll    int               // Scroll offset for long item lists in category drill-down
	// Leader key mode (like Vim's <space> leader)
	LeaderMode bool // True when waiting for next key after <space>
	// Project init
	ProjectPathInput string
	ProjectPathError string
	ProjectStack     string
	ProjectMemory    string
	ProjectEngram    bool
	ProjectCI        string
	ProjectLogLines  []string
	// Skill manager
	SkillList      []string
	InstalledSkills []string
	SkillSelected  []bool
	SkillScroll    int
	SkillLoading   bool
	SkillLoadError string
	SkillResultLog []string
}

// NewModel creates a new Model with initial state
func NewModel() Model {
	return Model{
		Screen:                  ScreenWelcome,
		PrevScreen:              ScreenWelcome,
		Width:                   80,
		Height:                  24,
		SystemInfo:              system.Detect(),
		Choices:                 UserChoices{},
		Steps:                   []InstallStep{},
		CurrentStep:             0,
		Cursor:                  0,
		ShowDetails:             false,
		LogLines:                []string{},
		SpinnerFrame:            0,
		KeymapCategories:        GetNvimKeymaps(),
		SelectedCategory:        0,
		KeymapScroll:            0,
		TmuxKeymapCategories:    GetTmuxKeymaps(),
		TmuxSelectedCategory:    0,
		TmuxKeymapScroll:        0,
		ZellijKeymapCategories:  GetZellijKeymaps(),
		ZellijSelectedCategory:  0,
		ZellijKeymapScroll:      0,
		GhosttyKeymapCategories: GetGhosttyKeymaps(),
		GhosttySelectedCategory: 0,
		GhosttyKeymapScroll:     0,
		LazyVimTopics:           GetLazyVimTopics(),
		SelectedLazyVimTopic:    0,
		LazyVimScroll:           0,
		ExistingConfigs:         []string{},
		AvailableBackups:        []system.BackupInfo{},
		SelectedBackup:          0,
		BackupDir:               "",
		Program:                 nil, // Will be set after tea.Program is created
		// Trainer initialization
		TrainerStats:       nil, // Will be loaded when entering trainer
		TrainerGameState:   nil,
		TrainerModules:     trainer.GetAllModules(),
		TrainerCursor:      0,
		TrainerInput:       "",
		TrainerLastCorrect: false,
		TrainerMessage:     "",
		// Project init
		ProjectPathInput: "",
		ProjectPathError: "",
		ProjectStack:     "",
		ProjectMemory:    "",
		ProjectEngram:    false,
		ProjectCI:        "",
		ProjectLogLines:  []string{},
		// Skill manager
		SkillList:      []string{},
		InstalledSkills: []string{},
		SkillSelected:  []bool{},
		SkillScroll:    0,
		SkillLoading:   false,
		SkillLoadError: "",
		SkillResultLog: []string{},
	}
}

// SetProgram sets the tea.Program reference for sending messages during installation
func (m *Model) SetProgram(p *tea.Program) {
	m.Program = p
}

// globalProgram holds a reference to the tea.Program for sending logs during installation
var globalProgram *tea.Program

// SetGlobalProgram sets the global program reference
func SetGlobalProgram(p *tea.Program) {
	globalProgram = p
}

// nonInteractiveMode indicates if we're running without TUI
var nonInteractiveMode bool

// SetNonInteractiveMode enables or disables non-interactive mode
func SetNonInteractiveMode(enabled bool) {
	nonInteractiveMode = enabled
}

// SendLog sends a log message to the TUI during installation
func SendLog(stepID string, log string) {
	if nonInteractiveMode {
		// In non-interactive mode, print to stdout if verbose
		if os.Getenv("GENTLEMAN_VERBOSE") == "1" {
			fmt.Printf("    %s\n", log)
		}
		return
	}
	if globalProgram != nil {
		globalProgram.Send(stepProgressMsg{
			stepID: stepID,
			log:    log,
		})
	}
}

// SendLogLine is an alias for SendLog for compatibility
func (m *Model) SendLog(stepID string, log string) {
	SendLog(stepID, log)
}

// GetCurrentOptions returns the options for the current screen
func (m Model) GetCurrentOptions() []string {
	switch m.Screen {
	case ScreenMainMenu:
		opts := []string{
			"🚀 Start Installation",
			"📚 Learn About Tools",
			"⌨️  Keymaps Reference",
			"📖 LazyVim Guide",
			"🎮 Vim Trainer",
		}
		// Add restore option if backups exist
		if len(m.AvailableBackups) > 0 {
			opts = append(opts, "🔄 Restore from Backup")
		}
		opts = append(opts, "📦 Initialize Project")
		opts = append(opts, "🎯 Skill Manager")
		opts = append(opts, "❌ Exit")
		return opts
	case ScreenKeymapsMenu:
		return []string{"Neovim", "Tmux", "Zellij", "Ghostty", "─────────────", "← Back"}
	case ScreenOSSelect:
		macLabel := "macOS"
		linuxLabel := "Linux"
		termuxLabel := "Termux"
		if m.SystemInfo.OS == system.OSMac {
			macLabel = "macOS (detected)"
		} else if m.SystemInfo.OS == system.OSTermux {
			termuxLabel = "Termux (detected)"
		} else if m.SystemInfo.OS == system.OSLinux || m.SystemInfo.OS == system.OSArch || m.SystemInfo.OS == system.OSDebian || m.SystemInfo.OS == system.OSFedora {
			linuxLabel = "Linux (detected)"
		}
		return []string{macLabel, linuxLabel, termuxLabel}
	case ScreenTerminalSelect:
		alacrittyLabel := "Alacritty"
		// On Debian/Ubuntu, Alacritty needs to be built from source (PPAs are unreliable)
		// This applies to ALL Debian-based systems, not just ARM
		if m.SystemInfo != nil && (m.SystemInfo.OS == system.OSDebian || m.SystemInfo.OS == system.OSLinux) && m.Choices.OS == "linux" {
			alacrittyLabel = "Alacritty ⏱️  (builds from source, installs Rust ~5-10 min)"
		}
		if m.Choices.OS == "mac" {
			return []string{alacrittyLabel, "WezTerm", "Kitty", "Ghostty", "None", "─────────────", "ℹ️  Learn about terminals"}
		}
		return []string{alacrittyLabel, "WezTerm", "Ghostty", "None", "─────────────", "ℹ️  Learn about terminals"}
	case ScreenFontSelect:
		return []string{"Yes, install Iosevka Term Nerd Font", "No, I already have it"}
	case ScreenShellSelect:
		return []string{"Fish", "Zsh", "Nushell", "─────────────", "ℹ️  Learn about shells"}
	case ScreenWMSelect:
		return []string{"Tmux", "Zellij", "None", "─────────────", "ℹ️  Learn about multiplexers"}
	case ScreenNvimSelect:
		return []string{"Yes, install Neovim with config", "No, skip Neovim", "─────────────", "ℹ️  Learn about Neovim", "⌨️  View Keymaps", "📖 LazyVim Guide"}
	case ScreenAIToolsSelect:
		return []string{"Claude Code", "OpenCode", "Gemini CLI", "GitHub Copilot", "─────────────", "✅ Confirm selection"}
	case ScreenAIFrameworkConfirm:
		return []string{"Yes, install AI Framework", "No, skip framework"}
	case ScreenAIFrameworkPreset:
		return []string{
			"🔧 Custom — Pick individual modules",
			"─────────────",
			"🎯 Minimal — Core + git commands only",
			"🖥️  Frontend — React, Vue, Angular, testing, security hooks",
			"⚙️  Backend — APIs, databases, microservices, security hooks",
			"🔄 Fullstack — Frontend + Backend + infra + all commands",
			"📊 Data — Data engineering, ML/AI, analytics",
			"📦 Complete — Everything included",
		}
	case ScreenAIFrameworkCategories:
		opts := make([]string, 0, len(moduleCategories)+2)
		for i, cat := range moduleCategories {
			selected := 0
			total := len(cat.Items)
			if bools, ok := m.AICategorySelected[cat.ID]; ok {
				for _, b := range bools {
					if b {
						selected++
					}
				}
			}
			opts = append(opts, fmt.Sprintf("%s %s (%d/%d selected)", cat.Icon, cat.Label, selected, total))
			_ = i
		}
		opts = append(opts, "─────────────")
		opts = append(opts, "✅ Confirm selection")
		return opts
	case ScreenAIFrameworkCategoryItems:
		if m.SelectedModuleCategory < 0 || m.SelectedModuleCategory >= len(moduleCategories) {
			return []string{}
		}
		cat := moduleCategories[m.SelectedModuleCategory]
		opts := make([]string, 0, len(cat.Items)+2)
		for _, item := range cat.Items {
			opts = append(opts, item.Label)
		}
		opts = append(opts, "─────────────")
		opts = append(opts, "← Back")
		return opts
	case ScreenBackupConfirm:
		return []string{
			"✅ Install with Backup (recommended)",
			"⚠️  Install without Backup",
			"❌ Cancel",
		}
	case ScreenRestoreBackup:
		opts := make([]string, len(m.AvailableBackups)+2)
		for i, backup := range m.AvailableBackups {
			// Format: timestamp + file count
			opts[i] = fmt.Sprintf("%s (%d items)", backup.Timestamp.Format("2006-01-02 15:04:05"), len(backup.Files))
		}
		opts[len(m.AvailableBackups)] = "─────────────"
		opts[len(m.AvailableBackups)+1] = "← Back"
		return opts
	case ScreenRestoreConfirm:
		return []string{
			"✅ Yes, restore this backup",
			"🗑️  Delete this backup",
			"❌ Cancel",
		}
	case ScreenGhosttyWarning:
		return []string{
			"⚠️  Continue with Ghostty anyway",
			"🔄 Choose a different terminal",
			"❌ Cancel installation",
		}
	case ScreenLearnTerminals:
		return []string{"Alacritty", "WezTerm", "Kitty", "Ghostty", "─────────────", "← Back"}
	case ScreenLearnShells:
		return []string{"Fish", "Zsh", "Nushell", "─────────────", "← Back"}
	case ScreenLearnWM:
		return []string{"Tmux", "Zellij", "─────────────", "← Back"}
	case ScreenLearnNvim:
		return []string{"View Features", "View Keymaps", "📖 LazyVim Guide", "─────────────", "← Back"}
	case ScreenKeymaps:
		categories := make([]string, len(m.KeymapCategories)+2)
		for i, cat := range m.KeymapCategories {
			categories[i] = cat.Name
		}
		categories[len(m.KeymapCategories)] = "─────────────"
		categories[len(m.KeymapCategories)+1] = "← Back"
		return categories
	case ScreenKeymapsTmux:
		categories := make([]string, len(m.TmuxKeymapCategories)+2)
		for i, cat := range m.TmuxKeymapCategories {
			categories[i] = cat.Name
		}
		categories[len(m.TmuxKeymapCategories)] = "─────────────"
		categories[len(m.TmuxKeymapCategories)+1] = "← Back"
		return categories
	case ScreenKeymapsZellij:
		categories := make([]string, len(m.ZellijKeymapCategories)+2)
		for i, cat := range m.ZellijKeymapCategories {
			categories[i] = cat.Name
		}
		categories[len(m.ZellijKeymapCategories)] = "─────────────"
		categories[len(m.ZellijKeymapCategories)+1] = "← Back"
		return categories
	case ScreenKeymapsGhostty:
		categories := make([]string, len(m.GhosttyKeymapCategories)+2)
		for i, cat := range m.GhosttyKeymapCategories {
			categories[i] = cat.Name
		}
		categories[len(m.GhosttyKeymapCategories)] = "─────────────"
		categories[len(m.GhosttyKeymapCategories)+1] = "← Back"
		return categories
	case ScreenLearnLazyVim:
		titles := GetLazyVimTopicTitles()
		result := make([]string, len(titles)+2)
		copy(result, titles)
		result[len(titles)] = "─────────────"
		result[len(titles)+1] = "← Back"
		return result
	// Project Init screens
	case ScreenProjectStack:
		return []string{"Angular", "Node.js", "Go", "Python", "Rust", "Java", "Ruby", "PHP", "Other"}
	case ScreenProjectMemory:
		return []string{"🧠 Obsidian Brain", "📋 VibeKanban", "🧠 Engram", "📝 Simple", "❌ None"}
	case ScreenProjectEngram:
		return []string{"Yes, add Engram too", "No, just Obsidian Brain"}
	case ScreenProjectCI:
		return []string{"GitHub Actions", "GitLab CI", "Woodpecker", "None"}
	case ScreenProjectConfirm:
		return []string{"✅ Confirm & Initialize", "❌ Cancel"}
	// Skill Manager screens
	case ScreenSkillMenu:
		return []string{"🔍 Browse Skills", "📥 Install Skills", "🗑️  Remove Skills", "─────────────", "← Back"}
	case ScreenSkillBrowse:
		opts := make([]string, 0, len(m.SkillList)+2)
		for _, skill := range m.SkillList {
			opts = append(opts, skill)
		}
		opts = append(opts, "─────────────")
		opts = append(opts, "← Back")
		return opts
	case ScreenSkillInstall:
		opts := make([]string, 0, len(m.SkillList)+2)
		for _, skill := range m.SkillList {
			opts = append(opts, skill)
		}
		opts = append(opts, "─────────────")
		opts = append(opts, "✅ Confirm installation")
		return opts
	case ScreenSkillRemove:
		opts := make([]string, 0, len(m.InstalledSkills)+2)
		for _, skill := range m.InstalledSkills {
			opts = append(opts, skill)
		}
		opts = append(opts, "─────────────")
		opts = append(opts, "✅ Confirm removal")
		return opts
	default:
		return []string{}
	}
}

// GetScreenTitle returns the title for the current screen
func (m Model) GetScreenTitle() string {
	switch m.Screen {
	case ScreenWelcome:
		return "Welcome to Javi.Dots Installer"
	case ScreenMainMenu:
		return "Main Menu"
	case ScreenOSSelect:
		return "Step 1: Select Your Operating System"
	case ScreenTerminalSelect:
		return "Step 2: Choose Terminal Emulator"
	case ScreenFontSelect:
		return "Step 3: Nerd Font Installation"
	case ScreenShellSelect:
		return "Step 4: Choose Your Shell"
	case ScreenWMSelect:
		return "Step 5: Choose Window Manager"
	case ScreenNvimSelect:
		return "Step 6: Neovim Configuration"
	case ScreenAIToolsSelect:
		return "Step 7: AI Coding Tools"
	case ScreenAIFrameworkConfirm:
		return "Step 8: AI Framework"
	case ScreenAIFrameworkPreset:
		return "Step 8: Choose Framework Preset"
	case ScreenAIFrameworkCategories:
		return "Step 8: Select Module Categories"
	case ScreenAIFrameworkCategoryItems:
		if m.SelectedModuleCategory >= 0 && m.SelectedModuleCategory < len(moduleCategories) {
			cat := moduleCategories[m.SelectedModuleCategory]
			return fmt.Sprintf("Step 8: %s %s", cat.Icon, cat.Label)
		}
		return "Step 8: Select Modules"
	case ScreenBackupConfirm:
		return "⚠️  Existing Configs Detected"
	case ScreenRestoreBackup:
		return "🔄 Restore from Backup"
	case ScreenRestoreConfirm:
		return "🔄 Confirm Restore"
	case ScreenGhosttyWarning:
		return "⚠️  Ghostty Compatibility Warning"
	case ScreenInstalling:
		return "Installing..."
	case ScreenComplete:
		return "Installation Complete!"
	case ScreenError:
		return "Error"
	case ScreenLearnTerminals:
		return "📚 Learn: Terminal Emulators"
	case ScreenLearnShells:
		return "📚 Learn: Shells"
	case ScreenLearnWM:
		return "📚 Learn: Window Managers"
	case ScreenLearnNvim:
		return "📚 Learn: Neovim"
	case ScreenKeymaps:
		return "⌨️  Neovim Keymaps Reference"
	case ScreenKeymapCategory:
		if m.SelectedCategory < len(m.KeymapCategories) {
			return "⌨️  " + m.KeymapCategories[m.SelectedCategory].Name
		}
		return "⌨️  Keymaps"
	case ScreenKeymapsMenu:
		return "⌨️  Keymaps Reference"
	case ScreenKeymapsTmux:
		return "⌨️  Tmux Keymaps"
	case ScreenKeymapsTmuxCat:
		if m.TmuxSelectedCategory < len(m.TmuxKeymapCategories) {
			return "⌨️  " + m.TmuxKeymapCategories[m.TmuxSelectedCategory].Name
		}
		return "⌨️  Tmux Keymaps"
	case ScreenKeymapsZellij:
		return "⌨️  Zellij Keymaps"
	case ScreenKeymapsZellijCat:
		if m.ZellijSelectedCategory < len(m.ZellijKeymapCategories) {
			return "⌨️  " + m.ZellijKeymapCategories[m.ZellijSelectedCategory].Name
		}
		return "⌨️  Zellij Keymaps"
	case ScreenKeymapsGhostty:
		return "⌨️  Ghostty Keymaps"
	case ScreenKeymapsGhosttyCat:
		if m.GhosttySelectedCategory < len(m.GhosttyKeymapCategories) {
			return "⌨️  " + m.GhosttyKeymapCategories[m.GhosttySelectedCategory].Name
		}
		return "⌨️  Ghostty Keymaps"
	case ScreenLearnLazyVim:
		return "📖 LazyVim Guide"
	case ScreenLazyVimTopic:
		if m.SelectedLazyVimTopic < len(m.LazyVimTopics) {
			return "📖 " + m.LazyVimTopics[m.SelectedLazyVimTopic].Title
		}
		return "📖 LazyVim"
	case ScreenTrainerMenu:
		return "🎮 Vim Trainer - Module Selection"
	case ScreenTrainerLesson:
		return "🎮 Vim Trainer - Lesson"
	case ScreenTrainerPractice:
		return "🎮 Vim Trainer - Practice"
	case ScreenTrainerBoss:
		return "🎮 Vim Trainer - Boss Fight!"
	case ScreenTrainerResult:
		return "🎮 Vim Trainer - Result"
	case ScreenTrainerBossResult:
		return "🎮 Vim Trainer - Boss Battle Complete"
	// Project Init screens
	case ScreenProjectPath:
		return "📦 Initialize Project — Path"
	case ScreenProjectStack:
		return "📦 Initialize Project — Stack"
	case ScreenProjectMemory:
		return "📦 Initialize Project — Memory Module"
	case ScreenProjectEngram:
		return "📦 Initialize Project — Engram Add-on"
	case ScreenProjectCI:
		return "📦 Initialize Project — CI/CD Provider"
	case ScreenProjectConfirm:
		return "📦 Initialize Project — Confirm"
	case ScreenProjectInstalling:
		return "📦 Initializing Project..."
	case ScreenProjectResult:
		return "📦 Project Initialization Result"
	// Skill Manager screens
	case ScreenSkillMenu:
		return "🎯 Skill Manager"
	case ScreenSkillBrowse:
		return "🎯 Skill Manager — Browse"
	case ScreenSkillInstall:
		return "🎯 Skill Manager — Install"
	case ScreenSkillRemove:
		return "🎯 Skill Manager — Remove"
	case ScreenSkillResult:
		return "🎯 Skill Manager — Result"
	default:
		return ""
	}
}

// GetScreenDescription returns a description for the current screen
func (m Model) GetScreenDescription() string {
	switch m.Screen {
	case ScreenOSSelect:
		detected := m.SystemInfo.OSName
		if m.SystemInfo.IsWSL {
			detected += " (WSL)"
		}
		return "Detected: " + detected
	case ScreenTerminalSelect:
		if m.SystemInfo.IsWSL {
			return "Note: Terminal emulators should be installed on Windows for WSL"
		}
		return "Select your preferred terminal emulator"
	case ScreenFontSelect:
		return "Iosevka Term Nerd Font is required for icons and glyphs"
	case ScreenShellSelect:
		return "Current shell: " + m.SystemInfo.UserShell
	case ScreenWMSelect:
		return "Terminal multiplexer for managing sessions"
	case ScreenNvimSelect:
		return "Includes LSP, TreeSitter, and Gentleman config"
	case ScreenAIToolsSelect:
		return "Toggle tools with Enter. Confirm when ready."
	case ScreenAIFrameworkConfirm:
		return "Agents, skills, hooks, and commands for AI coding tools"
	case ScreenAIFrameworkPreset:
		return "Presets bundle agents, skills, hooks, and commands by role"
	case ScreenAIFrameworkCategories:
		return "Select a category to configure its modules"
	case ScreenAIFrameworkCategoryItems:
		return "Toggle modules with Enter. Press Esc to go back."
	case ScreenGhosttyWarning:
		return "Ghostty installation may fail on Ubuntu/Debian.\nThe installer script only supports certain versions."
	// Project Init screens
	case ScreenProjectPath:
		return "Enter the path to your project directory"
	case ScreenProjectStack:
		if m.ProjectStack != "" && m.ProjectStack != "unknown" {
			return "Auto-detected: " + m.ProjectStack
		}
		return "Select your project's tech stack"
	case ScreenProjectMemory:
		return "Choose an AI memory module for your project"
	case ScreenProjectEngram:
		return "Add Engram persistent memory alongside Obsidian Brain?"
	case ScreenProjectCI:
		return "Select CI/CD provider for your project"
	case ScreenProjectConfirm:
		return "Review your choices before initializing"
	case ScreenProjectInstalling:
		return "Running init-project.sh..."
	case ScreenProjectResult:
		return "Initialization complete"
	// Skill Manager screens
	case ScreenSkillMenu:
		return "Manage skills from the Gentleman-Skills catalog"
	case ScreenSkillBrowse:
		return "Available skills from the catalog"
	case ScreenSkillInstall:
		return "Toggle skills to install with Enter, then confirm"
	case ScreenSkillRemove:
		return "Toggle skills to remove with Enter, then confirm"
	case ScreenSkillResult:
		return "Operation results"
	default:
		return ""
	}
}

// SetupInstallSteps creates the installation steps based on user choices
func (m *Model) SetupInstallSteps() {
	m.Steps = []InstallStep{}

	// Backup step if user chose to backup (not interactive - just file copies)
	if m.Choices.CreateBackup && len(m.ExistingConfigs) > 0 {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "backup",
			Name:        "Backup Existing Configs",
			Description: "Creating backup of your current configuration",
			Status:      StatusPending,
		})
	}

	// Always clone repo first (not interactive - just git clone)
	m.Steps = append(m.Steps, InstallStep{
		ID:          "clone",
		Name:        "Clone Repository",
		Description: "Downloading Javi.Dots",
		Status:      StatusPending,
	})

	// Homebrew (interactive - first install needs password)
	// Skip for Termux - it uses pkg instead
	if !m.SystemInfo.HasBrew && !m.SystemInfo.IsTermux {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "homebrew",
			Name:        "Install Homebrew",
			Description: "Package manager",
			Status:      StatusPending,
			Interactive: true,
		})
	}

	// Dependencies based on OS
	// Check both Choices.OS and SystemInfo for Termux detection (redundancy)
	isTermux := m.Choices.OS == "termux" || m.SystemInfo.IsTermux
	if m.Choices.OS == "linux" && !isTermux {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "deps",
			Name:        "Install Dependencies",
			Description: "Base packages",
			Status:      StatusPending,
			Interactive: true, // Needs sudo
		})
	} else if isTermux {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "deps",
			Name:        "Install Dependencies",
			Description: "Base packages (pkg)",
			Status:      StatusPending,
			Interactive: false, // Termux doesn't need sudo
		})
	} else if m.Choices.OS == "mac" && !m.SystemInfo.HasXcode {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "xcode",
			Name:        "Install Xcode CLI",
			Description: "Developer tools",
			Status:      StatusPending,
		})
	}

	// Terminal
	if m.Choices.Terminal != "none" && m.Choices.Terminal != "" {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "terminal",
			Name:        "Install " + m.Choices.Terminal,
			Description: "Terminal emulator",
			Status:      StatusPending,
			Interactive: m.Choices.OS == "linux", // Linux needs sudo for pacman/apt
		})
	}

	// Font (not interactive - brew doesn't need password after installed)
	if m.Choices.InstallFont {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "font",
			Name:        "Install Iosevka Nerd Font",
			Description: "Nerd font with icons",
			Status:      StatusPending,
		})
	}

	// Shell (not interactive - brew doesn't need password)
	m.Steps = append(m.Steps, InstallStep{
		ID:          "shell",
		Name:        "Install " + m.Choices.Shell,
		Description: "Shell and plugins",
		Status:      StatusPending,
	})

	// Window manager (not interactive - brew doesn't need password)
	if m.Choices.WindowMgr != "none" && m.Choices.WindowMgr != "" {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "wm",
			Name:        "Install " + m.Choices.WindowMgr,
			Description: "Terminal multiplexer",
			Status:      StatusPending,
		})
	}

	// Neovim (not interactive - brew doesn't need password)
	if m.Choices.InstallNvim {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "nvim",
			Name:        "Install Neovim",
			Description: "Editor with config",
			Status:      StatusPending,
		})
	}

	// AI Tools: Claude Code + OpenCode (not interactive)
	if len(m.Choices.AITools) > 0 {
		toolNames := strings.Join(m.Choices.AITools, " + ")
		m.Steps = append(m.Steps, InstallStep{
			ID:          "aitools",
			Name:        "Install AI Tools",
			Description: toolNames,
			Status:      StatusPending,
		})
	}

	// AI Framework (not interactive)
	if m.Choices.InstallAIFramework {
		presetLabel := m.Choices.AIFrameworkPreset
		if presetLabel == "" {
			presetLabel = "custom"
		}
		m.Steps = append(m.Steps, InstallStep{
			ID:          "aiframework",
			Name:        "Install AI Framework",
			Description: "Preset: " + presetLabel,
			Status:      StatusPending,
		})
	}

	// Set default shell (interactive - chsh needs password)
	m.Steps = append(m.Steps, InstallStep{
		ID:          "setshell",
		Name:        "Set Default Shell",
		Description: "Configure default shell",
		Status:      StatusPending,
		Interactive: true,
	})

	// Cleanup (not interactive - just file deletion)
	m.Steps = append(m.Steps, InstallStep{
		ID:          "cleanup",
		Name:        "Cleanup",
		Description: "Removing temporary files",
		Status:      StatusPending,
	})
}
