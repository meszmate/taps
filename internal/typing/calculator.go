package typing

import "math"

// RawWPM calculates raw words per minute: (totalKeystrokes / 5) / minutes
func RawWPM(totalKeystrokes int, elapsedSeconds float64) float64 {
	if elapsedSeconds <= 0 {
		return 0
	}
	minutes := elapsedSeconds / 60.0
	return (float64(totalKeystrokes) / 5.0) / minutes
}

// NetWPM calculates net words per minute: (correctChars / 5) / minutes
func NetWPM(correctChars int, elapsedSeconds float64) float64 {
	if elapsedSeconds <= 0 {
		return 0
	}
	minutes := elapsedSeconds / 60.0
	wpm := (float64(correctChars) / 5.0) / minutes
	if wpm < 0 {
		return 0
	}
	return wpm
}

// Accuracy calculates accuracy percentage
func Accuracy(correct, incorrect, extra int) float64 {
	total := correct + incorrect + extra
	if total == 0 {
		return 100
	}
	return float64(correct) / float64(total) * 100
}

// Consistency calculates typing consistency as 100 - CV (coefficient of variation)
// of per-second WPM samples
func Consistency(perSecondWPM []float64) float64 {
	if len(perSecondWPM) < 2 {
		return 100
	}

	mean := 0.0
	for _, v := range perSecondWPM {
		mean += v
	}
	mean /= float64(len(perSecondWPM))

	if mean == 0 {
		return 0
	}

	variance := 0.0
	for _, v := range perSecondWPM {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(perSecondWPM))

	stddev := math.Sqrt(variance)
	cv := (stddev / mean) * 100

	consistency := 100 - cv
	if consistency < 0 {
		return 0
	}
	return consistency
}
