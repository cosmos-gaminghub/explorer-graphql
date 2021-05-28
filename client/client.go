package client

import (
	"encoding/json"
	"fmt"
	"graphql-go-example/conf"
	"graphql-go-example/utils"
	"log"
)

const (
	DelegationUrl          = "%s/cosmos/staking/v1beta1/delegations/%s"
	CommissionUrl          = "%s/cosmos/distribution/v1beta1/validators/%s/commission"       // for commission
	RewardUrl              = "%s/cosmos/distribution/v1beta1/delegators/%s/rewards"          // for Reward
	UnbondingDelegationUrl = "%s/cosmos/staking/v1beta1/delegators/%s/unbonding_delegations" // for unbonding
	BalanceUrl             = "%s/cosmos/bank/v1beta1/balances/%s"                            // for avaiable
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
