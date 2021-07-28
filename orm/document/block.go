package document

import (
	"fmt"
	"log"
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"github.com/cosmos-gaminghub/exploder-graphql/orm"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmBlock = "block"

	Block_Field_Height           = "height"
	Block_Field_Hash             = "block_hash"
	Block_Field_Time             = "timestamp"
	Block_Field_NumTxs           = "num_txs"
	Block_Field_Meta             = "meta"
	Block_Field_Block            = "block"
	Block_Field_Validators       = "validators"
	Block_Field_ProposalAddress  = "proposer"
	Block_Field_Moniker          = "moniker"
	Block_Field_Operator_Address = "operator_address"
	Block_Field_Date_Time        = "date_time"
)

type Block struct {
	Height       int64     `bson:"height" json:"height"`
	Hash         string    `bson:"block_hash" json:"block_hash"`
	Time         time.Time `bson:"timestamp" json:"timestamp"`
	NumTxs       int64     `bson:"num_txs" json:"num_txs"`
	ProposalAddr string    `bson:"proposer" json:"proposer"`
}

func (b Block) String() string {
	return fmt.Sprintf(`
		Height      :%v
		Hash        :%v
		Time        :%v
		NumTxs      :%v
		ProposalAddr:%v
		`, b.Height, b.Hash, b.Time, b.NumTxs, b.ProposalAddr)
}

func (_ Block) QueryBlockByHeight(height int64) (bson.M, error) {

	var query = orm.NewQuery()
	defer query.Release()

	var condition = []bson.M{
		{
			"$match": bson.M{
				Block_Field_Height: height,
			},
		},
		{
			"$lookup": bson.M{
				"from":         CollectionNmValidator,
				"localField":   Block_Field_ProposalAddress,
				"foreignField": ValidatorFieldProposerHashAddr,
				"as":           Block_Field_Validators,
			},
		},
		{
			"$unwind": "$" + Block_Field_Validators,
		},
		{
			"$project": bson.M{
				"moniker":    "$" + Block_Field_Validators + ".description.moniker",
				"height":     1,
				"timestamp":  1,
				"num_txs":    1,
				"block_hash": 1,
				"proposer":   1,
				Block_Field_Date_Time: bson.M{
					"$dateToString": bson.M{"format": "%G-%m-%dT%H:%M:%SZ", "date": "$timestamp"},
				},
				"operator_address": "$" + Block_Field_Validators + "." + ValidatorFieldOperatorAddress,
			},
		},
	}

	result := bson.M{}
	err := query.SetResult(&result).
		SetCollection(CollectionNmBlock).
		PipeQuery(
			condition,
		)
	return result, err
}

func (_ Block) GetBlockListByOffsetAndSize(offset, size int) ([]bson.M, error) {
	var query = orm.NewQuery()
	defer query.Release()

	var condition = []bson.M{
		{"$sort": bson.M{Block_Field_Height: -1}},
		{
			"$skip": offset,
		},
		{
			"$limit": size,
		},
		{
			"$lookup": bson.M{
				"from":         CollectionNmValidator,
				"localField":   Block_Field_ProposalAddress,
				"foreignField": ValidatorFieldProposerHashAddr,
				"as":           Block_Field_Validators,
			},
		},
		{
			"$unwind": "$" + Block_Field_Validators,
		},
		{
			"$project": bson.M{
				"moniker":          "$" + Block_Field_Validators + ".description.moniker",
				"operator_address": "$" + Block_Field_Validators + "." + ValidatorFieldOperatorAddress,
				"height":           1, "timestamp": 1, "num_txs": 1, "block_hash": 1,
				"proposer": 1,
				Block_Field_Date_Time: bson.M{
					"$dateToString": bson.M{"format": "%G-%m-%dT%H:%M:%SZ", "date": "$timestamp"},
				},
			},
		},
	}

	results := []bson.M{}
	err := query.SetResult(&results).
		SetCollection(CollectionNmBlock).
		PipeQuery(
			condition,
		)
	return results, err
}

func (_ Block) GetBlockListByOffsetAndSizeByOperatorAddress(before int, size int, operatorAddress string) ([]Block, int, error) {
	validator, err := Validator{}.QueryValidatorDetailByOperatorAddr(operatorAddress)
	var totalRecord int
	if err != nil {
		return []Block{}, totalRecord, err
	}
	condition := bson.M{
		Block_Field_ProposalAddress: validator.ProposerAddr,
	}
	totalRecord, _ = Block{}.GetCountBlock(condition)

	if before != 0 {
		condition[Block_Field_Height] = bson.M{
			"$lt": before,
		}
	}

	var selector = bson.M{"height": 1, "timestamp": 1, "num_txs": 1, "block_hash": 1, "validators.pub_key": 1, "validators.address": 1,
		"validators.voting_power": 1, "block.last_commit.precommits.validator_address": 1, "meta.header.total_txs": 1, "proposer": 1}
	var blocks []Block

	sort := desc(Block_Field_Height)

	err = querylistByOffsetAndSize(CollectionNmBlock, selector, condition, sort, 0, size, &blocks)

	return blocks, totalRecord, err
}

