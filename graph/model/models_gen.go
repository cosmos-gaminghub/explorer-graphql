// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AccountDetail struct {
	IsValidator     bool   `json:"is_validator"`
	OperatorAddress string `json:"operator_address"`
}

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Balances struct {
	Balances []*Balance `json:"balances"`
}

type Block struct {
	Height          int    `json:"height"`
	Hash            string `json:"hash"`
	ProposerAddr    string `json:"proposer_addr"`
	NumTxs          int    `json:"num_txs"`
	Time            string `json:"time"`
	Moniker         string `json:"moniker"`
	OperatorAddress string `json:"operator_address"`
	TotalRecords    int    `json:"total_records"`
}

type Change struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Subspace string `json:"subspace"`
}

type Commission struct {
	Commission *Commissions `json:"commission"`
}

type CommissionInfo struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Commissions struct {
	Commission []*CommissionInfo `json:"commission"`
}

type Content struct {
	Title       string    `json:"title"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Amount      []*Amount `json:"amount"`
	Changes     []*Change `json:"changes"`
	Plan        *Plan     `json:"plan"`
}

type Delegation struct {
	Moniker          string `json:"moniker"`
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Amount           int    `json:"amount"`
}

type Deposit struct {
	Depositor string  `json:"depositor"`
	Amount    *string `json:"amount"`
	TxHash    string  `json:"tx_hash"`
	Time      string  `json:"time"`
}

type Entries struct {
	RedelegationEntry *RedelegationEntry `json:"redelegation_entry"`
	Balance           string             `json:"balance"`
}

type Entry struct {
	CreationHeight *string `json:"creation_height"`
	CompletionTime *string `json:"completion_time"`
	InitialBalance *string `json:"initial_balance"`
	Balance        *string `json:"balance"`
}

type Inflation struct {
	Inflation string `json:"inflation"`
}

type Plan struct {
	Name                string `json:"name"`
	Time                string `json:"time"`
	Height              string `json:"height"`
	Info                string `json:"info"`
	UpgradedClientState string `json:"upgraded_client_state"`
}

type PowerEvent struct {
	Height       int    `json:"height"`
	TxHash       string `json:"tx_hash"`
	Timestamp    string `json:"timestamp"`
	Amount       int    `json:"amount"`
	Type         string `json:"type"`
	TotalRecords int    `json:"total_records"`
}

type Price struct {
	Price            string `json:"price"`
	Volume24h        string `json:"volume_24h"`
	MarketCap        string `json:"market_cap"`
	PercentChange24h string `json:"percent_change_24h"`
}

type Proposal struct {
	ID             int       `json:"id"`
	Status         string    `json:"status"`
	VotingStart    string    `json:"voting_start"`
	VotingEnd      string    `json:"voting_end"`
	SubmitTime     string    `json:"submit_time"`
	Tally          *Tally    `json:"tally"`
	Content        *Content  `json:"content"`
	Proposer       string    `json:"proposer"`
	Moniker        string    `json:"moniker"`
	TotalDeposit   []*Amount `json:"total_deposit"`
	DepositEndTime string    `json:"deposit_end_time"`
}

type Redelegation struct {
	DelegatorAddress    string               `json:"delegator_address"`
	ValidatorDstAddress string               `json:"validator_dst_address"`
	ValidatorSrcAddress string               `json:"validator_src_address"`
	MonikerSrc          string               `json:"moniker_src"`
	MonikerDst          string               `json:"moniker_dst"`
	Entries             []*RedelegationEntry `json:"entries"`
}

type RedelegationEntry struct {
	CreationHeight int    `json:"creation_height"`
	CompletionTime string `json:"completion_time"`
	InitialBalance string `json:"initial_balance"`
	SharesDst      string `json:"shares_dst"`
}

type RedelegationResponse struct {
	Redelegation *Redelegation `json:"redelegation"`
	Entries      []*Entries    `json:"entries"`
}

type Redelegations struct {
	RedelegationResponses []*RedelegationResponse `json:"redelegation_responses"`
}

type Reward struct {
	ValidatorAddress string        `json:"validator_address"`
	Reward           []*RewardInfo `json:"reward"`
}

type RewardInfo struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type Rewards struct {
	Rewards []*Reward `json:"rewards"`
}

type StatsAsset struct {
	Price     string `json:"price"`
	MarketCap string `json:"market_cap"`
	Volume24h string `json:"volume_24h"`
	Timestamp string `json:"timestamp"`
}

type Status struct {
	BlockHeight       int                `json:"block_height"`
	BlockTime         int                `json:"block_time"`
	TotalTxsNum       int                `json:"total_txs_num"`
	BondedTokens      int                `json:"bonded_tokens"`
	TotalSupplyTokens *TotalSupplyTokens `json:"total_supply_tokens"`
	Timestamp         string             `json:"timestamp"`
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
	RawLog    string `json:"raw_log"`
}

type Unbonding struct {
	UnbondingResponses []*UnbondingResponse `json:"unbonding_responses"`
}

type UnbondingResponse struct {
	DelegatorAddress string   `json:"delegator_address"`
	ValidatorAddress string   `json:"validator_address"`
	Moniker          string   `json:"moniker"`
	Entries          []*Entry `json:"entries"`
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
	Moniker          string  `json:"moniker"`
	VotingPower      int     `json:"voting_power"`
	CumulativeShare  string  `json:"cumulative_share"`
	Uptime           int     `json:"uptime"`
	OverBlocks       int     `json:"over_blocks"`
	Commission       float64 `json:"commission"`
	OperatorAddress  string  `json:"operator_address"`
	AccAddress       string  `json:"acc_address"`
	Jailed           bool    `json:"jailed"`
	Status           string  `json:"status"`
	Website          string  `json:"website"`
	Rank             int     `json:"rank"`
	Details          string  `json:"details"`
	Identity         string  `json:"identity"`
	ImageURL         string  `json:"image_url"`
	TotalMissedBlock int     `json:"total_missed_block"`
}

type Vote struct {
	Voter   string `json:"voter"`
	Option  string `json:"option"`
	TxHash  string `json:"tx_hash"`
	Time    string `json:"time"`
	Moniker string `json:"moniker"`
}
