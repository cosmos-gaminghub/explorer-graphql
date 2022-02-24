package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cosmos-gaminghub/exploder-graphql/client"
	"github.com/cosmos-gaminghub/exploder-graphql/conf"
	"github.com/cosmos-gaminghub/exploder-graphql/graph/generated"
	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"github.com/cosmos-gaminghub/exploder-graphql/orm/document"
	"github.com/cosmos-gaminghub/exploder-graphql/utils"
)

func (r *queryResolver) Blocks(ctx context.Context, offset *int, size *int) ([]*model.Block, error) {
	results, err := document.Block{}.GetBlockListByOffsetAndSize(*offset, *size)
	if err != nil {
		return []*model.Block{}, nil
	}
	return document.Block{}.FormatBsonMForModel(results)
}

func (r *queryResolver) BlockDetail(ctx context.Context, height *int) (*model.Block, error) {
	result, err := document.Block{}.QueryBlockByHeight(int64(*height))
	if err != nil {
		return &model.Block{}, nil
	}
	return document.Block{}.FormatBsonMForModelBlockDetail(result)
}

func (r *queryResolver) BlockTxs(ctx context.Context, height *int) ([]*model.Tx, error) {
	txs, err := document.CommonTx{}.QueryTxByHeight(int64(*height))
	if err != nil {
		return []*model.Tx{}, nil
	}
	return document.CommonTx{}.FormatListTxsForModel(txs)
}

func (r *queryResolver) Txs(ctx context.Context, size *int) ([]*model.Tx, error) {
	txs, err := document.CommonTx{}.GetListTxBy(*size)
	if err != nil {
		return []*model.Tx{}, nil
	}
	return document.CommonTx{}.FormatListTxsForModel(txs)
}

func (r *queryResolver) TxDetail(ctx context.Context, txHash *string) (*model.Tx, error) {
	tx, err := document.CommonTx{}.QueryTxByHash(*txHash)

	if err != nil {
		fmt.Print(err.Error())
		return &model.Tx{}, nil
	}

	return document.CommonTx{}.FormatTxForModel(tx)
}

func (r *queryResolver) Validators(ctx context.Context) ([]*model.Validator, error) {
	validators, err := document.Validator{}.GetValidatorList()
	if err != nil {
		return []*model.Validator{}, nil
	}

	listConsensusAddress := document.Validator{}.GetListConsensusAddress(validators)
	upTimeCount, overBlocks := document.MissedBlock{}.GetMissedBlockCount(listConsensusAddress)

	validatorFormat := document.Validator{}.FormatListValidator(validators)
	var listValidator []*model.Validator
	for index, validator := range validatorFormat {
		uptime := overBlocks
		if document.IsActiveValidator(validator) {
			uptime = upTimeCount[validator.ConsensusAddres]
		}

		commision, _ := utils.ParseStringToFloat(validator.Commission.CommissionRate.Rate)
		t := &model.Validator{
			Moniker:         validator.Description.Moniker,
			OperatorAddress: validator.OperatorAddr,
			AccAddress:      validator.AccountAddr,
			VotingPower:     int(validator.Tokens),
			Commission:      commision,
			Jailed:          validator.Jailed,
			Status:          validator.Status,
			Uptime:          uptime,
			OverBlocks:      overBlocks,
			Website:         validator.Description.Website,
			Rank:            index + 1,
			Identity:        validator.Description.Identity,
		}
		listValidator = append(listValidator, t)
	}
	return listValidator, nil
}

func (r *queryResolver) ValidatorDetail(ctx context.Context, operatorAddress *string) (*model.Validator, error) {
	validator, err := document.Validator{}.QueryValidatorDetailByOperatorAddr(*operatorAddress)
	if err != nil {
		return &model.Validator{}, nil
	}
	upTimeCount, overBlocks := document.MissedBlock{}.GetMissedBlockCount([]string{validator.OperatorAddr})
	commision, _ := utils.ParseStringToFloat(validator.Commission.CommissionRate.Rate)

	var rank int
	validators, err := document.Validator{}.GetValidatorList()
	if err == nil {
		validatorFormat := document.Validator{}.FormatListValidator(validators)
		rank = document.Validator{}.GetIndexFromFormatListValidator(validatorFormat, validator.OperatorAddr)
	}

	uptime := overBlocks
	if document.IsActiveValidator(validator) {
		uptime = upTimeCount[validator.OperatorAddr]
	}
	return &model.Validator{
		Moniker:         validator.Description.Moniker,
		OperatorAddress: validator.OperatorAddr,
		AccAddress:      validator.AccountAddr,
		VotingPower:     int(validator.Tokens),
		Commission:      commision,
		Jailed:          validator.Jailed,
		Status:          validator.Status,
		Uptime:          uptime,
		OverBlocks:      overBlocks,
		Website:         validator.Description.Website,
		Details:         validator.Description.Details,
		Rank:            rank,
		Identity:        validator.Description.Identity,
	}, nil
}

