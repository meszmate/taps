package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// RainbowText renders each character with a gradient from start to end color
func RainbowText(text string, startHex, endHex string) string {
	if len(text) == 0 {
		return ""
	}

	startColor, _ := colorful.Hex(startHex)
	endColor, _ := colorful.Hex(endHex)

	runes := []rune(text)
	result := ""
	for i, r := range runes {
		t := float64(i) / float64(max(len(runes)-1, 1))
		c := startColor.BlendHsv(endColor, t)
		hex := c.Hex()
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(hex))
		result += style.Render(string(r))
	}
	return result
}

// FullRainbow renders text with a full spectrum rainbow
func FullRainbow(text string) string {
	if len(text) == 0 {
		return ""
	}

	runes := []rune(text)
	result := ""
	colors := []string{"#ff0000", "#ff7700", "#ffff00", "#00ff00", "#0077ff", "#8b00ff"}

	for i, r := range runes {
		t := float64(i) / float64(max(len(runes)-1, 1))
		// Interpolate across the color stops
		segmentFloat := t * float64(len(colors)-1)
		segIdx := int(segmentFloat)
		if segIdx >= len(colors)-1 {
			segIdx = len(colors) - 2
		}
		segT := segmentFloat - float64(segIdx)

		c1, _ := colorful.Hex(colors[segIdx])
		c2, _ := colorful.Hex(colors[segIdx+1])
		c := c1.BlendHsv(c2, segT)

		style := lipgloss.NewStyle().Foreground(lipgloss.Color(c.Hex()))
		result += style.Render(string(r))
	}
	return result
}
