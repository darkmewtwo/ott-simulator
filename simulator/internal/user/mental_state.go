package user

type MentalState struct {
	Energy       float64 `json:"energy"`
	Happiness    float64 `json:"happiness"`
	Boredom      float64 `json:"boredom"`
	Frustration  float64 `json:"frustration"`
	Immersion    float64 `json:"immersion"`
	Satisfaction float64 `json:"satisfaction"`
}
