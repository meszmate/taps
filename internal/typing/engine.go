package typing

import (
	"strings"
	"time"
)

type CharState int

const (
	CharUntyped CharState = iota
	CharCorrect
	CharIncorrect
	CharExtra
	CharMissed
)

type DisplayChar struct {
	Expected rune
	Typed    rune
	State    CharState
}

type Engine struct {
	Target         string
	Words          []string
	Chars          []DisplayChar
	CursorPos      int
	CurrentWord    int
	StartTime      time.Time
	Started        bool
	Finished       bool
	PerSecondWPM   []float64
	TotalTyped     int
	CorrectChars   int
	IncorrectChars int
	ExtraChars     int
	MissedChars    int
	ExtraByWord    map[int][]DisplayChar // extra chars per word index
	lastSampleTime time.Time
	wordStartIdx   []int // start index of each word in Chars
	wordEndIdx     []int // end index (exclusive) of each word in Chars

	// Config options
	StopOnError string // "off", "word", "letter"
	FreedomMode bool
	Difficulty  string // "normal", "expert", "master"

	// Track failed state for expert/master
	Failed       bool
	FailedReason string
}

func NewEngine(target string, stopOnError string, freedomMode bool, difficulty string) *Engine {
	runes := []rune(target)
	chars := make([]DisplayChar, len(runes))
	for i, r := range runes {
		chars[i] = DisplayChar{Expected: r, State: CharUntyped}
	}

	words := strings.Split(target, " ")

	// Calculate word boundaries
	wordStartIdx := make([]int, len(words))
	wordEndIdx := make([]int, len(words))
	pos := 0
	for i, w := range words {
		wordStartIdx[i] = pos
		wordEndIdx[i] = pos + len([]rune(w))
		pos += len([]rune(w)) + 1 // +1 for space
	}

	return &Engine{
		Target:      target,
		Words:       words,
		Chars:       chars,
		ExtraByWord: make(map[int][]DisplayChar),
		StopOnError: stopOnError,
		FreedomMode: freedomMode,
		Difficulty:  difficulty,
		wordStartIdx: wordStartIdx,
		wordEndIdx:   wordEndIdx,
	}
}

func (e *Engine) HandleKey(key rune) {
	if e.Finished || e.Failed {
		return
	}
	if !e.Started {
		e.Started = true
		e.StartTime = time.Now()
		e.lastSampleTime = e.StartTime
	}

	e.TotalTyped++

	if e.CursorPos >= len(e.Chars) {
		// We've gone past all characters - add as extra to last word
		e.ExtraChars++
		e.ExtraByWord[e.CurrentWord] = append(e.ExtraByWord[e.CurrentWord], DisplayChar{
			Typed: key,
			State: CharExtra,
		})
		return
	}

	expected := e.Chars[e.CursorPos].Expected

	if expected == ' ' {
		// Space pressed - move to next word
		if key == ' ' {
			// Mark any remaining chars in current word as missed
			if e.CurrentWord < len(e.Words) {
				end := e.wordEndIdx[e.CurrentWord]
				for i := e.CursorPos; i < end && i < len(e.Chars); i++ {
					if e.Chars[i].State == CharUntyped {
						e.Chars[i].State = CharMissed
						e.MissedChars++
					}
				}
			}
			e.Chars[e.CursorPos].State = CharCorrect
			e.Chars[e.CursorPos].Typed = key
			e.CorrectChars++
			e.CursorPos++
			e.CurrentWord++
		} else {
			// Typed a non-space where space expected - add as extra char
			e.ExtraChars++
			e.ExtraByWord[e.CurrentWord] = append(e.ExtraByWord[e.CurrentWord], DisplayChar{
				Typed: key,
				State: CharExtra,
			})
			if e.Difficulty == "master" {
				e.Failed = true
				e.FailedReason = "Wrong character (Master mode)"
			}
		}
	} else if key == ' ' {
		// Space pressed but not expected - skip to next word
		// Mark remaining chars as missed
		if e.CurrentWord < len(e.Words) {
			end := e.wordEndIdx[e.CurrentWord]
			for i := e.CursorPos; i < end && i < len(e.Chars); i++ {
				if e.Chars[i].State == CharUntyped {
					e.Chars[i].State = CharMissed
					e.MissedChars++
				}
			}
			// Mark space char
			spaceIdx := end
			if spaceIdx < len(e.Chars) {
				e.Chars[spaceIdx].State = CharCorrect
				e.Chars[spaceIdx].Typed = ' '
				e.CorrectChars++
				e.CursorPos = spaceIdx + 1
			}
		}
		e.CurrentWord++
		if e.Difficulty == "expert" {
			// Check if any chars in the word we just left were incorrect
			if e.CurrentWord > 0 {
				wordIdx := e.CurrentWord - 1
				if wordIdx < len(e.Words) {
					start := e.wordStartIdx[wordIdx]
					end := e.wordEndIdx[wordIdx]
					for i := start; i < end; i++ {
						if e.Chars[i].State == CharIncorrect || e.Chars[i].State == CharMissed {
							e.Failed = true
							e.FailedReason = "Wrong word (Expert mode)"
							return
						}
					}
				}
			}
		}
	} else if key == expected {
		e.Chars[e.CursorPos].State = CharCorrect
		e.Chars[e.CursorPos].Typed = key
		e.CorrectChars++
		e.CursorPos++
	} else {
		// Wrong character
		if e.StopOnError == "letter" {
			// Don't advance cursor
			return
		}
		e.Chars[e.CursorPos].State = CharIncorrect
		e.Chars[e.CursorPos].Typed = key
		e.IncorrectChars++
		if e.StopOnError != "word" {
			e.CursorPos++
		}
		if e.Difficulty == "master" {
			e.Failed = true
			e.FailedReason = "Wrong character (Master mode)"
		}
	}

	// Check if test is finished (word mode / quote mode)
	if e.CursorPos >= len(e.Chars) {
		e.Finished = true
	}
}

