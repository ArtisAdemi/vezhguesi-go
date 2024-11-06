package helper

func GetSentimentLabel(score float32) string {
	switch {
	case score > 0:
		return "positive"
	case score < 0:
		return "negative"
	default:
		return "neutral"
	}
}
