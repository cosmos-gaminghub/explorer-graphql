package document

import (
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionProposal = "proposal"

	ProposalFieldProposalId = "proposal_id"
)

// Transaction defines the structure for transaction information.
type Proposal struct {
	ProposalId       int              `bson:"proposal_id"`
	ProposalStatus   string           `bson:"proposal_status"`
	Content          Content          `bson:"content" json:"content"`
	SubmitTime       time.Time        `bson:"submit_time"`
	FinalTallyResult FinalTallyResult `bson:"final_tally_result" json:"final_tally_result"`
	VotingEndTime    time.Time        `bson:"voting_end_time"`
	VotingStartTime  time.Time        `bson:"voting_start_time"`
	Proposer         string           `bson:"proposer"`
	TotalDeposit     []Amount         `bson:"total_deposit"`
}

type Content struct {
	Type        string `bson:"type"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Amount      []struct {
		Denom  string `bson:"denom"`
		Amount string `bson:"amount"`
	} `bson:"amount"`
}

type FinalTallyResult struct {
	Yes        string `bson:"yes"`
	Abstain    string `bson:"abstain"`
	No         string `bson:"no"`
	NoWithVeto string `bson:"no_with_veto"`
}

func (_ Proposal) GetList() ([]Proposal, error) {
	var proposals []Proposal

	sort := desc(ProposalFieldProposalId)
	err := queryAll(CollectionProposal, nil, nil, sort, 0, &proposals)
	return proposals, err
}

func (_ Proposal) QueryProposalById(ProposalId int) (Proposal, error) {
	condition := bson.M{ProposalFieldProposalId: ProposalId}
	var proposal Proposal
	err := queryOne(CollectionProposal, nil, condition, &proposal)

	return proposal, err
}

func (_ Proposal) FormatListProposalForModel(proposals []Proposal) (listProposal []*model.Proposal, err error) {
	for _, proposal := range proposals {
		t, _ := Proposal{}.FormatProposalForModel(proposal, "")
		listProposal = append(listProposal, t)
	}
	return listProposal, nil
}

func (_ Proposal) FormatProposalForModel(proposal Proposal, moniker string) (listProposal *model.Proposal, err error) {
	submitTime, _ := proposal.SubmitTime.MarshalText()
	vottingStart, _ := proposal.VotingStartTime.MarshalText()
	votingEndTime, _ := proposal.VotingEndTime.MarshalText()
	return &model.Proposal{
		Proposer:    proposal.Proposer,
		Status:      proposal.ProposalStatus,
		VotingStart: string(vottingStart),
		VotingEnd:   string(votingEndTime),
		Content: &model.Content{
			Title:       proposal.Content.Title,
			Description: proposal.Content.Description,
			Type:        proposal.Content.Type,
			Amount:      formatAmountForModelContent(proposal.Content),
		},
		Tally: &model.Tally{
			Yes:        proposal.FinalTallyResult.Yes,
			Abstain:    proposal.FinalTallyResult.Abstain,
			No:         proposal.FinalTallyResult.No,
			NoWithVeto: proposal.FinalTallyResult.NoWithVeto,
		},
		ID:           proposal.ProposalId,
		SubmitTime:   string(submitTime),
		TotalDeposit: formatTotalDepostForModel(proposal.TotalDeposit),
		Moniker:      moniker,
	}, err
}

func formatAmountForModelContent(content Content) (listAmount []*model.Amount) {
	if len(content.Amount) > 0 {
		for _, item := range content.Amount {
			c := &model.Amount{
				Denom:  item.Denom,
				Amount: item.Amount,
			}
			listAmount = append(listAmount, c)
		}
	}
	return listAmount
}

func formatTotalDepostForModel(totalDeposit []Amount) (td []*model.Amount) {
	if len(totalDeposit) > 0 {
		for _, item := range totalDeposit {
			d := &model.Amount{
				Denom:  item.Denom,
				Amount: item.Amount,
			}
			td = append(td, d)
		}
	}
	return td
}

// func formatDepositForModel(proposal Proposal) (listDeposit []*model.Deposit) {
// 	if len(proposal.Deposit) > 0 {
// 		for _, item := range proposal.Deposit {
// 			d := &model.Deposit{
// 				ProposalID: item.ProposalID,
// 				Depositor:  item.Depositor,
// 				Amount:     []*model.Amount{},
// 			}
// 			d.Amount = formatAmountForModel(item)
// 			listDeposit = append(listDeposit, d)
// 		}
// 	}
// 	return listDeposit
// }

// func formatVoteForModel(proposal Proposal) (listVote []*model.Vote) {
// 	if len(proposal.Deposit) > 0 {
// 		for _, item := range proposal.Vote {
// 			v := &model.Vote{
// 				ProposalID: item.ProposalId,
// 				Voter:      item.Voter,
// 				Option:     item.Option,
// 			}
// 			listVote = append(listVote, v)
// 		}
// 	}
// 	return listVote
// }

func formatAmountForModel(deposit Deposit) (listAmount []*model.Amount) {
	if len(deposit.Amount) > 0 {
		for _, a := range deposit.Amount {
			listAmount = append(listAmount, &model.Amount{
				Denom:  a.Denom,
				Amount: a.Amount,
			})
		}
	}
	return listAmount
}
