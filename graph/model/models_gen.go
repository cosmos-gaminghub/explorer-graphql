// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Block struct {
	Height          int    `json:"height"`
	Hash            string `json:"hash"`
	ProposerAddress string `json:"proposer_address"`
	NumTxs          int    `json:"num_txs"`
	Time            string `json:"time"`
}

type Change struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Subspace string `json:"subspace"`
}

type Content struct {
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Changes     []*Change `json:"changes"`
}

type Delegation struct {
	Moniker          string `json:"moniker"`
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           int    `json:"amount"`
}

type Deposit struct {
	ProposalID string    `json:"proposal_id"`
	Depositor  string    `json:"depositor"`
	Amount     []*Amount `json:"amount"`
}

type PowerEvent struct {
	Height    int    `json:"height"`
	TxHash    string `json:"tx_hash"`
	Timestamp string `json:"timestamp"`
	Amount    int    `json:"amount"`
	Type      string `json:"type"`
}

type Proposal struct {
	ID          int        `json:"id"`
	Status      string     `json:"status"`
	VotingStart string     `json:"voting_start"`
	VotingEnd   string     `json:"voting_end"`
	SubmitTime  string     `json:"submit_time"`
	Deposit     []*Deposit `json:"deposit"`
	Vote        []*Vote    `json:"vote"`
	Tally       *Tally     `json:"tally"`
	Content     *Content   `json:"content"`
	Proposer    string     `json:"proposer"`
}

type Status struct {
	BlockHeight       int                `json:"block_height"`
	BlockTime         string             `json:"block_time"`
	TotalTxsNum       int                `json:"total_txs_num"`
	BondedTokens      int                `json:"bonded_tokens"`
	TotalSupplyTokens *TotalSupplyTokens `json:"total_supply_tokens"`
}

type Supply struct {
	Denom  *string `json:"denom"`
	Amount *string `json:"amount"`
}

type Tally struct {
	Yes        string `json:"yes"`
	Abstain    string `json:"abstain"`
	No         string `json:"no"`
	NoWithVeto string `json:"no_with_veto"`
}

type TotalSupplyTokens struct {
	Supply []*Supply `json:"supply"`
}

type Tx struct {
	TxHash    string `json:"tx_hash"`
	Status    int    `json:"status"`
	Fee       string `json:"fee"`
	Height    int    `json:"height"`
	Timestamp string `json:"timestamp"`
	Messages  string `json:"messages"`
	Logs      string `json:"logs"`
	Memo      string `json:"memo"`
	GasUsed   int    `json:"gas_used"`
	GasWanted int    `json:"gas_wanted"`
}

type Uptime struct {
	Height    int    `json:"height"`
	Timestamp string `json:"timestamp"`
}

type UptimeResult struct {
	LastHeight int       `json:"last_height"`
	Uptime     []*Uptime `json:"uptime"`
}

type Validator struct {
	Moniker         string  `json:"moniker"`
	VotingPower     int     `json:"voting_power"`
	CumulativeShare string  `json:"cumulative_share"`
	Uptime          int     `json:"uptime"`
	OverBlocks      int     `json:"over_blocks"`
	Commission      float64 `json:"commission"`
	OperatorAddress string  `json:"operator_address"`
	AccAddress      string  `json:"acc_address"`
	Jailed          bool    `json:"jailed"`
	Status          string  `json:"status"`
	Website         string  `json:"website"`
}

type Vote struct {
	ProposalID string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}
