package document

import (
	"fmt"
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/client"
	"github.com/cosmos-gaminghub/exploder-graphql/orm"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

const (
	CollectionNmValidator = "validator"

	ValidatorFieldVotingPower      = "voting_power"
	ValidatorFieldJailed           = "jailed"
	ValidatorFieldStatus           = "status"
	ValidatorFieldOperatorAddress  = "operator_address"
	ValidatorFieldDescription      = "description"
	ValidatorFieldConsensusAddr    = "consensus_pubkey"
	ValidatorFieldProposerHashAddr = "proposer_addr"
	ValidatorFieldTokens           = "tokens"
	ValidatorFieldDelegatorShares  = "delegator_shares"
	ValidatorFieldIcon             = "icons"
	ValidatorFieldValidatorAddress = "account_address"
	ValidatorStatusValUnbonded     = 0
	ValidatorStatusValUnbonding    = 1
	ValidatorStatusValBonded       = 2

	DefaultValidatorLimit = 125

	Bonded   = "BOND_STATUS_BONDED"
	Unbonded = "BOND_STATUS_UNBONDED"
)

// func (v Validator) GetValidatorStatus() string {

// 	if v.Jailed == false && v.Status == types.Bonded {
// 		return "Active"
// 	}

// 	if v.Status != types.Bonded && v.Jailed == false {
// 		return "Candidate"
// 	}

// 	return "Jailed"

// }

// func (v Validator) IsCandidatorWithStatus() bool {

// 	if v.Status != types.Bonded && v.Jailed == false {
// 		return true
// 	}

// 	return false

// }

type (
	UptimeChangeVo struct {
		Address string
		Time    string
		Uptime  float64
	}

	ValVotingPowerChangeVo struct {
		Height  int64
		Address string
		Power   int64
		Time    time.Time
		Change  string
	}

	ValUpTimeVo struct {
		Time   string `bson:"_id,omitempty"`
		Uptime float64
	}

	CountVo struct {
		Id    bson.ObjectId `bson:"_id,omitempty"`
		Count float64
	}
)

type Validator struct {
	OperatorAddr    string      `bson:"operator_address"`
	ConsensusPubkey string      `bson:"consensus_pubkey"`
	ConsensusAddres string      `bson:"consensus_address"`
	AccountAddr     string      `bson:"account_address"`
	Jailed          bool        `bson:"jailed"`
	Status          string      `bson:"status"`
	Tokens          int64       `bson:"tokens" json:"tokens"`
	DelegatorShares string      `bson:"delegator_shares"`
	Description     Description `bson:"description" json:"description"`
	UnbondingHeight string      `bson:"unbonding_height"`
	UnbondingTime   time.Time   `bson:"unbonding_time"`
	Commission      Commission  `bson:"commission" json:"commission"`
	ProposerAddr    string      `bson:"proposer_addr"`
	Icons           string      `bson:"icons"`
}

func (v Validator) GetValidatorList() ([]Validator, error) {
	var validatorsDocArr []Validator
	var query = orm.NewQuery()
	defer query.Release()

	var selector = bson.M{"description.moniker": 1, "description.identity": 1, "total_missed_block": 1, "operator_address": 1, "tokens": 1, "commission": 1, "jailed": 1, "status": 1}

	err := queryAll(CollectionNmValidator, selector, nil, desc(ValidatorFieldTokens), 0, &validatorsDocArr)
	return validatorsDocArr, err
}

func (v Validator) GetValidatorByProposerAddr(addr string) (Validator, error) {

	var selector = bson.M{"description.moniker": 1, "operator_address": 1}
	err := queryOne(CollectionNmValidator, selector, bson.M{"proposer_addr": addr}, &v)

	return v, err
}

type Description struct {
	ImageUrl string `bson:"imageurl" json:"imageurl"`
	Moniker  string `bson:"moniker" json:"moniker"`
	Identity string `bson:"identity" json:"identity"`
	Website  string `bson:"website" json:"website"`
	Details  string `bson:"details" json:"details"`
}

func (d Description) String() string {
	return fmt.Sprintf(`Moniker  :%v  Identity :%v Website  :%v Details  :%v`, d.Moniker, d.Identity, d.Website, d.Details)
}

type Commission struct {
	CommissionRate struct {
		Rate          string `bson:"rate"`
		MaxRate       string `bson:"max_rate"`
		MaxChangeRate string `bson:"max_change_rate"`
	}
	UpdateTime string `bson:"update_time"`
}

func (v Validator) Name() string {
	return CollectionNmValidator
}

func (_ Validator) GetAllValidator() ([]Validator, error) {
	var validators []Validator
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(CollectionNmValidator).
		SetResult(&validators)

	err := query.Exec()

	return validators, err
}

func (v Validator) QueryValidatorMonikerOpAddr(addrArrAsVa []string) ([]Validator, error) {
	var validators []Validator
	var selector = bson.M{
		"description.moniker": 1,
		"operator_address":    1,
	}

	err := queryAll(CollectionNmValidator, selector, bson.M{"operator_address": bson.M{"$in": addrArrAsVa}}, "", 0, &validators)
	return validators, err
}

func (v Validator) QueryValidatorsMonikerOpAddrConsensusPubkey() ([]Validator, error) {
	var validators []Validator
	var selector = bson.M{
		"description.moniker": 1,
		"operator_address":    1,
		"consensus_pubkey":    1,
		"status":              1,
		"voting_power":        1,
	}

	condition := bson.M{
		ValidatorFieldStatus: ValidatorStatusValBonded,
	}

	err := queryAll(CollectionNmValidator, selector, condition, "", 0, &validators)
	return validators, err
}

func (v Validator) QueryValidatorMonikerOpAddrByHashAddr(hashAddr []string) ([]Validator, error) {
	var validators []Validator
	var selector = bson.M{"description.moniker": 1, "operator_address": 1, "proposer_addr": 1}

	err := queryAll(CollectionNmValidator, selector, bson.M{"proposer_addr": bson.M{"$in": hashAddr}}, "", 0, &validators)
	return validators, err
}

func GetValidatorByAddr(addr string) (Validator, error) {
	db := getDb()
	c := db.C(CollectionNmValidator)
	defer db.Session.Close()
	var validator Validator
	err := c.Find(bson.M{ValidatorFieldOperatorAddress: addr}).One(&validator)

	return validator, err
}

func (_ Validator) GetBondedValidators() ([]Validator, error) {
	var (
		validators []Validator
	)

	selector := bson.M{
		ValidatorFieldTokens: "1",
	}
	condition := bson.M{
		ValidatorFieldStatus: ValidatorStatusValBonded,
	}

	err := queryAll(CollectionNmValidator, selector, condition, "", 0, &validators)

	return validators, err
}

func (_ Validator) GetBondedValidatorsSharesTokens() ([]Validator, error) {
	var (
		validators []Validator
	)

	selector := bson.M{
		ValidatorFieldVotingPower:     "1",
		ValidatorFieldOperatorAddress: "1",
		ValidatorFieldDelegatorShares: "1",
		ValidatorFieldTokens:          "1",
		ValidatorFieldDescription:     "1",
	}
	condition := bson.M{
		ValidatorFieldStatus: ValidatorStatusValBonded,
	}

	err := queryAll(CollectionNmValidator, selector, condition, "", 0, &validators)

	return validators, err
}

func (_ Validator) QueryValidatorListByAddrList(addrs []string) ([]Validator, error) {
	validatorArr := []Validator{}

	valCondition := bson.M{
		ValidatorFieldOperatorAddress: bson.M{"$in": addrs},
	}

	err := queryAll(CollectionNmValidator, nil, valCondition, "", 0, &validatorArr)

	return validatorArr, err
}

func (_ Validator) QueryMonikerAndValidatorAddrByHashAddr(addr string) (Validator, error) {

	selector := bson.M{
		ValidatorFieldOperatorAddress: 1,
		ValidatorFieldDescription:     1,
		ValidatorFieldIcon:            1,
	}
	condition := bson.M{ValidatorFieldProposerHashAddr: addr}
	var val Validator
	err := queryOne(CollectionNmValidator, selector, condition, &val)

	return val, err
}

func (_ Validator) QueryValidatorByConsensusAddr(addr string) (Validator, error) {
	var query = orm.NewQuery()
	defer query.Release()

	var result Validator
	condition := bson.M{}
	condition[ValidatorFieldConsensusAddr] = addr

	query.SetCollection(CollectionNmValidator).
		SetResult(&result).
		SetCondition(condition).
		SetSize(1)
	err := query.Exec()

	return result, err
}

func (_ Validator) QueryValidatorDetailByOperatorAddr(opAddr string) (Validator, error) {

	validator := Validator{}

	valCondition := bson.M{
		ValidatorFieldOperatorAddress: opAddr,
	}

	err := queryOne(CollectionNmValidator, nil, valCondition, &validator)

	return validator, err
}

func (_ Validator) QueryValidatorDetailByAccAddr(accAddress string) (Validator, error) {

	validator := Validator{}

	valCondition := bson.M{
		ValidatorFieldValidatorAddress: accAddress,
	}

	err := queryOne(CollectionNmValidator, nil, valCondition, &validator)

	return validator, err
}

func (_ Validator) QueryValidatorDetailByListAccAddr(listAccAddress []string) ([]Validator, error) {

	var validators []Validator
	var selector = bson.M{
		"description.moniker":          1,
		ValidatorFieldValidatorAddress: 1,
	}

	err := queryAll(CollectionNmValidator, selector, bson.M{ValidatorFieldValidatorAddress: bson.M{"$in": listAccAddress}}, "", 0, &validators)
	return validators, err
}

func (_ Validator) QueryTotalActiveValidatorVotingPower() (int64, error) {

	validators := []Validator{}
	condition := bson.M{ValidatorFieldJailed: false, ValidatorFieldStatus: Bonded}
	var selector = bson.M{ValidatorFieldTokens: 1}

	err := queryAll(CollectionNmValidator, selector, condition, "", 0, &validators)

	if err != nil {
		return 0, err
	}

	totalVotingPower := int64(0)
	for _, v := range validators {
		totalVotingPower += v.Tokens
	}
	return totalVotingPower, nil
}

func (_ Validator) Batch(txs []txn.Op) error {
	return orm.Batch(txs)
}

func (_ Validator) GetListOperatorAdress(validators []Validator) []string {
	var list []string
	for _, validator := range validators {
		list = append(list, validator.OperatorAddr)
	}
	return list
}

func (_ Validator) GetListConsensusAddress(validators []Validator) []string {
	var list []string
	for _, validator := range validators {
		list = append(list, validator.ConsensusAddres)
	}
	return list
}

func (_ Validator) GetListMapOperatorAndMoniker(delegationResult client.DelegationResult) map[string]string {
	var listOperatorAddress []string
	for _, validator := range delegationResult.DelegationResponses {
		listOperatorAddress = append(listOperatorAddress, validator.Delegation.ValidatorAddress)
	}
	var mapList = make(map[string]string)
	validators, err := Validator{}.QueryValidatorMonikerOpAddr(listOperatorAddress)
	if err != nil {
		return mapList
	}

	for _, validator := range validators {
		mapList[validator.OperatorAddr] = ""
		for _, item := range delegationResult.DelegationResponses {
			if item.Delegation.ValidatorAddress == validator.OperatorAddr {
				mapList[validator.OperatorAddr] = validator.Description.Moniker
				break
			}
		}
	}
	return mapList
}

func (_ Validator) GetListMapAccAndMoniker(listAccAddress []string) map[string]string {
	validators, err := Validator{}.QueryValidatorDetailByListAccAddr(listAccAddress)
	var mapList = make(map[string]string)
	if err != nil {
		return mapList
	}

	for _, validator := range validators {
		mapList[validator.AccountAddr] = validator.Description.Moniker
	}
	return mapList
}

func (_ Validator) MapOperatorAndMoniker(validators []Validator) map[string]string {
	var mapList = make(map[string]string)
	for _, validator := range validators {
		mapList[validator.OperatorAddr] = validator.Description.Moniker
	}
	return mapList
}

func (_ Validator) GetSumBondedToken() (int64, error) {
	result := []bson.M{}
	var query = orm.NewQuery()
	defer query.Release()

	selector := bson.M{ValidatorFieldStatus: 1, ValidatorFieldTokens: 1}
	condition := bson.M{
		ValidatorFieldJailed: false,
		ValidatorFieldStatus: Bonded,
	}
	query.SetResult(&result).
		SetCollection(CollectionNmValidator).
		SetSelector(selector).
		PipeQuery(
			[]bson.M{
				{"$match": condition},
				{"$group": bson.M{
					"_id":   "",
					"total": bson.M{"$sum": "$tokens"},
				}},
			},
		)
	if len(result) == 0 {
		return 0, nil
	}
	return result[0]["total"].(int64), nil
}

func (_ Validator) FormatListValidator(validators []Validator) (result []Validator) {
	for _, item := range validators {
		if item.Jailed == false {
			result = append(result, item)
		}
	}

	for _, item := range validators {
		if item.Jailed == true {
			result = append(result, item)
		}
	}
	return result
}

func (_ Validator) GetIndexFromFormatListValidator(validators []Validator, operator_address string) (rank int) {
	for index, item := range validators {
		if item.OperatorAddr == operator_address {
			return index + 1
		}
	}
	return rank
}

func IsActiveValidator(validator Validator) bool {
	return !validator.Jailed && validator.Status == Bonded
}
