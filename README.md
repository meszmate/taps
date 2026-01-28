# Taps

A terminal typing test built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

![Go](https://img.shields.io/github/go-mod/go-version/meszmate/taps)
![License](https://img.shields.io/github/license/meszmate/taps)
![Release](https://img.shields.io/github/v/release/meszmate/taps)

## Features

- **Test modes** — time (15/30/60/120s), word count (10/25/50/100), quote, and zen (freeform)
- **Live feedback** — per-character coloring (correct, incorrect, extra, missed), live WPM and accuracy
- **Results screen** — net/raw WPM, accuracy, consistency, character breakdown, WPM-over-time graph
- **10 built-in themes** — Default Dark, Dracula, Nord, Gruvbox, Catppuccin Mocha, Solarized Dark, Tokyo Night, One Dark, Rose Pine, Serika Dark
- **History tracking** — every completed test saved locally with personal bests and averages
- **Configurable** — punctuation, numbers, difficulty (normal/expert/master), cursor style, tape mode, focus mode, and more

## Install

### From release

Download the latest binary from [Releases](https://github.com/meszmate/taps/releases).

### From source

```bash
go install github.com/meszmate/taps/cmd/taps@latest
```

### Build locally

```bash
git clone https://github.com/meszmate/taps.git
cd taps
make build
./bin/taps
```

## Usage

Run `taps` to open the main menu.

### Menu controls

| Key | Action |
|-----|--------|
| `1-4` | Select mode (time/words/quote/zen) |
| `arrows` | Navigate menu / change duration or word count |
| `p` | Toggle punctuation |
| `n` | Toggle numbers |
| `enter` | Confirm selection |

### During a test

| Key | Action |
|-----|--------|
| `tab` | Restart test |
| `esc` | Back to menu |
| `ctrl+w` | Delete current word |
| `ctrl+c` | Quit |

### Results screen

| Key | Action |
|-----|--------|
| `tab` | Restart same test |
| `enter` | New test |
| `esc` | Back to menu |

## Configuration

Settings are persisted to `~/.config/taps/config.json`. All options can be changed from the in-app settings screen.

| Setting | Options |
|---------|---------|
| Mode | time, words, quote, zen |
| Duration | 15, 30, 60, 120 seconds |
| Word count | 10, 25, 50, 100 |
| Language | english (200 words), english_1k (1000 words) |
| Punctuation | on/off |
| Numbers | on/off |
| Difficulty | normal, expert (fail on wrong word), master (fail on wrong char) |
| Theme | 10 built-in themes |
| Cursor style | line, block, underline |
| Live WPM | on/off |
| Live accuracy | on/off |
| Stop on error | off, word, letter |
| Freedom mode | on/off (backspace to previous words) |
| Tape mode | on/off (single-line horizontal scroll) |
| Focus mode | on/off (minimal UI during test) |

## Themes

| Theme | Colors |
|-------|--------|
| Default Dark | Blue accent on dark |
| Dracula | Purple/pink on dark |
| Nord | Arctic blue palette |
| Gruvbox Dark | Warm retro |
| Catppuccin Mocha | Pastel on dark |
| Solarized Dark | Teal/orange on dark blue |
| Tokyo Night | Blue/purple on dark |
| One Dark | Atom editor style |
| Rose Pine | Muted rose/gold |
| Serika Dark | Monkeytype default |

Custom themes can be defined in `config.json` by providing 8 hex colors (background, foreground, sub, main, caret, correct, error, extra_error).

## Data

Test history is stored at `~/.local/share/taps/history.json`.

## License

[MIT](LICENSE)
