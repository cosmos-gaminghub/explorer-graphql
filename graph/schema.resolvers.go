package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/cosmos-gaminghub/exploder-graphql/client"
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

	listOperatorAddress := document.Validator{}.GetListOperatorAdress(validators)
	upTimeCount, overBlocks := document.MissedBlock{}.GetMissedBlockCount(listOperatorAddress)

	validatorFormat := document.Validator{}.FormatListValidator(validators)
	var listValidator []*model.Validator
	for index, validator := range validatorFormat {
		commision, _ := utils.ParseStringToFloat(validator.Commission.CommissionRate.Rate)
		t := &model.Validator{
			Moniker:         validator.Description.Moniker,
			OperatorAddress: validator.OperatorAddr,
			AccAddress:      validator.AccountAddr,
			VotingPower:     int(validator.Tokens),
			Commission:      commision,
			Jailed:          validator.Jailed,
			Status:          validator.Status,
			Uptime:          upTimeCount[validator.OperatorAddr],
			OverBlocks:      overBlocks,
			Website:         validator.Description.Website,
			Rank:            index + 1,
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
	return &model.Validator{
		Moniker:         validator.Description.Moniker,
		OperatorAddress: validator.OperatorAddr,
		AccAddress:      validator.AccountAddr,
		VotingPower:     int(validator.Tokens),
		Commission:      commision,
		Jailed:          validator.Jailed,
		Status:          validator.Status,
		Uptime:          upTimeCount[validator.OperatorAddr],
		OverBlocks:      overBlocks,
		Website:         validator.Description.Website,
		Details:         validator.Description.Details,
		Rank:            rank,
	}, nil
}

func (r *queryResolver) Uptimes(ctx context.Context, operatorAddress *string) (*model.UptimeResult, error) {
	block, err := document.Block{}.QueryLatestBlockFromDB()
	if err != nil {
		return &model.UptimeResult{}, nil
	}
	missedBlocks, err := document.MissedBlock{}.GetListMissedBlock(block.Height, *operatorAddress)
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
	blocks, err := document.Block{}.GetBlockListByOffsetAndSizeByOperatorAddress(*before, *size, operatorAddress)
	if err != nil {
		return []*model.Block{}, nil
	}
	return document.Block{}.FormatListBlockForModel(blocks)
}

func (r *queryResolver) PowerEvents(ctx context.Context, before *int, size *int, operatorAddress string) ([]*model.PowerEvent, error) {
	txs, err := document.CommonTx{}.GetListTxByAddress(*before, *size, operatorAddress)
	if err != nil {
		return []*model.PowerEvent{}, nil
	}
	return document.CommonTx{}.FormatListTxsForModelPowerEvent(txs, operatorAddress)
}

func (r *queryResolver) AccountTransactions(ctx context.Context, accAddress *string) ([]*model.Tx, error) {
	txs, err := document.CommonTx{}.GetListTxByAccountAddress(*accAddress)
	if err != nil {
		return []*model.Tx{}, nil
	}
	return document.CommonTx{}.FormatListTxsForModel(txs)
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
	return document.Proposal{}.FormatProposalForModel(proposal)
}

func (r *queryResolver) Status(ctx context.Context) (*model.Status, error) {
	lastBlock, err := document.Block{}.QueryOneBlockOrderByHeightDesc()
	if err != nil {
		return &model.Status{}, nil
	}

	bondedToken, err := document.Validator{}.GetSumBondedToken()
	if err != nil {
		return &model.Status{}, nil
	}

	totalNumTxs, err := document.CommonTx{}.GetCountTxs(nil)
	if err != nil {
		return &model.Status{}, nil
	}

	totalSupplyToken, err := client.GetSupply()
	if err != nil {
		return &model.Status{}, nil
	}

	bytes, _ := lastBlock.Time.MarshalText()
	return &model.Status{
		BlockHeight:       int(lastBlock.Height),
		BlockTime:         string(bytes),
		BondedTokens:      int(bondedToken),
		TotalTxsNum:       totalNumTxs,
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

func (r *queryResolver) Delegations(ctx context.Context, accAddress *string) ([]*model.Delegation, error) {
	delegationResult, err := client.GetDelegation(*accAddress)
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

func (r *queryResolver) Unbonding(ctx context.Context, accAddress *string) (*model.Unbonding, error) {
	return client.GetUnbonding(*accAddress)
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
	for _, item := range txs {
		bytes, _ := item.Timestamp.MarshalText()
		option := document.CommonTx{}.GetValueOfLog(item.Logs, "proposal_vote", "option")
		voter := document.CommonTx{}.GetValueOfLog(item.Logs, "message", "sender")
		t := &model.Vote{
			TxHash: item.TxHash,
			Time:   string(bytes),
			Option: option,
			Voter:  voter,
		}
		listVote = append(listVote, t)
	}
	return listVote, nil
}

func (r *queryResolver) Price(ctx context.Context, slug string) (*model.Price, error) {
	return client.GetConcurrencyQuoteLastest(slug)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
