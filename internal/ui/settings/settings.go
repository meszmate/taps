package settings

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/taps/internal/config"
	"github.com/meszmate/taps/internal/ui/styles"
	"github.com/meszmate/taps/internal/ui/theme"
)

type BackToMenuMsg struct{}
type ThemeChangedMsg struct {
	ThemeName string
}

type settingType int

const (
	settingSelector settingType = iota
	settingToggle
)

type setting struct {
	label   string
	typ     settingType
	options []string
	getVal  func(c *config.Config) string
	setVal  func(c *config.Config, val string)
}

type Model struct {
	Config   *config.Config
	Styles   *styles.Styles
	cursor   int
	settings []setting
	width    int
	height   int
	scroll   int
}

func New(cfg *config.Config, s *styles.Styles) Model {
	settings := []setting{
		{
			label:   "Mode",
			typ:     settingSelector,
			options: []string{"time", "words", "quote", "zen"},
			getVal:  func(c *config.Config) string { return c.Mode },
			setVal:  func(c *config.Config, v string) { c.Mode = v },
		},
		{
			label:   "Time Duration",
			typ:     settingSelector,
			options: []string{"15", "30", "60", "120"},
			getVal:  func(c *config.Config) string { return fmt.Sprintf("%d", c.Duration) },
			setVal: func(c *config.Config, v string) {
				var d int
				fmt.Sscanf(v, "%d", &d)
				c.Duration = d
			},
		},
		{
			label:   "Word Count",
			typ:     settingSelector,
			options: []string{"10", "25", "50", "100"},
			getVal:  func(c *config.Config) string { return fmt.Sprintf("%d", c.WordCount) },
			setVal: func(c *config.Config, v string) {
				var w int
				fmt.Sscanf(v, "%d", &w)
				c.WordCount = w
			},
		},
		{
			label:   "Language",
			typ:     settingSelector,
			options: []string{"english", "english_1k"},
			getVal:  func(c *config.Config) string { return c.Language },
			setVal:  func(c *config.Config, v string) { c.Language = v },
		},
		{
			label:   "Punctuation",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.Punctuation {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.Punctuation = v == "on" },
		},
		{
			label:   "Numbers",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.Numbers {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.Numbers = v == "on" },
		},
		{
			label:   "Difficulty",
			typ:     settingSelector,
			options: []string{"normal", "expert", "master"},
			getVal:  func(c *config.Config) string { return c.Difficulty },
			setVal:  func(c *config.Config, v string) { c.Difficulty = v },
		},
		{
			label:   "Theme",
			typ:     settingSelector,
			options: theme.ThemeNames(),
			getVal:  func(c *config.Config) string { return c.Theme },
			setVal:  func(c *config.Config, v string) { c.Theme = v },
		},
		{
			label:   "Cursor Style",
			typ:     settingSelector,
			options: []string{"line", "block", "underline"},
			getVal:  func(c *config.Config) string { return c.CursorStyle },
			setVal:  func(c *config.Config, v string) { c.CursorStyle = v },
		},
		{
			label:   "Live WPM",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.LiveWPM {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.LiveWPM = v == "on" },
		},
		{
			label:   "Live Accuracy",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.LiveAccuracy {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.LiveAccuracy = v == "on" },
		},
		{
			label:   "Stop on Error",
			typ:     settingSelector,
			options: []string{"off", "word", "letter"},
			getVal:  func(c *config.Config) string { return c.StopOnError },
			setVal:  func(c *config.Config, v string) { c.StopOnError = v },
		},
		{
			label:   "Freedom Mode",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.FreedomMode {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.FreedomMode = v == "on" },
		},
		{
			label:   "Tape Mode",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.TapeMode {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.TapeMode = v == "on" },
		},
		{
			label:   "Show All Lines",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.ShowAllLines {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.ShowAllLines = v == "on" },
		},
		{
			label:   "Focus Mode",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.FocusMode {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.FocusMode = v == "on" },
		},
		{
			label:   "Sound on Error",
			typ:     settingToggle,
			options: []string{"off", "on"},
			getVal: func(c *config.Config) string {
				if c.SoundOnError {
					return "on"
				}
				return "off"
			},
			setVal: func(c *config.Config, v string) { c.SoundOnError = v == "on" },
		},
		{
			label:   "Quote Length",
			typ:     settingSelector,
			options: []string{"short", "medium", "long"},
			getVal:  func(c *config.Config) string { return c.QuoteLength },
			setVal:  func(c *config.Config, v string) { c.QuoteLength = v },
		},
	}

	return Model{
		Config:   cfg,
		Styles:   s,
		settings: settings,
	}
}

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
		case "esc", "q":
			_ = m.Config.Save()
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.settings)-1 {
				m.cursor++
				visibleLines := m.height - 6
				if visibleLines < 5 {
					visibleLines = 5
				}
				if m.cursor >= m.scroll+visibleLines {
					m.scroll = m.cursor - visibleLines + 1
				}
			}
		case "left", "h":
			m.cycleSetting(-1)
		case "right", "l", "enter":
			cmd := m.cycleSetting(1)
			return m, cmd
		}
	}
	return m, nil
}

func (m *Model) cycleSetting(dir int) tea.Cmd {
	s := m.settings[m.cursor]
	current := s.getVal(m.Config)
	idx := 0
	for i, opt := range s.options {
		if opt == current {
			idx = i
			break
		}
	}
	idx += dir
	if idx < 0 {
		idx = len(s.options) - 1
	}
	if idx >= len(s.options) {
		idx = 0
	}
	s.setVal(m.Config, s.options[idx])

	// If theme changed, notify
	if s.label == "Theme" {
		return func() tea.Msg { return ThemeChangedMsg{ThemeName: s.options[idx]} }
	}
	return nil
}

func (m Model) View() string {
	t := m.Styles.Theme
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)
	b.WriteString(titleStyle.Render("Settings"))
	b.WriteString("\n\n")

	visibleLines := m.height - 6
	if visibleLines < 5 {
		visibleLines = 20
	}

	endIdx := m.scroll + visibleLines
	if endIdx > len(m.settings) {
		endIdx = len(m.settings)
	}

	for i := m.scroll; i < endIdx; i++ {
		s := m.settings[i]
		labelStyle := lipgloss.NewStyle().Foreground(t.Sub).Width(18)
		cursor := "  "
		if i == m.cursor {
			labelStyle = lipgloss.NewStyle().Foreground(t.Main).Width(18).Bold(true)
			cursor = "> "
		}

		b.WriteString(cursor)
		b.WriteString(labelStyle.Render(s.label))
		b.WriteString("  ")

		currentVal := s.getVal(m.Config)
		for j, opt := range s.options {
			optStyle := lipgloss.NewStyle().Foreground(t.Sub)
			if opt == currentVal {
				optStyle = lipgloss.NewStyle().Foreground(t.Main).Bold(true)
			}
			b.WriteString(optStyle.Render(opt))
			if j < len(s.options)-1 {
				b.WriteString("  ")
			}
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(t.Sub)
	b.WriteString(helpStyle.Render("arrows navigate | left/right change | esc back"))

	content := b.String()
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
