package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	Mode         string `json:"mode"`
	Duration     int    `json:"duration"`
	WordCount    int    `json:"word_count"`
	Language     string `json:"language"`
	Punctuation  bool   `json:"punctuation"`
	Numbers      bool   `json:"numbers"`
	Difficulty   string `json:"difficulty"`
	Theme        string `json:"theme"`
	CursorStyle  string `json:"cursor_style"`
	LiveWPM      bool   `json:"live_wpm"`
	LiveAccuracy bool   `json:"live_accuracy"`
	StopOnError  string `json:"stop_on_error"`
	FreedomMode  bool   `json:"freedom_mode"`
	TapeMode     bool   `json:"tape_mode"`
	ShowAllLines bool   `json:"show_all_lines"`
	FocusMode    bool   `json:"focus_mode"`
	SoundOnError bool   `json:"sound_on_error"`
	QuoteLength  string `json:"quote_length"`
	CustomTheme  *CustomThemeConfig `json:"custom_theme,omitempty"`
}

type CustomThemeConfig struct {
	Name       string `json:"name"`
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Sub        string `json:"sub"`
	Main       string `json:"main"`
	Caret      string `json:"caret"`
	Correct    string `json:"correct"`
	Error      string `json:"error"`
	ExtraError string `json:"extra_error"`
}

func DefaultConfig() *Config {
	return &Config{
		Mode:         DefaultMode,
		Duration:     DefaultDuration,
		WordCount:    DefaultWordCount,
		Language:     DefaultLanguage,
		Punctuation:  false,
		Numbers:      false,
		Difficulty:   DefaultDifficulty,
		Theme:        DefaultTheme,
		CursorStyle:  DefaultCursorStyle,
		LiveWPM:      true,
		LiveAccuracy: true,
		StopOnError:  DefaultStopOnError,
		FreedomMode:  false,
		TapeMode:     false,
		ShowAllLines: false,
		FocusMode:    false,
		SoundOnError: false,
		QuoteLength:  DefaultQuoteLength,
	}
}

func configPath() (string, error) {
	return xdg.ConfigFile("taps/config.json")
}

func Load() *Config {
	cfg := DefaultConfig()
	p, err := configPath()
	if err != nil {
		return cfg
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, cfg)
	return cfg
}

func (c *Config) Save() error {
	p, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}
