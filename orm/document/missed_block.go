package document

import (
	"graphql-go-example/orm"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type MissedBlock struct {
	Height       int64     `bson:"height"`
	OperatorAddr string    `bson:"operator_address"`
	Timestamp    time.Time `bson:"timestamp"`
}

type Uptime struct {
	Id    string `bson:"_id"`
	Count int    `bson:"count"`
}

const (
	CollectionNmMissedBlock = "missed_block"

	MissedBlockFieldOperatorAddress = "operator_address"
	MissedBlockFieldHeight          = "height"
	MissedBlockFieldTimeStamp       = "timestamp"

	DefaultOverBlocks = 100
)

func (_ MissedBlock) GetMissedBlockCount(listOperatorAddress []string) (map[string]int, int) {
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
		MissedBlockFieldOperatorAddress: bson.M{"$in": listOperatorAddress},
		MissedBlockFieldHeight: bson.M{
			"$gt":  fromHeight,
			"$lte": toHeight,
		},
	}

	selector := bson.M{MissedBlockFieldOperatorAddress: 1, MissedBlockFieldHeight: 1}

	query.Reset().
		SetResult(&result).
		SetCollection(CollectionNmMissedBlock).
		SetSelector(selector).
		PipeQuery(
			[]bson.M{
				{"$match": condition},
				{"$group": bson.M{
					"_id":   "$operator_address",
					"count": bson.M{"$sum": 1},
				}},
			},
		)
	for _, missedBlock := range result {
		upTimeMap[missedBlock.Id] = overBlock - missedBlock.Count
	}
	return upTimeMap, overBlock
}

func (_ MissedBlock) GetListMissedBlock(Height int64, OperatorAddr string) ([]MissedBlock, error) {
	var result []MissedBlock

	toHeight := Height
	fromHeight := toHeight - DefaultOverBlocks
	if fromHeight < 0 {
		fromHeight = 0
	}

	condition := bson.M{
		MissedBlockFieldOperatorAddress: OperatorAddr,
		MissedBlockFieldHeight: bson.M{
			"$gt":  fromHeight,
			"$lte": toHeight,
		},
	}

	selector := bson.M{MissedBlockFieldHeight: 1}
	queryAll(CollectionNmMissedBlock, selector, condition, "", 0, &result)
	return result, nil
}