func (_ Block) GetBlockListByPage(offset, size int, total bool) (int, []Block, error) {

	var selector = bson.M{"height": 1, "time": 1, "num_txs": 1, "hash": 1, "validators.address": 1, "validators.voting_power": 1, "block.last_commit.precommits.validator_address": 1, "meta.header.total_txs": 1}

	var blocks []Block

	sort := desc(Block_Field_Height)
	var cnt, err = pageQuery(CollectionNmBlock, selector, bson.M{"height": bson.M{"$gt": 0}}, sort, offset, size, total, &blocks)

	return cnt, blocks, err
}

func (_ Block) GetRecentBlockList() ([]Block, error) {
	var blocks []Block
	var selector = bson.M{"height": 1, "time": 1, "num_txs": 1}

	sort := desc(Block_Field_Height)
	err := queryAll(CollectionNmBlock, selector, nil, sort, 10, &blocks)
	return blocks, err
}

func (_ Block) QueryOneBlockOrderByHeightAsc() (Block, error) {

	var blocks []Block

	err := queryAll(CollectionNmBlock, nil, nil, asc(Block_Field_Height), 1, &blocks)

	if len(blocks) == 1 {
		return blocks[0], err
	}

	return Block{}, err
}

func (_ Block) QueryLatestBlockFromDB() (Block, error) {

	var block Block
	var selector = bson.M{"height": 1}

	sort := desc(Block_Field_Height)
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(CollectionNmBlock).
		SetCondition(nil).
		SetSelector(selector).
		SetSort(sort).
		SetResult(&block)

	err := query.Exec()
	if err == nil {
		return block, nil
	} else {
		log.Fatalf("query db error")
	}

	return Block{}, err
}

func (_ Block) QueryBlockOrderByHeightDesc(size int) ([]Block, error) {

	db := orm.GetDatabase()
	defer db.Session.Close()

	var blocks []Block

	sort := desc(Block_Field_Height)
	err := querylistByOffsetAndSize(CollectionNmBlock, nil, nil, sort, 0, size, &blocks)
	return blocks, err
}

func (_ Block) QueryBlocksByDurationWithHeightAsc(startTime, endTime time.Time) ([]Block, error) {
	db := orm.GetDatabase()
	defer db.Session.Close()

	blocks := []Block{}
	err := db.C(CollectionNmBlock).Find(bson.M{"time": bson.M{"$gte": startTime, "$lt": endTime}}).Sort("height").All(&blocks)
	return blocks, err
}

func (_ Block) QueryValidatorsByHeightList(hArr []int64) ([]Block, error) {

	var selector = bson.M{Block_Field_Height: 1, Block_Field_Validators: 1}

	sort := desc(Block_Field_Height)
	var blocks []Block
	err := queryAll(CollectionNmBlock, selector, bson.M{"height": bson.M{"$in": hArr}}, sort, 0, &blocks)
	return blocks, err
}

func (_ Block) FormatListBlockForModel(blocks []Block, totalRecord int) ([]*model.Block, error) {
	var listBlock []*model.Block
	for _, block := range blocks {
		bytes, _ := block.Time.MarshalText()
		t := &model.Block{
			Height:       int(block.Height),
			Hash:         block.Hash,
			ProposerAddr: block.ProposalAddr,
			NumTxs:       int(block.NumTxs),
			Time:         string(bytes),
			TotalRecords: totalRecord,
		}
		listBlock = append(listBlock, t)
	}
	return listBlock, nil
}

func (_ Block) FormatBsonMForModel(results []bson.M) ([]*model.Block, error) {
	var listBlock []*model.Block
	for _, block := range results {
		t, _ := Block{}.FormatBsonMForModelBlockDetail(block)
		listBlock = append(listBlock, t)
	}
	return listBlock, nil
}

func (_ Block) FormatBsonMForModelBlockDetail(result bson.M) (*model.Block, error) {
	return &model.Block{
		Height:          int(result[Block_Field_Height].(int64)),
		Hash:            fmt.Sprintf("%v", result[Block_Field_Hash]),
		ProposerAddr:    fmt.Sprintf("%v", result[Block_Field_ProposalAddress]),
		NumTxs:          int(result[Block_Field_NumTxs].(int64)),
		Time:            fmt.Sprintf("%v", result[Block_Field_Date_Time]),
		Moniker:         fmt.Sprintf("%v", result[Block_Field_Moniker]),
		OperatorAddress: fmt.Sprintf("%v", result[Block_Field_Operator_Address]),
	}, nil
}