func (r *queryResolver) Uptimes(ctx context.Context, operatorAddress *string) (*model.UptimeResult, error) {
	block, err := document.Block{}.QueryLatestBlockFromDB()
	if err != nil {
		return &model.UptimeResult{}, nil
	}

	validator, err := document.GetValidatorByAddr(*operatorAddress)
	if err != nil {
		return &model.UptimeResult{}, nil
	}

	missedBlocks, err := document.MissedBlock{}.GetListMissedBlock(block.Height, validator.ConsensusAddres)
	if err != nil {
		return &model.UptimeResult{}, nil
	}
	var uptimeList []*model.Uptime
	for _, missedBlock := range missedBlocks {
		bytes, _ := missedBlock.Timestamp.MarshalText()
		uptimeList = append(uptimeList, &model.Uptime{
			Height:    int(missedBlock.Height),
			Timestamp: string(bytes),
		})
	}
	uptime := &model.UptimeResult{
		LastHeight: int(block.Height),
		Uptime:     uptimeList,
	}
	return uptime, nil
}

func (r *queryResolver) ProposedBlocks(ctx context.Context, before *int, size *int, operatorAddress string) ([]*model.Block, error) {
	blocks, totalRecord, err := document.Block{}.GetBlockListByOffsetAndSizeByOperatorAddress(*before, *size, operatorAddress)
	if err != nil {
		return []*model.Block{}, nil
	}
	return document.Block{}.FormatListBlockForModel(blocks, totalRecord)
}

func (r *queryResolver) PowerEvents(ctx context.Context, before *int, size *int, operatorAddress string) ([]*model.PowerEvent, error) {
	txs, err := document.CommonTx{}.GetListTxByAddress(*before, *size, operatorAddress)
	if err != nil {
		return []*model.PowerEvent{}, nil
	}
	return document.CommonTx{}.FormatListTxsForModelPowerEvent(txs, operatorAddress)
}

func (r *queryResolver) AccountTransactions(ctx context.Context, accAddress string, before int, size int) ([]*model.Tx, error) {
	listTxHash, err := document.AccountTransaction{}.GetListTxsByAddress(before, size, accAddress)
	if err != nil {
		return []*model.Tx{}, nil
	}
	txs, err := document.CommonTx{}.QueryByListByTxhash(listTxHash)
	if err != nil {
		return []*model.Tx{}, nil
	}
	return document.CommonTx{}.FormatListTxsForModel(txs)
}

func (r *queryResolver) AccountDetail(ctx context.Context, accAddress string) (*model.AccountDetail, error) {
	_, err := document.Validator{}.QueryValidatorDetailByAccAddr(accAddress)
	operatorAddress := utils.Convert(conf.Get().AddresPrefix+"valoper", accAddress)
	if err != nil {
		return &model.AccountDetail{
			IsValidator:     false,
			OperatorAddress: operatorAddress,
		}, nil
	}
	return &model.AccountDetail{
		IsValidator:     true,
		OperatorAddress: operatorAddress,
	}, nil
}

func (r *queryResolver) Proposals(ctx context.Context) ([]*model.Proposal, error) {
	proposals, err := document.Proposal{}.GetList()
	if err != nil {
		return []*model.Proposal{}, nil
	}
	return document.Proposal{}.FormatListProposalForModel(proposals)
}

func (r *queryResolver) ProposalDetail(ctx context.Context, proposalID int) (*model.Proposal, error) {
	proposal, err := document.Proposal{}.QueryProposalById(proposalID)
	if err != nil {
		return &model.Proposal{}, nil
	}

	validator, err := document.Validator{}.QueryValidatorDetailByAccAddr(proposal.Proposer)
	var moniker string
	if err == nil {
		moniker = validator.Description.Moniker
	}

	return document.Proposal{}.FormatProposalForModel(proposal, moniker)
}

