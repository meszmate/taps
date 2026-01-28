package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/taps/internal/ui/theme"
)

type Styles struct {
	Theme *theme.Theme

	// Base styles
	App          lipgloss.Style
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	Highlight    lipgloss.Style
	Dim          lipgloss.Style
	Error        lipgloss.Style
	Success      lipgloss.Style
	Selected     lipgloss.Style
	Unselected   lipgloss.Style
	Key          lipgloss.Style
	Value        lipgloss.Style
	Border       lipgloss.Style
	StatusBar    lipgloss.Style
	BigNumber    lipgloss.Style
	Label        lipgloss.Style
}

func New(t *theme.Theme) *Styles {
	return &Styles{
		Theme: t,
		App: lipgloss.NewStyle().
			Background(t.Background).
			Foreground(t.Foreground),
		Title: lipgloss.NewStyle().
			Foreground(t.Main).
			Bold(true),
		Subtitle: lipgloss.NewStyle().
			Foreground(t.Sub),
		Highlight: lipgloss.NewStyle().
			Foreground(t.Main).
			Bold(true),
		Dim: lipgloss.NewStyle().
			Foreground(t.Sub),
		Error: lipgloss.NewStyle().
			Foreground(t.Error),
		Success: lipgloss.NewStyle().
			Foreground(t.Correct),
		Selected: lipgloss.NewStyle().
			Foreground(t.Main).
			Bold(true),
		Unselected: lipgloss.NewStyle().
			Foreground(t.Sub),
		Key: lipgloss.NewStyle().
			Foreground(t.Main),
		Value: lipgloss.NewStyle().
			Foreground(t.Foreground),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Sub).
			Padding(1, 2),
		StatusBar: lipgloss.NewStyle().
			Foreground(t.Sub).
			Padding(0, 1),
		BigNumber: lipgloss.NewStyle().
			Foreground(t.Main).
			Bold(true),
		Label: lipgloss.NewStyle().
			Foreground(t.Sub),
	}
}
