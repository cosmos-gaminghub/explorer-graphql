package document

import (
	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionAccountTransaction = "account_transaction"

	AccountTransaction_Field_Height    = "height"
	AccountTransaction_Account_Address = "account_address"
	AccountTransaction_Txhash          = "tx_hash"
)

type AccountTransaction struct {
	Height      int64  `bson:"height"`
	AccountAddr string `bson:"account_address"`
	TxHash      string `bson:"tx_hash"`
}

func (_ AccountTransaction) GetListTxsByAddress(before int, size int, address string) (listTxHash []string, err error) {
	var data []AccountTransaction

	selector := bson.M{
		AccountTransaction_Txhash: 1,
	}
	query := bson.M{AccountTransaction_Account_Address: address}
	if before != 0 {
		query[AccountTransaction_Field_Height] = bson.M{
			"$lt": before,
		}
	}

	err = querylistByOffsetAndSize(CollectionAccountTransaction, selector, query, desc(AccountTransaction_Field_Height), 0, size, &data)
	for _, item := range data {
		listTxHash = append(listTxHash, item.TxHash)
	}
	return listTxHash, err
}
