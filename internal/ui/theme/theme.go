package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name       string
	Background lipgloss.Color
	Foreground lipgloss.Color
	Sub        lipgloss.Color // subdued text
	Main       lipgloss.Color // accent/highlight
	Caret      lipgloss.Color
	Correct    lipgloss.Color
	Error      lipgloss.Color
	ExtraError lipgloss.Color
}

func (t *Theme) BgColor() string    { return string(t.Background) }
func (t *Theme) FgColor() string    { return string(t.Foreground) }
func (t *Theme) SubColor() string   { return string(t.Sub) }
func (t *Theme) MainColor() string  { return string(t.Main) }
func (t *Theme) CaretColor() string { return string(t.Caret) }

func GetTheme(name string) *Theme {
	if t, ok := ThemeCatalog[name]; ok {
		return t
	}
	return ThemeCatalog["default_dark"]
}

func CustomTheme(name, bg, fg, sub, main, caret, correct, err, extra string) *Theme {
	return &Theme{
		Name:       name,
		Background: lipgloss.Color(bg),
		Foreground: lipgloss.Color(fg),
		Sub:        lipgloss.Color(sub),
		Main:       lipgloss.Color(main),
		Caret:      lipgloss.Color(caret),
		Correct:    lipgloss.Color(correct),
		Error:      lipgloss.Color(err),
		ExtraError: lipgloss.Color(extra),
	}
}

func ThemeNames() []string {
	return []string{
		"default_dark",
		"dracula",
		"nord",
		"gruvbox_dark",
		"catppuccin_mocha",
		"solarized_dark",
		"tokyo_night",
		"one_dark",
		"rose_pine",
		"serika_dark",
	}
}
