package history

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/taps/internal/history"
	"github.com/meszmate/taps/internal/ui/styles"
)

type BackToMenuMsg struct{}

type Model struct {
	Styles  *styles.Styles
	Results []history.TestResult
	Stats   history.Stats
	cursor  int
	scroll  int
	width   int
	height  int
}

func New(s *styles.Styles) Model {
	results, _ := history.Load()
	stats := history.CalculateStats(results)

	return Model{
		Styles:  s,
		Results: results,
		Stats:   stats,
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
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scroll {
					m.scroll = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < len(m.Results)-1 {
				m.cursor++
				visibleLines := m.height - 12
				if visibleLines < 5 {
					visibleLines = 10
				}
				if m.cursor >= m.scroll+visibleLines {
					m.scroll = m.cursor - visibleLines + 1
				}
			}
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	t := m.Styles.Theme
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)
	b.WriteString(titleStyle.Render("History"))
	b.WriteString("\n\n")

	// Stats overview
	statLabel := lipgloss.NewStyle().Foreground(t.Sub)
	statValue := lipgloss.NewStyle().Foreground(t.Foreground).Bold(true)

	b.WriteString(statLabel.Render("total tests "))
	b.WriteString(statValue.Render(fmt.Sprintf("%d", m.Stats.TotalTests)))
	b.WriteString("  ")
	b.WriteString(statLabel.Render("avg wpm "))
	b.WriteString(statValue.Render(fmt.Sprintf("%.0f", m.Stats.AverageWPM)))
	b.WriteString("  ")
	b.WriteString(statLabel.Render("best wpm "))
	b.WriteString(statValue.Render(fmt.Sprintf("%.0f", m.Stats.BestWPM)))
	b.WriteString("  ")
	b.WriteString(statLabel.Render("last 10 avg "))
	b.WriteString(statValue.Render(fmt.Sprintf("%.0f", m.Stats.Last10Avg)))
	b.WriteString("\n\n")

	if len(m.Results) == 0 {
		dimStyle := lipgloss.NewStyle().Foreground(t.Sub)
		b.WriteString(dimStyle.Render("No test history yet. Complete a test to see results here."))
	} else {
		// Header
		headerStyle := lipgloss.NewStyle().Foreground(t.Sub).Bold(true)
		b.WriteString(headerStyle.Render(fmt.Sprintf("  %-12s %-8s %-8s %-10s %-10s %s", "Date", "Mode", "WPM", "Accuracy", "Consist.", "Config")))
		b.WriteString("\n")

		visibleLines := m.height - 12
		if visibleLines < 5 {
			visibleLines = 15
		}

		endIdx := m.scroll + visibleLines
		if endIdx > len(m.Results) {
			endIdx = len(m.Results)
		}

		// Show results in reverse chronological order
		reversed := make([]int, len(m.Results))
		for i := range m.Results {
			reversed[i] = len(m.Results) - 1 - i
		}

		for i := m.scroll; i < endIdx && i < len(reversed); i++ {
			ri := reversed[i]
			r := m.Results[ri]
			cursor := "  "
			lineStyle := lipgloss.NewStyle().Foreground(t.Sub)
			if i == m.cursor {
				cursor = "> "
				lineStyle = lipgloss.NewStyle().Foreground(t.Foreground)
			}

			date := r.Date.Format("01/02 15:04")
			cfgParts := []string{r.Language}
			if r.Punctuation {
				cfgParts = append(cfgParts, "punct")
			}
			if r.Numbers {
				cfgParts = append(cfgParts, "num")
			}

			line := fmt.Sprintf("%-12s %-8s %-8.0f %-10.1f%% %-10.1f%% %s",
				date, r.Mode, r.NetWPM, r.Accuracy, r.Consistency, strings.Join(cfgParts, ","))
			b.WriteString(lineStyle.Render(cursor + line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	helpStyle := lipgloss.NewStyle().Foreground(t.Sub)
	b.WriteString(helpStyle.Render("up/down scroll | esc back"))

	content := b.String()
	if m.width > 0 && m.height > 0 {
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}