func (r *queryResolver) Status(ctx context.Context) (*model.Status, error) {
	blocks, err := document.Block{}.QueryBlockOrderByHeightDesc(2)
	if err != nil {
		return &model.Status{}, nil
	}

	bondedToken, err := document.Validator{}.GetSumBondedToken()
	if err != nil {
		return &model.Status{}, nil
	}

	totalNumTxs, err := document.Block{}.GetCountTxs()
	if err != nil {
		return &model.Status{}, nil
	}

	totalSupplyToken, err := client.GetSupply()
	if err != nil {
		return &model.Status{}, nil
	}

	BlockTime := blocks[0].Time.UnixNano() - blocks[1].Time.UnixNano()
	bytes, _ := blocks[0].Time.MarshalText()
	return &model.Status{
		BlockHeight:       int(blocks[0].Height),
		Timestamp:         string(bytes),
		BlockTime:         int(BlockTime),
		BondedTokens:      int(bondedToken),
		TotalTxsNum:       int(totalNumTxs),
		TotalSupplyTokens: totalSupplyToken,
	}, nil
}

func (r *queryResolver) Inflation(ctx context.Context) (*model.Inflation, error) {
	return client.GetInflation()
}

func (r *queryResolver) Balances(ctx context.Context, accAddress string) (*model.Balances, error) {
	return client.GetBalances(accAddress)
}

func (r *queryResolver) Rewards(ctx context.Context, accAddress string) (*model.Rewards, error) {
	return client.GetRewards(accAddress)
}

func (r *queryResolver) Commission(ctx context.Context, operatorAddress string) (*model.Commission, error) {
	return client.GetCommission(operatorAddress)
}

func (r *queryResolver) Delegations(ctx context.Context, accAddress string) ([]*model.Delegation, error) {
	delegationResult, err := client.GetDelegation(accAddress)
	if err != nil {
		return []*model.Delegation{}, nil
	}

	var listDelegation []*model.Delegation
	var mapOperatorMoniker = document.Validator{}.GetListMapOperatorAndMoniker(delegationResult)
	for _, delegation := range delegationResult.DelegationResponses {
		amount, err := utils.ParseIntBase(delegation.Balance.Amount)
		if !err {
			fmt.Println("Can not parse amount delegations from string to int")
			amount = 0
		}
		var moniker string
		if v, ok := mapOperatorMoniker[delegation.Delegation.ValidatorAddress]; ok {
			moniker = v
		}
		t := &model.Delegation{
			DelegatorAddress: delegation.Delegation.DelegatorAdrress,
			ValidatorAddress: delegation.Delegation.ValidatorAddress,
			Amount:           amount,
			Moniker:          moniker,
		}
		listDelegation = append(listDelegation, t)
	}
	return listDelegation, nil
}

func (r *queryResolver) Unbonding(ctx context.Context, accAddress string) (*model.Unbonding, error) {
	result, err := client.GetUnbonding(accAddress)
	listOpAddress := client.GetListAccAddressFromUnbonding(result)
	validators, _ := document.Validator{}.QueryValidatorMonikerOpAddr(listOpAddress)
	mapAccAndMoniker := document.Validator{}.MapOperatorAndMoniker(validators)
	for _, item := range result.UnbondingResponses {
		if v, found := mapAccAndMoniker[item.ValidatorAddress]; found {
			item.Moniker = v
		}
	}
	return result, err
}

func (r *queryResolver) Redelegations(ctx context.Context, accAddress string) (*model.Redelegations, error) {
	result, err := client.GetRedelegation(accAddress)
	listOpAddress := client.GetListAccAddressFromRedelegation(result)
	validators, _ := document.Validator{}.QueryValidatorMonikerOpAddr(listOpAddress)
	mapAccAndMoniker := document.Validator{}.MapOperatorAndMoniker(validators)
	for _, item := range result.RedelegationResponses {
		if v, found := mapAccAndMoniker[item.Redelegation.ValidatorSrcAddress]; found {
			item.Redelegation.MonikerSrc = v
		}

		if v, found := mapAccAndMoniker[item.Redelegation.ValidatorDstAddress]; found {
			item.Redelegation.MonikerDst = v
		}
	}
	return result, err
}

