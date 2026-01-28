package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

type TestResult struct {
	Date        time.Time `json:"date"`
	Mode        string    `json:"mode"`
	Duration    int       `json:"duration"`
	WordCount   int       `json:"word_count"`
	Language    string    `json:"language"`
	Punctuation bool      `json:"punctuation"`
	Numbers     bool      `json:"numbers"`
	Difficulty  string    `json:"difficulty"`
	NetWPM      float64   `json:"net_wpm"`
	RawWPM      float64   `json:"raw_wpm"`
	Accuracy    float64   `json:"accuracy"`
	Consistency float64   `json:"consistency"`
	Correct     int       `json:"correct"`
	Incorrect   int       `json:"incorrect"`
	Extra       int       `json:"extra"`
	Missed      int       `json:"missed"`
	QuoteLength string    `json:"quote_length,omitempty"`
}

func historyPath() (string, error) {
	return xdg.DataFile("taps/history.json")
}

func Load() ([]TestResult, error) {
	p, err := historyPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var results []TestResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func Save(results []TestResult) error {
	p, err := historyPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

func Append(result TestResult) error {
	results, err := Load()
	if err != nil {
		results = nil
	}
	results = append(results, result)
	return Save(results)
}
