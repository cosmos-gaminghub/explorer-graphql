package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/cosmos-gaminghub/exploder-graphql/conf"
	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"github.com/cosmos-gaminghub/exploder-graphql/utils"
)

const (
	DelegationUrl          = "%s/cosmos/staking/v1beta1/delegations/%s"
	CommissionUrl          = "%s/cosmos/distribution/v1beta1/validators/%s/commission"       // for commission
	RewardUrl              = "%s/cosmos/distribution/v1beta1/delegators/%s/rewards"          // for Reward
	UnbondingDelegationUrl = "%s/cosmos/staking/v1beta1/delegators/%s/unbonding_delegations" // for unbonding
	BalanceUrl             = "%s/cosmos/bank/v1beta1/balances/%s"                            // for avaiable
	SupplyUrl              = "%s/cosmos/bank/v1beta1/supply"                                 // for supply tokens
	InflationUrl           = "%s/cosmos/mint/v1beta1/inflation"                              // for inflation
	RedelegationUrl        = "%s/cosmos/staking/v1beta1/delegators/%s/redelegations"         // for Redelegation
)

// GetDelegation from lcd api
func GetDelegation(accAddress string) (DelegationResult, error) {
	url := fmt.Sprintf(DelegationUrl, conf.Get().LcdUrl, accAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		return DelegationResult{}, errors.New("Get delegation error")
	}

	var result DelegationResult
	if err := json.Unmarshal(resBytes, &result); err != nil {
		log.Fatalln("Unmarshal delegation error")
		return result, errors.New("Unmarshal delegation error")
	}

	return result, nil
}

// GetDelegation from lcd api
func GetSupply() (*model.TotalSupplyTokens, error) {
	url := fmt.Sprintf(SupplyUrl, conf.Get().LcdUrl)
	resBytes, err := utils.Get(url)
	if err != nil {
		return &model.TotalSupplyTokens{}, errors.New("Get supply error")
	}

	var result *model.TotalSupplyTokens
	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal supply error")
	}

	return result, nil
}

// GetDelegation from lcd api
func GetInflation() (*model.Inflation, error) {
	url := fmt.Sprintf(InflationUrl, conf.Get().LcdUrl)
	resBytes, err := utils.Get(url)
	if err != nil {
		return &model.Inflation{}, errors.New("Get inflation error")
	}

	var result *model.Inflation
	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal inflation error")
	}

	return result, nil
}

func GetBalances(accAddress string) (result *model.Balances, err error) {
	url := fmt.Sprintf(BalanceUrl, conf.Get().LcdUrl, accAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, errors.New("Get balances error")
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal balances error")
	}

	return result, nil
}

func GetRewards(accAddress string) (result *model.Rewards, err error) {
	url := fmt.Sprintf(RewardUrl, conf.Get().LcdUrl, accAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, errors.New("Get rewards error")
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal rewards error")
	}

	return result, nil
}

func GetCommission(operatorAddress string) (result *model.Commission, err error) {
	url := fmt.Sprintf(CommissionUrl, conf.Get().LcdUrl, operatorAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, errors.New("Get commision error")
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal commision error")
	}

	return result, nil
}

func GetUnbonding(accAddress string) (result *model.Unbonding, err error) {
	url := fmt.Sprintf(UnbondingDelegationUrl, conf.Get().LcdUrl, accAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, errors.New("Get unbonding error")
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal unbonding error")
	}

	return result, nil
}

func GetRedelegation(accAddress string) (result *model.Redelegations, err error) {
	url := fmt.Sprintf(RedelegationUrl, conf.Get().LcdUrl, accAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, errors.New("Get redelegations error")
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		return result, errors.New("Unmarshal redelegations error")
	}

	return result, nil
}
