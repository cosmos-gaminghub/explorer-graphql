package document

import (
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/orm"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmStatsAsset = "stats_asset"

	StatsAsset_Field_Time = "timestamp"
)

// StatAssetInfoList1H defines the schema for asset statistics in an hourly basis
type StatAssetInfoList20Minute struct {
	Price     float32   `bson:"price"`
	Marketcap float32   `bson:"market_cap"`
	Volume24H float32   `bson:"volumne_24h"`
	Timestamp time.Time `bson:"timestamp"`
}

func (_ StatAssetInfoList20Minute) GetList() ([]StatAssetInfoList20Minute, error) {
	var statsAssets []StatAssetInfoList20Minute

	sort := desc(StatsAsset_Field_Time)
	err := queryAll(CollectionNmStatsAsset, nil, nil, sort, 72, &statsAssets)

	return statsAssets, err
}

func (_ StatAssetInfoList20Minute) QueryLatestStatAssetFromDB() (StatAssetInfoList20Minute, error) {

	var statsAssets StatAssetInfoList20Minute

	sort := desc(StatsAsset_Field_Time)
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(CollectionNmStatsAsset).
		SetCondition(nil).
		SetSort(sort).
		SetResult(&statsAssets)

	err := query.Exec()
	if err == nil {
		return statsAssets, nil
	}

	return StatAssetInfoList20Minute{}, err
}

func (_ StatAssetInfoList20Minute) QueryNewestFromTime(time time.Time) (StatAssetInfoList20Minute, error) {

	var statsAssets StatAssetInfoList20Minute

	sort := asc(StatsAsset_Field_Time)
	condition := bson.M{
		StatsAsset_Field_Time: bson.M{
			"$gte": time.AddDate(0, -1, 0),
		},
	}
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(CollectionNmStatsAsset).
		SetCondition(condition).
		SetSort(sort).
		SetResult(&statsAssets)

	err := query.Exec()
	if err == nil {
		return statsAssets, nil
	}

	return StatAssetInfoList20Minute{}, err
}
