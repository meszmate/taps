package test

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/meszmate/taps/internal/typing"
)

func (m Model) View() string {
	t := m.Styles.Theme
	var b strings.Builder

	maxWidth := m.Width
	if maxWidth <= 0 {
		maxWidth = 80
	}
	textWidth := maxWidth - 4
	if textWidth < 40 {
		textWidth = 40
	}
	if textWidth > 100 {
		textWidth = 100
	}

	// Top stats bar
	if !m.Config.FocusMode || !m.Engine.Started {
		b.WriteString(m.renderTopBar())
		b.WriteString("\n\n")
	}

	// Render typed text
	if m.Config.TapeMode {
		b.WriteString(m.renderTapeMode(textWidth))
	} else if m.Config.ShowAllLines {
		b.WriteString(m.renderAllLines(textWidth))
	} else {
		b.WriteString(m.render3Lines(textWidth))
	}

	// Failed message
	if m.Engine.Failed {
		b.WriteString("\n\n")
		errStyle := lipgloss.NewStyle().Foreground(t.Error).Bold(true)
		b.WriteString(errStyle.Render(m.Engine.FailedReason))
	}

	// Bottom help
	if !m.Config.FocusMode {
		b.WriteString("\n\n")
		helpStyle := lipgloss.NewStyle().Foreground(t.Sub)
		b.WriteString(helpStyle.Render("tab restart | esc menu"))
	}

	content := b.String()
	if m.Width > 0 && m.Height > 0 {
		content = lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
	}

	return content
}

func (m Model) renderTopBar() string {
	t := m.Styles.Theme
	var parts []string

	// Timer or progress
	switch m.Mode {
	case "time":
		timerStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)
		parts = append(parts, timerStyle.Render(fmt.Sprintf("%ds", m.Timer)))
	case "words":
		typed, total := m.Engine.WordProgress()
		progressStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)
		parts = append(parts, progressStyle.Render(fmt.Sprintf("%d/%d", typed, total)))
	case "quote":
		pct := m.Engine.Progress() * 100
		progressStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)
		parts = append(parts, progressStyle.Render(fmt.Sprintf("%.0f%%", pct)))
	case "zen":
		zenStyle := lipgloss.NewStyle().Foreground(t.Main).Bold(true)
		parts = append(parts, zenStyle.Render("zen"))
	}

	// Live WPM
	if m.Config.LiveWPM && m.Engine.Started {
		wpmStyle := lipgloss.NewStyle().Foreground(t.Sub)
		wpm := m.Engine.CurrentNetWPM()
		parts = append(parts, wpmStyle.Render(fmt.Sprintf("%.0f wpm", wpm)))
	}

	// Live accuracy
	if m.Config.LiveAccuracy && m.Engine.Started {
		accStyle := lipgloss.NewStyle().Foreground(t.Sub)
		acc := m.Engine.CurrentAccuracy()
		parts = append(parts, accStyle.Render(fmt.Sprintf("%.0f%% acc", acc)))
	}

	return strings.Join(parts, "  ")
}

// renderChars renders chars from startIdx to endIdx with per-character coloring and cursor
func (m Model) renderChars(startIdx, endIdx int) string {
	t := m.Styles.Theme
	var b strings.Builder

	for i := startIdx; i < endIdx && i < len(m.Engine.Chars); i++ {
		ch := m.Engine.Chars[i]

		// Cursor
		if i == m.Engine.CursorPos {
			b.WriteString(m.renderCursor(ch))
			continue
		}

		var style lipgloss.Style
		switch ch.State {
		case typing.CharCorrect:
			style = lipgloss.NewStyle().Foreground(t.Correct)
		case typing.CharIncorrect:
			style = lipgloss.NewStyle().Foreground(t.Error)
		case typing.CharMissed:
			style = lipgloss.NewStyle().Foreground(t.Sub).Strikethrough(true)
		case typing.CharExtra:
			style = lipgloss.NewStyle().Foreground(t.ExtraError)
		default: // untyped
			style = lipgloss.NewStyle().Foreground(t.Sub)
		}

		r := ch.Expected
		if ch.State == typing.CharIncorrect && ch.Typed != 0 {
			r = ch.Typed
		}
		b.WriteString(style.Render(string(r)))

		// Render extra chars after the last char of a word
		if i+1 < len(m.Engine.Chars) && m.Engine.Chars[i+1].Expected == ' ' {
			wordIdx := m.findWordForCharIdx(i)
			if extras, ok := m.Engine.ExtraByWord[wordIdx]; ok {
				extraStyle := lipgloss.NewStyle().Foreground(t.ExtraError)
				for _, ex := range extras {
					b.WriteString(extraStyle.Render(string(ex.Typed)))
				}
			}
		}
	}

	return b.String()
}

