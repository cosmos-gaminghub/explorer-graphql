package document

import (
	"time"
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