func (r *queryResolver) Deposit(ctx context.Context, proposalID int) ([]*model.Deposit, error) {
	txs, _ := document.CommonTx{}.QueryProposalDeposit(proposalID)
	var listDeposit []*model.Deposit
	for _, item := range txs {
		bytes, _ := item.Timestamp.MarshalText()
		amount := document.CommonTx{}.GetValueOfLog(item.Logs, "transfer", "amount")
		depositor := document.CommonTx{}.GetValueOfLog(item.Logs, "transfer", "sender")
		t := &model.Deposit{
			TxHash:    item.TxHash,
			Time:      string(bytes),
			Amount:    &amount,
			Depositor: depositor,
		}
		listDeposit = append(listDeposit, t)
	}
	return listDeposit, nil
}

func (r *queryResolver) Vote(ctx context.Context, before *int, size *int, proposalID int) ([]*model.Vote, error) {
	txs, _ := document.CommonTx{}.QueryProposalVote(*before, *size, proposalID)
	var listVote []*model.Vote
	var listAccAddress []string
	for _, item := range txs {
		listAccAddress = append(listAccAddress, document.CommonTx{}.GetValueOfLog(item.Logs, "message", "sender"))
	}
	mapAccAndMoniker := document.Validator{}.GetListMapAccAndMoniker(listAccAddress)
	for _, item := range txs {
		bytes, _ := item.Timestamp.MarshalText()
		option := document.CommonTx{}.GetValueOfLog(item.Logs, "proposal_vote", "option")
		voter := document.CommonTx{}.GetValueOfLog(item.Logs, "message", "sender")

		var moniker string
		if v, ok := mapAccAndMoniker[voter]; ok {
			moniker = v
		}

		t := &model.Vote{
			TxHash:  item.TxHash,
			Time:    string(bytes),
			Option:  option,
			Voter:   voter,
			Moniker: moniker,
		}
		listVote = append(listVote, t)
	}
	return listVote, nil
}

func (r *queryResolver) Price(ctx context.Context, slug string) (*model.Price, error) {
	last, _ := document.StatAssetInfoList20Minute{}.QueryLatestStatAssetFromDB()
	first, _ := document.StatAssetInfoList20Minute{}.QueryNewestFromTime(last.Timestamp)
	change := last.Price - first.Price
	percent_change_24h := (change / last.Price) * 100
	return &model.Price{
		Price:            fmt.Sprintf("%f", last.Price),
		PercentChange24h: fmt.Sprintf("%f", percent_change_24h),
		Volume24h:        fmt.Sprintf("%f", last.Volume24H),
		MarketCap:        fmt.Sprintf("%f", last.Marketcap),
	}, nil
}

func (r *queryResolver) StatsAssets(ctx context.Context) ([]*model.StatsAsset, error) {
	statsAssets, err := document.StatAssetInfoList20Minute{}.GetList()

	if err != nil {
		return []*model.StatsAsset{}, nil
	}
	var listAssets []*model.StatsAsset
	for _, item := range statsAssets {
		bytes, _ := item.Timestamp.MarshalText()
		t := &model.StatsAsset{
			Price:     fmt.Sprintf("%f", item.Price),
			MarketCap: fmt.Sprintf("%f", item.Marketcap),
			Volume24h: fmt.Sprintf("%f", item.Volume24H),
			Timestamp: string(bytes),
		}
		listAssets = append(listAssets, t)
	}
	return listAssets, nil
}

func (r *queryResolver) Delegators(ctx context.Context, operatorAddress string, offset int) (*model.DelegatorResponse, error) {
	delegationResult, err := client.GetDelegators(operatorAddress, offset)
	if err != nil {
		return &model.DelegatorResponse{}, nil
	}

	var delegators []*model.Delegator
	for _, delegation := range delegationResult.DelegationResponses {
		t := &model.Delegator{
			DelegatorAddress: delegation.Delegation.DelegatorAdrress,
			Amount:           delegation.Balance.Amount,
		}
		delegators = append(delegators, t)
	}

	totalCount, bError := utils.ParseIntBase(delegationResult.Pagination.Total)
	if !bError {
		fmt.Println("Can not parse amount delegations from string to int")
		totalCount = 0
	}
	return &model.DelegatorResponse{
		TotalCount: totalCount,
		Delegators: delegators,
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
