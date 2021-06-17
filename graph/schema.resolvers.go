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
	block, err := document.Block{}.QueryBlockByHeight(int64(*height))
	if err != nil {
		return &model.Block{}, nil
	}
	return &model.Block{
		Height:       int(block.Height),
		Hash:         block.Hash,
		Time:         block.Time.String(),
		ProposerAddr: block.ProposalAddr,
		NumTxs:       int(block.NumTxs)}, nil
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
	var listValidator []*model.Validator
	for index, validator := range validators {
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
		uptimeList = append(uptimeList, &model.Uptime{
			Height:    int(missedBlock.Height),
			Timestamp: missedBlock.Timestamp.String(),
		})
	}
	uptime := &model.UptimeResult{
		LastHeight: int(block.Height),
		Uptime:     uptimeList,
	}
	return uptime, nil
}

func (r *queryResolver) ProposedBlocks(ctx context.Context, offset *int, size *int, operatorAddress *string) ([]*model.Block, error) {
	blocks, err := document.Block{}.GetBlockListByOffsetAndSizeByOperatorAddress(*offset, *size, *operatorAddress)
	if err != nil {
		return []*model.Block{}, nil
	}
	return document.Block{}.FormatListBlockForModel(blocks)
}

func (r *queryResolver) PowerEvents(ctx context.Context, offset *int, size *int, operatorAddress string) ([]*model.PowerEvent, error) {
	txs, err := document.CommonTx{}.GetListTxByAddress(*offset, *size, operatorAddress)
	if err != nil {
		return []*model.PowerEvent{}, nil
	}

	var listTx []*model.PowerEvent
	for _, tx := range txs {
		t := &model.PowerEvent{
			TxHash:    tx.TxHash,
			Height:    int(tx.Height),
			Timestamp: tx.Timestamp,
			Amount:    int(document.CommonTx{}.GetAmountFromLogs(tx.Logs, operatorAddress)),
			Type:      document.CommonTx{}.GetTypeTextFromLogs(tx.Logs, operatorAddress),
		}
		listTx = append(listTx, t)
	}
	return listTx, nil
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

	totalNumTxs, err := document.CommonTx{}.GetCountTxs()
	if err != nil {
		return &model.Status{}, nil
	}

	totalSupplyToken, err := client.GetSupply()
	if err != nil {
		return &model.Status{}, nil
	}
	return &model.Status{
		BlockHeight:       int(lastBlock.Height),
		BlockTime:         lastBlock.Time.String(),
		BondedTokens:      int(bondedToken),
		TotalTxsNum:       totalNumTxs,
		TotalSupplyTokens: totalSupplyToken,
	}, nil
}

func (r *queryResolver) Inflation(ctx context.Context) (*model.Inflation, error) {
	return client.GetInflation()
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
