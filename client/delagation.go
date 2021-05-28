package client

// Validator defines the structure for validator information.
type DelegationResult struct {
	DelegationResponses []struct {
		Delegation Delegation `json:"delegation"`
		Balance    Balance    `json:"balance"`
	} `json:"delegation_responses"`
}

type Delegation struct {
	DelegatorAdrress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
}

type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}
