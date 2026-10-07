package user

type Personality struct {
	Patience      float64 `json:"patience"`
	Curiosity     float64 `json:"curiosity"`
	AttentionSpan float64 `json:"attention_span"`
	Exploration   float64 `json:"exploration"`
	BingeWatching float64 `json:"binge_watching"`
	Consistency   float64 `json:"consistency"`
	Impulsiveness float64 `json:"impulsiveness"`
}
