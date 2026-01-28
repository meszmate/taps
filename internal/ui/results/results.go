package results

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/meszmate/taps/internal/typing"
	"github.com/meszmate/taps/internal/ui/styles"
	"github.com/meszmate/taps/internal/ui/theme"
)

type RestartMsg struct {
	Mode        string
	Duration    int
	WordCount   int
	QuoteLength string
}
type NewTestMsg struct{}
type BackToMenuMsg struct{}

type TestConfig struct {
	Mode        string
	Duration    int
	WordCount   int
	Language    string
	Punctuation bool
	Numbers     bool
	Difficulty  string
	QuoteLength string
}

type Model struct {
	Styles      *styles.Styles
	Engine      *typing.Engine
	Mode        string
	TCfg        TestConfig
	NetWPM      float64
	RawWPM      float64
	Accuracy    float64
	Consistency float64
	Width       int
	Height      int
}

func New(s *styles.Styles, engine *typing.Engine, mode string, tcfg TestConfig) Model {
	elapsed := engine.ElapsedSeconds()
	return Model{
		Styles:      s,
		Engine:      engine,
		Mode:        mode,
		TCfg:        tcfg,
		NetWPM:      typing.NetWPM(engine.CorrectChars, elapsed),
		RawWPM:      typing.RawWPM(engine.TotalTyped, elapsed),
		Accuracy:    typing.Accuracy(engine.CorrectChars, engine.IncorrectChars, engine.ExtraChars),
		Consistency: typing.Consistency(engine.PerSecondWPM),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			return m, func() tea.Msg {
				return RestartMsg{
					Mode:        m.TCfg.Mode,
					Duration:    m.TCfg.Duration,
					WordCount:   m.TCfg.WordCount,
					QuoteLength: m.TCfg.QuoteLength,
				}
			}
		case "enter":
			return m, func() tea.Msg { return NewTestMsg{} }
		case "esc":
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}
	return m, nil
}

func (m Model) View() string {
	t := m.Styles.Theme
	var b strings.Builder

	// WPM large display with rainbow
	wpmStr := fmt.Sprintf("%.0f", m.NetWPM)
	b.WriteString(theme.RainbowText(wpmStr, string(t.Main), string(t.Caret)))
	b.WriteString(" ")
	labelStyle := lipgloss.NewStyle().Foreground(t.Sub)
	b.WriteString(labelStyle.Render("wpm"))
	b.WriteString("\n\n")

	// Stats grid
	statLabel := lipgloss.NewStyle().Foreground(t.Sub)
	statValue := lipgloss.NewStyle().Foreground(t.Foreground).Bold(true)

	stats := []struct {
		label string
		value string
	}{
		{"raw", fmt.Sprintf("%.0f wpm", m.RawWPM)},
		{"accuracy", fmt.Sprintf("%.1f%%", m.Accuracy)},
		{"consistency", fmt.Sprintf("%.1f%%", m.Consistency)},
	}

	for _, s := range stats {
		b.WriteString(statLabel.Render(s.label+" "))
		b.WriteString(statValue.Render(s.value))
		b.WriteString("  ")
	}
	b.WriteString("\n\n")

	// Character breakdown
	charLabel := lipgloss.NewStyle().Foreground(t.Sub)
	correctStyle := lipgloss.NewStyle().Foreground(t.Correct)
	errorStyle := lipgloss.NewStyle().Foreground(t.Error)
	extraStyle := lipgloss.NewStyle().Foreground(t.ExtraError)
	missedStyle := lipgloss.NewStyle().Foreground(t.Sub)

	b.WriteString(charLabel.Render("characters  "))
	b.WriteString(correctStyle.Render(fmt.Sprintf("%d", m.Engine.CorrectChars)))
	b.WriteString(" / ")
	b.WriteString(errorStyle.Render(fmt.Sprintf("%d", m.Engine.IncorrectChars)))
	b.WriteString(" / ")
	b.WriteString(extraStyle.Render(fmt.Sprintf("%d", m.Engine.ExtraChars)))
	b.WriteString(" / ")
	b.WriteString(missedStyle.Render(fmt.Sprintf("%d", m.Engine.MissedChars)))
	b.WriteString("\n")
	b.WriteString(charLabel.Render("            correct / incorrect / extra / missed"))
	b.WriteString("\n\n")

	// WPM graph
	if len(m.Engine.PerSecondWPM) > 1 {
		graphWidth := m.Width - 20
		if graphWidth < 30 {
			graphWidth = 30
		}
		if graphWidth > 80 {
			graphWidth = 80
		}
		graph := asciigraph.Plot(m.Engine.PerSecondWPM,
			asciigraph.Width(graphWidth),
			asciigraph.Height(8),
			asciigraph.Caption("raw wpm over time"),
		)
		graphStyle := lipgloss.NewStyle().Foreground(t.Sub)
		b.WriteString(graphStyle.Render(graph))
		b.WriteString("\n\n")
	}

	// Test config summary
	cfgStyle := lipgloss.NewStyle().Foreground(t.Sub)
	cfgParts := []string{m.TCfg.Mode}
	if m.TCfg.Mode == "time" {
		cfgParts = append(cfgParts, fmt.Sprintf("%ds", m.TCfg.Duration))
	}
	if m.TCfg.Mode == "words" {
		cfgParts = append(cfgParts, fmt.Sprintf("%d words", m.TCfg.WordCount))
	}
	cfgParts = append(cfgParts, m.TCfg.Language)
	if m.TCfg.Punctuation {
		cfgParts = append(cfgParts, "punctuation")
	}
	if m.TCfg.Numbers {
		cfgParts = append(cfgParts, "numbers")
	}
	if m.TCfg.Difficulty != "normal" {
		cfgParts = append(cfgParts, m.TCfg.Difficulty)
	}
	b.WriteString(cfgStyle.Render(strings.Join(cfgParts, " | ")))
	b.WriteString("\n\n")

	// If failed
	if m.Engine.Failed {
		failStyle := lipgloss.NewStyle().Foreground(t.Error).Bold(true)
		b.WriteString(failStyle.Render("Test failed: "+m.Engine.FailedReason))
		b.WriteString("\n\n")
	}

	// Keybinds
	helpStyle := lipgloss.NewStyle().Foreground(t.Sub)
	b.WriteString(helpStyle.Render("tab restart | enter new test | esc menu"))

	content := b.String()
	if m.Width > 0 && m.Height > 0 {
		content = lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