func (_ Block) GetCountBlock(condition bson.M) (int, error) {
	result := []bson.M{}
	var query = orm.NewQuery()
	defer query.Release()
	query.SetResult(&result).
		SetCollection(CollectionNmBlock).
		PipeQuery(
			[]bson.M{
				{"$match": condition},
				{"$group": bson.M{
					"_id":   "",
					"count": bson.M{"$sum": 1},
				}},
			},
		)
	if len(result) == 0 {
		return 0, nil
	}
	return result[0]["count"].(int), nil
}

func (_ Block) GetCountTxs() (int64, error) {
	result := bson.M{}
	var query = orm.NewQuery()
	defer query.Release()
	query.SetResult(&result).
		SetCollection(CollectionNmBlock).
		PipeQuery(
			[]bson.M{
				{
					"$match": bson.M{
						Block_Field_NumTxs: bson.M{"$gt": 0},
					},
				},
				{"$group": bson.M{
					"_id":   nil,
					"count": bson.M{"$sum": "$num_txs"},
				}},
			},
		)
	return result["count"].(int64), nil
}

type BlockMeta struct {
	BlockID BlockID `bson:"block_id"`
	Header  Header  `bson:"header"`
}

type BlockID struct {
	Hash        string        `bson:"hash"`
	PartsHeader PartSetHeader `bson:"parts"`
}

type PartSetHeader struct {
	Total int    `bson:"total"`
	Hash  string `bson:"hash"`
}

type Header struct {
	// basic block info
	ChainID string    `bson:"chain_id"`
	Height  int64     `bson:"height"`
	Time    time.Time `bson:"time"`
	NumTxs  int64     `bson:"num_txs"`

	// prev block info
	LastBlockID BlockID `bson:"last_block_id"`
	TotalTxs    int64   `bson:"total_txs"`

	// hashes of block data
	LastCommitHash string `bson:"last_commit_hash"` // commit from validators from the last block
	DataHash       string `bson:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash  string `bson:"validators_hash"`   // validators for the current block
	ConsensusHash   string `bson:"consensus_hash"`    // consensus params for current block
	AppHash         string `bson:"app_hash"`          // state after txs from the previous block
	LastResultsHash string `bson:"last_results_hash"` // root hash of all results from the txs from the previous block

	// consensus info
	EvidenceHash string `bson:"evidence_hash"` // evidence included in the block
}

type BlockContent struct {
	LastCommit Commit `bson:"last_commit"`
}

type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	BlockID    BlockID `bson:"block_id"`
	Precommits []Vote  `bson:"precommits"`
}

type Signature struct {
	Type  string `bson:"type"`
	Value string `bson:"value"`
}

type TmValidator struct {
	Address     string `bson:"address"`
	PubKey      string `bson:"pub_key"`
	VotingPower int64  `bson:"voting_power"`
	Accum       int64  `bson:"accum"`
}

type BlockResults struct {
	DeliverTx  []ResponseDeliverTx `bson:"deliver_tx"`
	EndBlock   ResponseEndBlock    `bson:""end_block""`
	BeginBlock ResponseBeginBlock  `bson:""begin_block""`
}

type ResponseDeliverTx struct {
	Code      uint32   `bson:"code"`
	Data      string   `bson:"data"`
	Log       string   `bson:"log"`
	Info      string   `bson:"info"`
	GasWanted int64    `bson:"gas_wanted"`
	GasUsed   int64    `bson:"gas_used"`
	Tags      []KvPair `bson:"tags"`
	Codespace string   `bson:"codespace"`
}

type ResponseEndBlock struct {
	ValidatorUpdates      []ValidatorUpdate `bson:"validator_updates"`
	ConsensusParamUpdates ConsensusParams   `bson:"consensus_param_updates"`
	Tags                  []KvPair          `bson:"tags"`
}

type ValidatorUpdate struct {
	PubKey string `bson:"pub_key"`
	Power  int64  `bson:"power"`
}

type ConsensusParams struct {
	BlockSize BlockSizeParams `bson:"block_size"`
	Evidence  EvidenceParams  `bson:"evidence"`
	Validator ValidatorParams `bson:"validator"`
}

type BlockSizeParams struct {
	MaxBytes int64 `bson:"max_bytes"`
	MaxGas   int64 `bson:"max_gas"`
}

type EvidenceParams struct {
	MaxAge int64 `bson:"max_age"`
}

type ValidatorParams struct {
	PubKeyTypes []string `bson:"pub_key_types`
}

type ResponseBeginBlock struct {
	Tags []KvPair `bson:"tags"`
}

type KvPair struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func (d Block) Name() string {
	return CollectionNmBlock
}

func (d Block) PkKvPair() map[string]interface{} {
	return bson.M{Block_Field_Height: d.Height}
}

type ResValidatorPreCommits struct {
	Address       string `bson:"_id"`
	PreCommitsNum int64  `bson:"num"`
}
