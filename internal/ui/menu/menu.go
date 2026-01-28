package menu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/taps/internal/config"
	"github.com/meszmate/taps/internal/ui/styles"
	"github.com/meszmate/taps/internal/ui/theme"
)

type action int

const (
	actionStart action = iota
	actionSettings
	actionHistory
	actionQuit
)

type menuItem struct {
	label  string
	action action
}

type Model struct {
	Config     *config.Config
	Styles     *styles.Styles
	cursor     int
	items      []menuItem
	modeIdx    int
	modes      []string
	durIdx     int
	durations  []int
	wcIdx      int
	wordCounts []int
	width      int
	height     int
}

func New(cfg *config.Config, s *styles.Styles) Model {
	modes := []string{"time", "words", "quote", "zen"}
	durations := []int{15, 30, 60, 120}
	wordCounts := []int{10, 25, 50, 100}

	modeIdx := 0
	for i, m := range modes {
		if m == cfg.Mode {
			modeIdx = i
			break
		}
	}
	durIdx := 0
	for i, d := range durations {
		if d == cfg.Duration {
			durIdx = i
			break
		}
	}
	wcIdx := 0
	for i, w := range wordCounts {
		if w == cfg.WordCount {
			wcIdx = i
			break
		}
	}

	return Model{
		Config: cfg,
		Styles: s,
		items: []menuItem{
			{label: "Start Test", action: actionStart},
			{label: "Settings", action: actionSettings},
			{label: "History", action: actionHistory},
			{label: "Quit", action: actionQuit},
		},
		modes:      modes,
		modeIdx:    modeIdx,
		durations:  durations,
		durIdx:     durIdx,
		wordCounts: wordCounts,
		wcIdx:      wcIdx,
	}
}

type StartTestMsg struct {
	Mode        string
	Duration    int
	WordCount   int
	QuoteLength string
}
type OpenSettingsMsg struct{}
type OpenHistoryMsg struct{}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "left", "h":
			m.handleLeft()
		case "right", "l":
			m.handleRight()
		case "enter":
			return m, m.handleEnter()
		case "1":
			m.modeIdx = 0
			m.Config.Mode = m.modes[0]
		case "2":
			m.modeIdx = 1
			m.Config.Mode = m.modes[1]
		case "3":
			m.modeIdx = 2
			m.Config.Mode = m.modes[2]
		case "4":
			m.modeIdx = 3
			m.Config.Mode = m.modes[3]
		case "p":
			m.Config.Punctuation = !m.Config.Punctuation
		case "n":
			m.Config.Numbers = !m.Config.Numbers
		}
	}
	return m, nil
}

func (m *Model) handleLeft() {
	mode := m.modes[m.modeIdx]
	switch mode {
	case "time":
		if m.durIdx > 0 {
			m.durIdx--
			m.Config.Duration = m.durations[m.durIdx]
		}
	case "words":
		if m.wcIdx > 0 {
			m.wcIdx--
			m.Config.WordCount = m.wordCounts[m.wcIdx]
		}
	}
}

func (m *Model) handleRight() {
	mode := m.modes[m.modeIdx]
	switch mode {
	case "time":
		if m.durIdx < len(m.durations)-1 {
			m.durIdx++
			m.Config.Duration = m.durations[m.durIdx]
		}
	case "words":
		if m.wcIdx < len(m.wordCounts)-1 {
			m.wcIdx++
			m.Config.WordCount = m.wordCounts[m.wcIdx]
		}
	}
}

func (m Model) handleEnter() tea.Cmd {
	item := m.items[m.cursor]
	switch item.action {
	case actionStart:
		return func() tea.Msg {
			return StartTestMsg{
				Mode:        m.modes[m.modeIdx],
				Duration:    m.durations[m.durIdx],
				WordCount:   m.wordCounts[m.wcIdx],
				QuoteLength: m.Config.QuoteLength,
			}
		}
	case actionSettings:
		return func() tea.Msg { return OpenSettingsMsg{} }
	case actionHistory:
		return func() tea.Msg { return OpenHistoryMsg{} }
	case actionQuit:
		return tea.Quit
	}
	return nil
}

