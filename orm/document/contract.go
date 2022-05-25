package document

import (
	"fmt"
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"github.com/cosmos-gaminghub/exploder-graphql/orm"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionContract = "contracts"

	ContractAddressField          = "contract_address"
	ContractTxhashField           = "txhash"
	ContractLabelField            = "label"
	ContractCodeIdFieled          = "code_id"
	ContractField                 = "contract"
	ContractAdminField            = "admin"
	ContractCreatorField          = "creator"
	ContractExecutedCountField    = "executed_count"
	ContractInstantiatedAtField   = "instantiated_at"
	ContractPermissionField       = "permission"
	ContractPermittedAddressField = "permitted_address"
	ContractLastExecutedAtField   = "last_executed_at"
	ContractTxHashField           = "txhash"
	ContractVersionField          = "version"
	ContractMessagesField         = "messages"

	Contract_Field_Tx = "txs"
)

type Contract struct {
	CodeId           int       `bson:"code_id"`
	Contract         string    `bson:"contract"`
	ContractAddress  string    `bson:"contract_address"`
	Admin            string    `bson:"admin"`
	Creator          string    `bson:"creator"`
	ExecutedCount    int       `bson:"executed_count"`
	InstantiatedAt   time.Time `bson:"instantiated_at"`
	Label            string    `bson:"label"`
	LastExecutedAt   time.Time `bson:"last_executed_at"`
	Permission       string    `bson:"permission"`
	PermittedAddress string    `bson:"permitted_address"`
	TxHash           string    `bson:"txhash"`
	Version          string    `bson:"version"`
}

func (d Contract) Name() string {
	return CollectionContract
}

func (_ Contract) FindByContractAddress(contractAddress string) (bson.M, error) {
	var query = orm.NewQuery()
	defer query.Release()

	var condition = []bson.M{
		{
			"$match": bson.M{
				ContractAddressField: contractAddress,
			},
		},
		{
			"$lookup": bson.M{
				"from":         CollectionNmCommonTx,
				"localField":   ContractTxhashField,
				"foreignField": Tx_Field_Hash,
				"as":           Contract_Field_Tx,
			},
		},
		{
			"$unwind": "$" + Contract_Field_Tx,
		},
		{
			"$project": bson.M{
				"messages":          "$" + Contract_Field_Tx + ".messages",
				"code_id":           1,
				"contract":          1,
				"contract_address":  1,
				"admin":             1,
				"creator":           1,
				"executed_count":    1,
				"instantiated_at":   1,
				"label":             1,
				"last_executed_at":  1,
				"permission":        1,
				"permitted_address": 1,
				"txhash":            1,
				"version":           1,
			},
		},
	}

	result := bson.M{}
	err := query.SetResult(&result).
		SetCollection(CollectionContract).
		PipeQuery(
			condition,
		)
	return result, err
}

func (_ Contract) GetContractByLimitAndOffset(offset int, size int, keyword *string) ([]Contract, error) {
	query := bson.M{}
	if keyword != nil {
		query[ContractLabelField] = bson.RegEx{
			Pattern: *keyword,
			Options: "i",
		}
	}

	data := []Contract{}
	err := querylistByOffsetAndSize(CollectionContract, nil, query, "", offset, size, &data)
	return data, err
}

func (_ Contract) GetContractByCodeId(codeId int) ([]bson.M, error) {
	query := bson.M{}
	query[ContractCodeIdFieled] = codeId

	data := []bson.M{}
	var selector = bson.M{
		ContractAddressField: 1,
	}
	err := queryAll(CollectionContract, selector, query, "", 0, &data)
	return data, err
}

func (_ Contract) GetListContractAddressFromBson(bson []bson.M) (result []string) {
	for _, item := range bson {
		result = append(result, item[ContractAddressField].(string))
	}
	return result
}

func (_ Contract) FormatForModel(result []Contract) ([]*model.Contract, error) {
	var listContract []*model.Contract
	for _, item := range result {
		listContract = append(listContract, Contract{}.FormatForModelItem(item))
	}
	return listContract, nil
}

func (_ Contract) FormatForModelItem(item Contract) *model.Contract {
	return &model.Contract{
		CodeID:          item.CodeId,
		ContractAddress: item.ContractAddress,
		Label:           item.Label,
		Contract:        item.Contract,
		Creator:         item.Creator,
		ExecutedCount:   item.ExecutedCount,
		InstantiatedAt:  item.InstantiatedAt.String(),
		LastExecutedAt:  item.LastExecutedAt.String(),
		Txhash:          item.TxHash,
		Version:         item.Version,
	}
}

func (_ Contract) FormatBsonMForModelContractDetail(result bson.M) (*model.Contract, error) {
	return &model.Contract{
		CodeID:           int(result[ContractCodeIdFieled].(int)),
		Contract:         fmt.Sprintf("%v", result[ContractField]),
		ContractAddress:  fmt.Sprintf("%v", result[ContractAddressField]),
		Admin:            fmt.Sprintf("%v", result[ContractAdminField]),
		Creator:          fmt.Sprintf("%v", result[ContractCreatorField]),
		ExecutedCount:    int(result[ContractExecutedCountField].(int)),
		InstantiatedAt:   fmt.Sprintf("%v", result[ContractInstantiatedAtField]),
		Label:            fmt.Sprintf("%v", result[ContractLabelField]),
		LastExecutedAt:   fmt.Sprintf("%v", result[ContractLastExecutedAtField]),
		Permission:       fmt.Sprintf("%v", result[ContractPermissionField]),
		PermittedAddress: fmt.Sprintf("%v", result[ContractPermittedAddressField]),
		Txhash:           fmt.Sprintf("%v", result[ContractTxHashField]),
		Version:          fmt.Sprintf("%v", result[ContractVersionField]),
		Messages:         fmt.Sprintf("%v", result[ContractMessagesField]),
	}, nil
}
