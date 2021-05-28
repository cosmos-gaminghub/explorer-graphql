package client

// Validator defines the structure for validator information.
type CommissionResult struct {
	Commission struct {
		Commission []Commission `json:"commission"`
	} `json:"commission"`
}
type Commission struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
