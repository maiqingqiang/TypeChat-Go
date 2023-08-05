package main

import "fmt"

type Sentiment int

const (
	Negative Sentiment = iota
	Neutral
	Positive
)

func (s Sentiment) String() string {
	switch s {
	case Negative:
		return "negative"
	case Neutral:
		return "neutral"
	case Positive:
		return "positive"
	default:
		return ""
	}
}

type SentimentResponse struct {
	Sentiment int `json:"sentiment"`
}

func (s SentimentResponse) String() string {
	return fmt.Sprintf("%s", s.Sentiment)
}
