package history

import "sort"

type Stats struct {
	TotalTests   int
	AverageWPM   float64
	BestWPM      float64
	TotalWords   int
	Last10Avg    float64
	PersonalBest *TestResult
}

func CalculateStats(results []TestResult) Stats {
	if len(results) == 0 {
		return Stats{}
	}

	s := Stats{
		TotalTests: len(results),
	}

	totalWPM := 0.0
	for i := range results {
		r := &results[i]
		totalWPM += r.NetWPM
		if r.NetWPM > s.BestWPM {
			s.BestWPM = r.NetWPM
			s.PersonalBest = r
		}
		s.TotalWords += r.Correct / 5 // approximate words
	}
	s.AverageWPM = totalWPM / float64(len(results))

	// Last 10 average
	sorted := make([]TestResult, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})

	count := 10
	if count > len(sorted) {
		count = len(sorted)
	}
	last10Total := 0.0
	for i := 0; i < count; i++ {
		last10Total += sorted[i].NetWPM
	}
	s.Last10Avg = last10Total / float64(count)

	return s
}

func PersonalBestForConfig(results []TestResult, mode, language string, duration, wordCount int) *TestResult {
	var best *TestResult
	for i := range results {
		r := &results[i]
		if r.Mode != mode || r.Language != language {
			continue
		}
		if mode == "time" && r.Duration != duration {
			continue
		}
		if mode == "words" && r.WordCount != wordCount {
			continue
		}
		if best == nil || r.NetWPM > best.NetWPM {
			best = r
		}
	}
	return best
}