func (e *Engine) HandleBackspace() {
	if e.Finished || e.Failed || !e.Started {
		return
	}

	// Check for extra chars in current word first
	if extras, ok := e.ExtraByWord[e.CurrentWord]; ok && len(extras) > 0 {
		e.ExtraByWord[e.CurrentWord] = extras[:len(extras)-1]
		e.ExtraChars--
		e.TotalTyped--
		return
	}

	if e.CursorPos <= 0 {
		return
	}

	// Check if we're at the start of a word
	if e.CurrentWord < len(e.wordStartIdx) && e.CursorPos == e.wordStartIdx[e.CurrentWord] {
		if !e.FreedomMode {
			return // Can't go back to previous word without freedom mode
		}
		// Go back to previous word
		if e.CurrentWord > 0 {
			e.CurrentWord--
			// Go back past the space
			e.CursorPos--
			if e.Chars[e.CursorPos].State == CharCorrect {
				e.CorrectChars--
			}
			e.Chars[e.CursorPos].State = CharUntyped
			e.Chars[e.CursorPos].Typed = 0
		}
		return
	}

	e.CursorPos--
	ch := &e.Chars[e.CursorPos]
	switch ch.State {
	case CharCorrect:
		e.CorrectChars--
	case CharIncorrect:
		e.IncorrectChars--
	case CharMissed:
		e.MissedChars--
	}
	ch.State = CharUntyped
	ch.Typed = 0
}

func (e *Engine) HandleCtrlBackspace() {
	if e.Finished || e.Failed || !e.Started {
		return
	}

	// Delete entire current word progress
	if e.CurrentWord < len(e.wordStartIdx) {
		start := e.wordStartIdx[e.CurrentWord]
		// Remove extras first
		if extras, ok := e.ExtraByWord[e.CurrentWord]; ok {
			e.ExtraChars -= len(extras)
			delete(e.ExtraByWord, e.CurrentWord)
		}
		for i := e.CursorPos - 1; i >= start; i-- {
			ch := &e.Chars[i]
			switch ch.State {
			case CharCorrect:
				e.CorrectChars--
			case CharIncorrect:
				e.IncorrectChars--
			}
			ch.State = CharUntyped
			ch.Typed = 0
		}
		e.CursorPos = start
	}
}

func (e *Engine) SampleWPM() {
	if !e.Started || e.Finished {
		return
	}
	elapsed := time.Since(e.StartTime).Seconds()
	if elapsed <= 0 {
		return
	}
	raw := RawWPM(e.TotalTyped, elapsed)
	e.PerSecondWPM = append(e.PerSecondWPM, raw)
	e.lastSampleTime = time.Now()
}

func (e *Engine) ElapsedSeconds() float64 {
	if !e.Started {
		return 0
	}
	if e.Finished {
		// Use last sample time as approximate end
		return time.Since(e.StartTime).Seconds()
	}
	return time.Since(e.StartTime).Seconds()
}

func (e *Engine) CurrentRawWPM() float64 {
	return RawWPM(e.TotalTyped, e.ElapsedSeconds())
}

func (e *Engine) CurrentNetWPM() float64 {
	return NetWPM(e.CorrectChars, e.ElapsedSeconds())
}

func (e *Engine) CurrentAccuracy() float64 {
	return Accuracy(e.CorrectChars, e.IncorrectChars, e.ExtraChars)
}

func (e *Engine) CurrentConsistency() float64 {
	return Consistency(e.PerSecondWPM)
}

func (e *Engine) Finish() {
	if e.Finished {
		return
	}
	e.Finished = true
	// Count remaining untyped as missed
	for i := e.CursorPos; i < len(e.Chars); i++ {
		if e.Chars[i].State == CharUntyped {
			e.Chars[i].State = CharMissed
			e.MissedChars++
		}
	}
}

func (e *Engine) Progress() float64 {
	if len(e.Chars) == 0 {
		return 0
	}
	return float64(e.CursorPos) / float64(len(e.Chars))
}

func (e *Engine) WordProgress() (typed, total int) {
	return e.CurrentWord, len(e.Words)
}
