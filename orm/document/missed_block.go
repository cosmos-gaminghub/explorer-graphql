package document

import (
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/orm"
	"gopkg.in/mgo.v2/bson"
)

type MissedBlock struct {
	Height          int64     `bson:"height"`
	ConsensusAddres string    `bson:"consensus_address"`
	Timestamp       time.Time `bson:"timestamp"`
}

type Uptime struct {
	Id    string `bson:"_id"`
	Count int    `bson:"count"`
}

const (
	CollectionNmMissedBlock = "missed_block"

	MissedBlockFieldConsensusAddress = "consensus_address"
	MissedBlockFieldHeight           = "height"
	MissedBlockFieldTimeStamp        = "timestamp"

	DefaultOverBlocks = 100
)

func (_ MissedBlock) GetMissedBlockCount(listConsensusAddress []string) (map[string]int, int) {
	var result []Uptime
	var query = orm.NewQuery()
	defer query.Release()

	var upTimeMap = make(map[string]int)
	block, err := Block{}.QueryLatestBlockFromDB()
	if err != nil {
		return upTimeMap, DefaultOverBlocks
	}

	toHeight := block.Height
	fromHeight := toHeight - DefaultOverBlocks
	if fromHeight < 0 {
		fromHeight = 0
	}
	overBlock := int(toHeight - fromHeight)

	condition := bson.M{
		MissedBlockFieldConsensusAddress: bson.M{"$in": listConsensusAddress},
		MissedBlockFieldHeight: bson.M{
			"$gt":  fromHeight,
			"$lte": toHeight,
		},
	}

	selector := bson.M{MissedBlockFieldConsensusAddress: 1, MissedBlockFieldHeight: 1}

	query.Reset().
		SetResult(&result).
		SetCollection(CollectionNmMissedBlock).
		SetSelector(selector).
		PipeQuery(
			[]bson.M{
				{"$match": condition},
				{"$group": bson.M{
					"_id":   "$consensus_address",
					"count": bson.M{"$sum": 1},
				}},
			},
		)
	for _, missedBlock := range result {
		upTimeMap[missedBlock.Id] = missedBlock.Count
	}
	return upTimeMap, overBlock
}

func (_ MissedBlock) GetListMissedBlock(Height int64, ConsensusAddress string) ([]MissedBlock, error) {
	var result []MissedBlock

	toHeight := Height
	fromHeight := toHeight - DefaultOverBlocks
	if fromHeight < 0 {
		fromHeight = 0
	}

	condition := bson.M{
		MissedBlockFieldConsensusAddress: ConsensusAddress,
		MissedBlockFieldHeight: bson.M{
			"$gt":  fromHeight,
			"$lte": toHeight,
		},
	}

	selector := bson.M{MissedBlockFieldHeight: 1}
	queryAll(CollectionNmMissedBlock, selector, condition, "", 0, &result)
	return result, nil
}
