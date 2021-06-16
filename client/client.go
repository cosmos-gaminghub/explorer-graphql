package client

import (
	"encoding/json"
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
)

// GetDelegation from lcd api
func GetDelegation(accAddress string) (DelegationResult, error) {
	url := fmt.Sprintf(DelegationUrl, conf.Get().LcdUrl, accAddress)
	resBytes, err := utils.Get(url)
	if err != nil {
		log.Fatalln("Get delegation error")
		return DelegationResult{}, err
	}

	var result DelegationResult
	if err := json.Unmarshal(resBytes, &result); err != nil {
		log.Fatalln("Unmarshal delegation error")
		return result, err
	}

	return result, nil
}

// GetDelegation from lcd api
func GetSupply() (*model.TotalSupplyTokens, error) {
	url := fmt.Sprintf(SupplyUrl, conf.Get().LcdUrl)
	resBytes, err := utils.Get(url)
	if err != nil {
		log.Fatalln("Get supply error")
		return &model.TotalSupplyTokens{}, err
	}

	var result *model.TotalSupplyTokens
	if err := json.Unmarshal(resBytes, &result); err != nil {
		log.Fatalln("Unmarshal supply error")
		return result, err
	}

	return result, nil
}

// GetDelegation from lcd api
func GetInflation() (*model.Inflation, error) {
	url := fmt.Sprintf(InflationUrl, conf.Get().LcdUrl)
	resBytes, err := utils.Get(url)
	if err != nil {
		log.Fatalln("Get inflation error")
		return &model.Inflation{}, err
	}

	var result *model.Inflation
	if err := json.Unmarshal(resBytes, &result); err != nil {
		log.Fatalln("Unmarshal inflation error")
		return result, err
	}

	return result, nil
}