func (m Model) View() string {
	t := m.Styles.Theme
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(t.Main).
		Bold(true).
		MarginBottom(1)
	b.WriteString(titleStyle.Render(theme.RainbowText("Taps", string(t.Main), string(t.Caret))))
	b.WriteString("\n\n")

	// Mode selector
	modeLabel := lipgloss.NewStyle().Foreground(t.Sub).Render("mode  ")
	b.WriteString(modeLabel)
	for i, mode := range m.modes {
		style := lipgloss.NewStyle().Foreground(t.Sub)
		if i == m.modeIdx {
			style = lipgloss.NewStyle().Foreground(t.Main).Bold(true)
		}
		b.WriteString(style.Render(mode))
		if i < len(m.modes)-1 {
			b.WriteString("  ")
		}
	}
	b.WriteString("\n")

	// Duration / word count selector based on mode
	mode := m.modes[m.modeIdx]
	switch mode {
	case "time":
		durLabel := lipgloss.NewStyle().Foreground(t.Sub).Render("time  ")
		b.WriteString(durLabel)
		for i, d := range m.durations {
			style := lipgloss.NewStyle().Foreground(t.Sub)
			if i == m.durIdx {
				style = lipgloss.NewStyle().Foreground(t.Main).Bold(true)
			}
			b.WriteString(style.Render(fmt.Sprintf("%ds", d)))
			if i < len(m.durations)-1 {
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	case "words":
		wcLabel := lipgloss.NewStyle().Foreground(t.Sub).Render("words ")
		b.WriteString(wcLabel)
		for i, w := range m.wordCounts {
			style := lipgloss.NewStyle().Foreground(t.Sub)
			if i == m.wcIdx {
				style = lipgloss.NewStyle().Foreground(t.Main).Bold(true)
			}
			b.WriteString(style.Render(fmt.Sprintf("%d", w)))
			if i < len(m.wordCounts)-1 {
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	case "quote":
		ql := lipgloss.NewStyle().Foreground(t.Sub).Render("length ")
		b.WriteString(ql)
		for i, l := range []string{"short", "medium", "long"} {
			style := lipgloss.NewStyle().Foreground(t.Sub)
			if l == m.Config.QuoteLength {
				style = lipgloss.NewStyle().Foreground(t.Main).Bold(true)
			}
			b.WriteString(style.Render(l))
			if i < 2 {
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	}

	// Toggles
	b.WriteString("\n")
	punctStyle := lipgloss.NewStyle().Foreground(t.Sub)
	if m.Config.Punctuation {
		punctStyle = lipgloss.NewStyle().Foreground(t.Main)
	}
	b.WriteString(punctStyle.Render(fmt.Sprintf("@ punctuation %s", boolIcon(m.Config.Punctuation))))

	b.WriteString("  ")
	numStyle := lipgloss.NewStyle().Foreground(t.Sub)
	if m.Config.Numbers {
		numStyle = lipgloss.NewStyle().Foreground(t.Main)
	}
	b.WriteString(numStyle.Render(fmt.Sprintf("# numbers %s", boolIcon(m.Config.Numbers))))
	b.WriteString("\n\n")

	// Menu items
	for i, item := range m.items {
		style := lipgloss.NewStyle().Foreground(t.Sub)
		cursor := "  "
		if i == m.cursor {
			style = lipgloss.NewStyle().Foreground(t.Main).Bold(true)
			cursor = "> "
		}
		b.WriteString(style.Render(cursor + item.label))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(t.Sub)
	b.WriteString(helpStyle.Render("1-4 mode | arrows select | p punctuation | n numbers | enter confirm | ctrl+c quit"))

	// Center the content
	content := b.String()
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

func boolIcon(v bool) string {
	if v {
		return "on"
	}
	return "off"
}