func (m Model) findWordForCharIdx(charIdx int) int {
	pos := 0
	for i, w := range m.Engine.Words {
		wordEnd := pos + len([]rune(w))
		if charIdx >= pos && charIdx < wordEnd {
			return i
		}
		pos = wordEnd + 1 // +1 for space
	}
	return 0
}

func (m Model) renderCursor(ch typing.DisplayChar) string {
	t := m.Styles.Theme
	cursorChar := string(ch.Expected)

	switch m.Config.CursorStyle {
	case "block":
		return lipgloss.NewStyle().
			Background(t.Caret).
			Foreground(t.Background).
			Render(cursorChar)
	case "underline":
		return lipgloss.NewStyle().
			Foreground(t.Caret).
			Underline(true).
			Render(cursorChar)
	default: // "line"
		return lipgloss.NewStyle().
			Foreground(t.Caret).
			Render("|" + cursorChar)
	}
}

// wordWrapIndices splits chars into lines by word boundaries within maxWidth
func (m Model) wordWrapLines(maxWidth int) []struct{ start, end int } {
	var lines []struct{ start, end int }
	chars := m.Engine.Chars

	lineStart := 0
	lineWidth := 0

	for i := 0; i < len(chars); i++ {
		lineWidth++

		if lineWidth >= maxWidth {
			// Find last space to break at
			breakIdx := i
			for j := i; j > lineStart; j-- {
				if chars[j].Expected == ' ' {
					breakIdx = j
					break
				}
			}
			lines = append(lines, struct{ start, end int }{lineStart, breakIdx})
			lineStart = breakIdx
			if breakIdx < len(chars) && chars[breakIdx].Expected == ' ' {
				lineStart = breakIdx + 1
			}
			lineWidth = i - lineStart + 1
		}
	}

	if lineStart < len(chars) {
		lines = append(lines, struct{ start, end int }{lineStart, len(chars)})
	}

	return lines
}

func (m Model) renderAllLines(maxWidth int) string {
	lines := m.wordWrapLines(maxWidth)
	var b strings.Builder
	for i, line := range lines {
		b.WriteString(m.renderChars(line.start, line.end))
		if i < len(lines)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m Model) render3Lines(maxWidth int) string {
	lines := m.wordWrapLines(maxWidth)
	if len(lines) == 0 {
		return ""
	}

	// Find which line the cursor is on
	cursorLine := 0
	for i, line := range lines {
		if m.Engine.CursorPos >= line.start && m.Engine.CursorPos < line.end {
			cursorLine = i
			break
		}
		// If cursor is at the end
		if m.Engine.CursorPos >= line.end && i == len(lines)-1 {
			cursorLine = i
		}
	}

	// Show 3 lines centered on cursor line
	startLine := cursorLine - 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := startLine + 3
	if endLine > len(lines) {
		endLine = len(lines)
		startLine = endLine - 3
		if startLine < 0 {
			startLine = 0
		}
	}

	var b strings.Builder
	for i := startLine; i < endLine; i++ {
		b.WriteString(m.renderChars(lines[i].start, lines[i].end))
		if i < endLine-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m Model) renderTapeMode(maxWidth int) string {
	// Single line horizontal scroll centered on cursor
	cursorPos := m.Engine.CursorPos
	halfWidth := maxWidth / 2

	start := cursorPos - halfWidth
	if start < 0 {
		start = 0
	}
	end := start + maxWidth
	if end > len(m.Engine.Chars) {
		end = len(m.Engine.Chars)
		start = end - maxWidth
		if start < 0 {
			start = 0
		}
	}

	return m.renderChars(start, end)
}
