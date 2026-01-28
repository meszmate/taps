package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/meszmate/taps/internal/config"
	"github.com/meszmate/taps/internal/history"
	"github.com/meszmate/taps/internal/ui/menu"
	"github.com/meszmate/taps/internal/ui/results"
	"github.com/meszmate/taps/internal/ui/settings"
	"github.com/meszmate/taps/internal/ui/styles"
	"github.com/meszmate/taps/internal/ui/test"
	"github.com/meszmate/taps/internal/ui/theme"

	historyui "github.com/meszmate/taps/internal/ui/history"
)

type screen int

const (
	screenMenu screen = iota
	screenTest
	screenResults
	screenSettings
	screenHistory
)

type Model struct {
	screen     screen
	config     *config.Config
	styles     *styles.Styles
	theme      *theme.Theme
	menu       menu.Model
	test       test.Model
	results    results.Model
	settings   settings.Model
	history    historyui.Model
	windowSize tea.WindowSizeMsg
}

func New() Model {
	cfg := config.Load()
	t := theme.GetTheme(cfg.Theme)
	s := styles.New(t)

	m := Model{
		screen: screenMenu,
		config: cfg,
		styles: s,
		theme:  t,
		menu:   menu.New(cfg, s),
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.screen {
	case screenMenu:
		return m.updateMenu(msg)
	case screenTest:
		return m.updateTest(msg)
	case screenResults:
		return m.updateResults(msg)
	case screenSettings:
		return m.updateSettings(msg)
	case screenHistory:
		return m.updateHistory(msg)
	}
	return m, nil
}

// sendSize returns a command that sends the current window size to the new screen
func (m Model) sendSize() tea.Cmd {
	ws := m.windowSize
	return func() tea.Msg { return ws }
}

func (m Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.menu, cmd = m.menu.Update(msg)

	switch msg.(type) {
	case menu.StartTestMsg:
		stMsg := msg.(menu.StartTestMsg)
		m.test = test.New(m.config, m.styles, stMsg.Mode, stMsg.Duration, stMsg.WordCount, stMsg.QuoteLength)
		m.test.Width = m.windowSize.Width
		m.test.Height = m.windowSize.Height
		m.screen = screenTest
		return m, nil
	case menu.OpenSettingsMsg:
		m.settings = settings.New(m.config, m.styles)
		m.screen = screenSettings
		return m, m.sendSize()
	case menu.OpenHistoryMsg:
		m.history = historyui.New(m.styles)
		m.screen = screenHistory
		return m, m.sendSize()
	}

	return m, cmd
}

func (m Model) updateTest(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.test, cmd = m.test.Update(msg)

	switch msg := msg.(type) {
	case test.TestFinishedMsg:
		result := history.TestResult{
			Date:        msg.Engine.StartTime,
			Mode:        msg.Config.Mode,
			Duration:    msg.Config.Duration,
			WordCount:   msg.Config.WordCount,
			Language:    msg.Config.Language,
			Punctuation: msg.Config.Punctuation,
			Numbers:     msg.Config.Numbers,
			Difficulty:  msg.Config.Difficulty,
			NetWPM:      msg.Engine.CurrentNetWPM(),
			RawWPM:      msg.Engine.CurrentRawWPM(),
			Accuracy:    msg.Engine.CurrentAccuracy(),
			Consistency: msg.Engine.CurrentConsistency(),
			Correct:     msg.Engine.CorrectChars,
			Incorrect:   msg.Engine.IncorrectChars,
			Extra:       msg.Engine.ExtraChars,
			Missed:      msg.Engine.MissedChars,
			QuoteLength: msg.Config.QuoteLength,
		}
		_ = history.Append(result)

		tcfg := results.TestConfig{
			Mode:        msg.Config.Mode,
			Duration:    msg.Config.Duration,
			WordCount:   msg.Config.WordCount,
			Language:    msg.Config.Language,
			Punctuation: msg.Config.Punctuation,
			Numbers:     msg.Config.Numbers,
			Difficulty:  msg.Config.Difficulty,
			QuoteLength: msg.Config.QuoteLength,
		}
		m.results = results.New(m.styles, msg.Engine, msg.Mode, tcfg)
		m.results.Width = m.windowSize.Width
		m.results.Height = m.windowSize.Height
		m.screen = screenResults
		return m, nil

	case test.BackToMenuMsg:
		m.menu = menu.New(m.config, m.styles)
		m.screen = screenMenu
		return m, m.sendSize()
	}

	return m, cmd
}

func (m Model) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.results, cmd = m.results.Update(msg)

	switch msg := msg.(type) {
	case results.RestartMsg:
		m.test = test.New(m.config, m.styles, msg.Mode, msg.Duration, msg.WordCount, msg.QuoteLength)
		m.test.Width = m.windowSize.Width
		m.test.Height = m.windowSize.Height
		m.screen = screenTest
		return m, nil
	case results.NewTestMsg:
		m.menu = menu.New(m.config, m.styles)
		m.screen = screenMenu
		return m, m.sendSize()
	case results.BackToMenuMsg:
		m.menu = menu.New(m.config, m.styles)
		m.screen = screenMenu
		return m, m.sendSize()
	}

	return m, cmd
}

func (m Model) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.settings, cmd = m.settings.Update(msg)

	switch msg := msg.(type) {
	case settings.BackToMenuMsg:
		m.menu = menu.New(m.config, m.styles)
		m.screen = screenMenu
		return m, m.sendSize()
	case settings.ThemeChangedMsg:
		m.theme = theme.GetTheme(msg.ThemeName)
		m.styles = styles.New(m.theme)
		m.settings.Styles = m.styles
		return m, nil
	}

	return m, cmd
}

func (m Model) updateHistory(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.history, cmd = m.history.Update(msg)

	switch msg.(type) {
	case historyui.BackToMenuMsg:
		m.menu = menu.New(m.config, m.styles)
		m.screen = screenMenu
		return m, m.sendSize()
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.screen {
	case screenMenu:
		return m.menu.View()
	case screenTest:
		return m.test.View()
	case screenResults:
		return m.results.View()
	case screenSettings:
		return m.settings.View()
	case screenHistory:
		return m.history.View()
	}
	return ""
}
