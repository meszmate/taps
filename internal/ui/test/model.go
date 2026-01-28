package test

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/meszmate/taps/internal/config"
	"github.com/meszmate/taps/internal/typing"
	"github.com/meszmate/taps/internal/ui/styles"
)

type TickMsg time.Time
type WPMSampleMsg time.Time

type TestFinishedMsg struct {
	Engine *typing.Engine
	Mode   string
	Config TestConfig
}

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

type BackToMenuMsg struct{}

type Model struct {
	Config  *config.Config
	Styles  *styles.Styles
	Engine  *typing.Engine
	Mode    string
	TCfg    TestConfig
	Timer   int // remaining seconds (time mode)
	Width   int
	Height  int
	started bool
}

func New(cfg *config.Config, s *styles.Styles, mode string, duration, wordCount int, quoteLength string) Model {
	var target string
	tcfg := TestConfig{
		Mode:        mode,
		Duration:    duration,
		WordCount:   wordCount,
		Language:    cfg.Language,
		Punctuation: cfg.Punctuation,
		Numbers:     cfg.Numbers,
		Difficulty:  cfg.Difficulty,
		QuoteLength: quoteLength,
	}

	switch mode {
	case "time":
		target = typing.GenerateWordsForTime(cfg.Language, cfg.Punctuation, cfg.Numbers)
	case "words":
		target = typing.GenerateWords(wordCount, cfg.Language, cfg.Punctuation, cfg.Numbers)
	case "quote":
		q := typing.GetRandomQuote(quoteLength)
		target = q.Text
	case "zen":
		target = typing.GenerateWordsForTime(cfg.Language, cfg.Punctuation, cfg.Numbers)
	default:
		target = typing.GenerateWords(50, cfg.Language, cfg.Punctuation, cfg.Numbers)
	}

	engine := typing.NewEngine(target, cfg.StopOnError, cfg.FreedomMode, cfg.Difficulty)

	return Model{
		Config: cfg,
		Styles: s,
		Engine: engine,
		Mode:   mode,
		TCfg:   tcfg,
		Timer:  duration,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func wpmSampleCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return WPMSampleMsg(t)
	})
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case TickMsg:
		if !m.Engine.Started || m.Engine.Finished {
			return m, nil
		}
		if m.Mode == "time" {
			m.Timer--
			if m.Timer <= 0 {
				m.Engine.Finish()
				return m, m.finishCmd()
			}
			return m, tickCmd()
		}

	case WPMSampleMsg:
		if m.Engine.Started && !m.Engine.Finished {
			m.Engine.SampleWPM()
			return m, wpmSampleCmd()
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return BackToMenuMsg{} }
		case "tab":
			// Quick restart
			newM := New(m.Config, m.Styles, m.Mode, m.TCfg.Duration, m.TCfg.WordCount, m.TCfg.QuoteLength)
			newM.Width = m.Width
			newM.Height = m.Height
			return newM, nil
		case "backspace", "ctrl+h":
			m.Engine.HandleBackspace()
		case "ctrl+w":
			m.Engine.HandleCtrlBackspace()
		default:
			if len(msg.Runes) == 1 {
				key := msg.Runes[0]
				wasStarted := m.Engine.Started
				m.Engine.HandleKey(key)

				// Start timer on first keystroke
				if !wasStarted && m.Engine.Started {
					var cmds []tea.Cmd
					if m.Mode == "time" {
						cmds = append(cmds, tickCmd())
					}
					cmds = append(cmds, wpmSampleCmd())
					return m, tea.Batch(cmds...)
				}

				// Check if finished
				if m.Engine.Finished {
					return m, m.finishCmd()
				}

				// Check if failed (expert/master)
				if m.Engine.Failed {
					return m, m.finishCmd()
				}
			}
		}
	}

	return m, nil
}

func (m Model) finishCmd() tea.Cmd {
	return func() tea.Msg {
		return TestFinishedMsg{
			Engine: m.Engine,
			Mode:   m.Mode,
			Config: m.TCfg,
		}
	}
}
