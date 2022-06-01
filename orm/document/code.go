package document

import (
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionCode = "codes"

	CodeIdField = "code_id"
)

type Code struct {
	CodeId           int       `bson:"code_id"`
	Contract         string    `bson:"contract"`
	DataHash         string    `bson:"data_hash"`
	CreatedAt        time.Time `bson:"created_at"`
	Creator          string    `bson:"creator"`
	InstantiateCount int       `bson:"instantiate_count"`
	Permission       string    `bson:"permission"`
	PermittedAddress string    `bson:"permitted_address"`
	TxHash           string    `bson:"txhash"`
	Version          string    `bson:"version"`
}

func (_ Code) FindByCodeId(codeId int) (Code, error) {
	var code Code
	condition := bson.M{
		CodeIdField: codeId,
	}

	err := queryOne(CollectionCode, nil, condition, &code)
	return code, err
}

func (_ Code) GetCodeListByOffsetAndSize(before int, size int) ([]Code, error) {
	var data []Code

	query := bson.M{}
	if before != 0 {
		query[CodeIdField] = bson.M{
			"$gt": before,
		}
	}

	err := querylistByOffsetAndSize(CollectionCode, nil, query, asc(CodeIdField), 0, size, &data)
	return data, err
}

func (_ Code) FormatForModel(results []Code) ([]*model.Code, error) {
	var listCode []*model.Code
	for _, item := range results {
		t := Code{}.FormatModelCodeItem(item)
		listCode = append(listCode, t)
	}
	return listCode, nil
}

func (_ Code) FormatModelCodeItem(item Code) *model.Code {
	return &model.Code{
		CodeID:           item.CodeId,
		Contract:         item.Contract,
		DataHash:         item.DataHash,
		CreatedAt:        item.CreatedAt.String(),
		Creator:          item.Creator,
		InstantiateCount: item.InstantiateCount,
		Permission:       item.Permission,
		PermittedAddress: item.PermittedAddress,
		Txhash:           item.TxHash,
		Version:          item.Version,
	}
}
