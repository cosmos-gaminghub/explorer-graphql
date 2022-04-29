package document

import (
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionContract = "contracts"

	ContractAddressField = "contract_address"
	ContractLabelField   = "label"
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

func (_ Contract) FindByContractAddress(contractAddress string) (Contract, error) {
	var contract Contract
	condition := bson.M{
		ContractAddressField: contractAddress,
	}

	err := queryOne(CollectionContract, nil, condition, &contract)
	return contract, err
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

func (_ Contract) FormatForModel(result []Contract) ([]*model.Contract, error) {
	var listContract []*model.Contract
	for _, item := range result {
		listContract = append(listContract, &model.Contract{
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
		})
	}
	return listContract, nil
}
