// The following is a schema definition for determining the sentiment of a some user input.

package main

import "fmt"

// Sentiment Define the enum type, this value is int
type Sentiment int

// Define the enum constants for Sentiment
const (
	Negative Sentiment = iota
	Neutral
	Positive
)

// Use switch statement to handle the enum type for Sentiment
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
	Sentiment Sentiment `json:"sentiment"` // The sentiment of the Sentiment enum type
}

func (s SentimentResponse) String() string {
	return fmt.Sprintf("%s", s.Sentiment)
}
