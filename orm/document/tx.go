package document

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos-gaminghub/exploder-graphql/graph/model"
	"github.com/cosmos-gaminghub/exploder-graphql/orm"
	"github.com/cosmos-gaminghub/exploder-graphql/utils"

	"gopkg.in/mgo.v2/bson"
)

const (
	CollectionNmCommonTx = "txs"
	TxStatusSuccess      = "success"
	TxStatusFail         = "fail"

	Tx_Field_Time        = "timestamp"
	Tx_Field_Height      = "height"
	Tx_Field_Hash        = "txhash"
	Tx_Field_From        = "from"
	Tx_Field_To          = "to"
	Tx_Field_Signers     = "signers"
	Tx_Field_Amount      = "amount"
	Tx_Field_Type        = "logs.events.type"
	Tx_Field_Value       = "logs.events.attributes.value"
	Tx_Field_Event_Type  = "logs.events.type"
	Tx_Field_Event_Value = "logs.events.attributes.value"
	Tx_Field_Event_Key   = "logs.events.attributes.key"
	Tx_Field_Fee         = "fee"
	Tx_Field_Memo        = "memo"
	Tx_Field_Status      = "status"
	Tx_Field_Code        = "code"
	Tx_Field_Log         = "log"
	Tx_Field_GasUsed     = "gas_used"
	Tx_Field_GasPrice    = "gas_price"
	Tx_Field_ActualFee   = "actual_fee"
	Tx_Field_ProposalId  = "proposal_id"
	Tx_Field_Tags        = "tags"
	Tx_Field_Msgs        = "msgs"

	Tx_Field_Msgs_UdInfo         = "msgs.msg.ud_info.source"
	Tx_Field_Msgs_Moniker        = "msgs.msg.moniker"
	Tx_Field_Msgs_UdInfo_Symbol  = "msgs.msg.ud_info.symbol"
	Tx_Field_Msgs_UdInfo_Gateway = "msgs.msg.ud_info.gateway"
	Tx_Field_Msgs_Hashcode       = "msgs.msg.hash_lock"
	Tx_AssetType_Native          = "native"
	Tx_AssetType_Gateway         = "gateway"

	Tx_Asset_TxType_Issue                = "IssueToken"
	Tx_Asset_TxType_Edit                 = "EditToken"
	Tx_Asset_TxType_Mint                 = "MintToken"
	Tx_Asset_TxType_TransferOwner        = "TransferTokenOwner"
	Tx_Asset_TxType_TransferGatewayOwner = "TransferGatewayOwner"

	TypeUnBond     = "unbond"
	TypeDelegate   = "delegate"
	TypeReDelegate = "redelegate"

	TypeForDeposit = "proposal_deposit"
	TypeForVote    = "proposal_vote"
)

type Signer struct {
	AddrHex    string `bson:"addr_hex"`
	AddrBech32 string `bson:"addr_bech32"`
}

type Coin struct {
	Denom  string  `bson:"denom"`
	Amount float64 `bson:"amount"`
}

func (c Coin) Add(a Coin) Coin {
	if c.Denom == a.Denom {
		return Coin{
			Denom:  c.Denom,
			Amount: c.Amount + a.Amount,
		}
	}
	return c
}

type Coins []Coin

type ActualFee struct {
	Denom  string  `bson:"denom"`
	Amount float64 `bson:"amount"`
}

type CommonTx struct {
	Height     int64     `bson:"height"`
	TxHash     string    `bson:"txhash"`
	Code       uint32    `bson:"code"`
	Memo       string    `bson:"memo"`
	GasWanted  int64     `bson:"gas_wanted"`
	GasUsed    int64     `bson:"gas_used"`
	Timestamp  time.Time `bson:"timestamp"`
	Logs       []Log     `bson:"logs" json:"logs"`
	Fee        Fee       `bson:"fee" json:"fee"`
	Signatures []string  `bson:"signatures" json:"signatures"`
	Messages   string    `bson:"messages" json:"messages"`
}

type Fee struct {
	Amount []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `json:"amount"`
	GasLimit string `json:"gas_limit"`
	Granter  string `json:"granter"`
	Payer    string `json:"payer"`
}

type Messages struct {
	Type        string `json:"@type"`
	FromAddress string `bson:"from_address" json:"from_address"`
	ToAddress   string `bson:"to_address" json:"to_address"`
	Amount      []struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	} `bson:"amount" json:"amount"`
}

func (tx CommonTx) String() string {
	return ""

}

type Log struct {
	MsgIndex int     `bson:"msg_index"`
	Log      string  `bson:"log"`
	Events   []Event `bson:"events"`
}

type Event struct {
	Type       string           `bson:"type"`
	Attributes []EventAttribute `bson:"attributes"`
}

type EventAttribute struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func (_ CommonTx) GetListTxBy(size int) ([]CommonTx, error) {
	var data []CommonTx

	err := queryAll(CollectionNmCommonTx, nil, nil, desc(Tx_Field_Height), size, &data)
	return data, err
}

func (_ CommonTx) GetListTxByAddress(before int, size int, operatorAddress string) ([]CommonTx, error) {
	var data []CommonTx
	query := bson.M{Tx_Field_Value: operatorAddress}
	typeArr := []string{TypeDelegate, TypeUnBond, TypeReDelegate}
	query[Tx_Field_Type] = bson.M{
		"$in": typeArr,
	}
	if before != 0 {
		query[Tx_Field_Height] = bson.M{
			"$lt": before,
		}
	}

	err := querylistByOffsetAndSize(CollectionNmCommonTx, nil, query, desc(Tx_Field_Height), 0, size, &data)
	return data, err
}

func (_ CommonTx) GetListTxByAccountAddress(accAddress string) ([]CommonTx, error) {
	var data []CommonTx
	query := bson.M{Tx_Field_Value: accAddress}

	err := queryAll(CollectionNmCommonTx, nil, query, desc(Tx_Field_Height), 0, &data)
	return data, err
}

func (_ CommonTx) GetAmountFromLogs(logs []Log, operatorAddress string) int64 {
	typeArr := []string{TypeDelegate, TypeUnBond, TypeReDelegate}
	var amount int64
	for _, log := range logs {
		for _, event := range log.Events {
			if utils.Contains(typeArr, event.Type) {
				for _, attribute := range event.Attributes {
					if attribute.Key == "amount" {
						value, err := utils.ParseInt(attribute.Value)
						if err {
							amount = value
						}
						break
					}
				}
			}
		}
	}
	return amount
}

func (_ CommonTx) GetTypeTextFromLogs(logs []Log, operatorAddress string) string {
	typeArr := []string{TypeDelegate, TypeUnBond, TypeReDelegate}
	typeText := "add"
	for _, log := range logs {
		for _, event := range log.Events {
			if utils.Contains(typeArr, event.Type) && event.Type == TypeUnBond {
				typeText = "minus"
				break
			}
		}
	}
	return typeText
}

func (_ CommonTx) FormatListTxsForModelPowerEvent(txs []CommonTx, operatorAddress string) ([]*model.PowerEvent, error) {
	query := bson.M{Tx_Field_Value: operatorAddress}
	typeArr := []string{TypeDelegate, TypeUnBond, TypeReDelegate}
	query[Tx_Field_Type] = bson.M{
		"$in": typeArr,
	}
	totalRecord, _ := CommonTx{}.GetCountTxs(query)

	var listTx []*model.PowerEvent
	for _, tx := range txs {
		bytes, _ := tx.Timestamp.MarshalText()
		t := &model.PowerEvent{
			TxHash:       tx.TxHash,
			Height:       int(tx.Height),
			Timestamp:    string(bytes),
			Amount:       int(CommonTx{}.GetAmountFromLogs(tx.Logs, operatorAddress)),
			Type:         CommonTx{}.GetTypeTextFromLogs(tx.Logs, operatorAddress),
			TotalRecords: totalRecord,
		}
		listTx = append(listTx, t)
	}
	return listTx, nil
}

func (_ CommonTx) FormatListTxsForModel(txs []CommonTx) ([]*model.Tx, error) {
	var listTx []*model.Tx
	for _, tx := range txs {
		t, _ := CommonTx{}.FormatTxForModel(tx)
		listTx = append(listTx, t)
	}
	return listTx, nil
}

func (_ CommonTx) FormatTxForModel(tx CommonTx) (*model.Tx, error) {
	fee, err := json.Marshal(tx.Fee)
	if err != nil {
		fee = []byte{}
	}

	logs, err := json.Marshal(tx.Logs)
	if err != nil {
		logs = []byte{}
	}
	bytes, _ := tx.Timestamp.MarshalText()
	t := &model.Tx{
		TxHash:    tx.TxHash,
		Height:    int(tx.Height),
		Timestamp: string(bytes),
		Status:    int(tx.Code),
		Fee:       string(fee),
		Messages:  tx.Messages,
		Logs:      string(logs),
		GasUsed:   int(tx.GasUsed),
		GasWanted: int(tx.GasWanted),
		Memo:      tx.Memo,
	}
	return t, nil
}

func (_ CommonTx) QueryByPage(query bson.M, pageNum, pageSize int, istotal bool) (int, []CommonTx, error) {
	var data []CommonTx

	total, err := pageQuery(CollectionNmCommonTx, nil, query, desc(Tx_Field_Time), pageNum, pageSize, istotal, &data)

	return total, data, err
}

func (_ CommonTx) QueryHashActualFeeType() ([]CommonTx, error) {

	var selector = bson.M{"time": 1, "tx_hash": 1, "actual_fee": 1, "type": 1}
	var txs []CommonTx

	err := queryAll(CollectionNmCommonTx, selector, nil, desc(Tx_Field_Time), 10, &txs)
	return txs, err
}

// func (_ CommonTx) QueryHashTimeByProposalIdVoters(proposalid int64, voters []string) ([]CommonTx, error) {

// 	var selector = bson.M{Tx_Field_Time: 1, Tx_Field_Hash: 1, Tx_Field_From: 1, Tx_Field_ProposalId: 1}
// 	var query = bson.M{Tx_Field_Type: types.TxTypeVote, Tx_Field_Status: TxStatusSuccess, Tx_Field_ProposalId: proposalid}
// 	if len(voters) > 0 {
// 		query[Tx_Field_From] = bson.M{"$in": voters}
// 	}
// 	var txs []CommonTx

// 	err := queryAll(CollectionNmCommonTx, selector, query, desc(Tx_Field_Time), 0, &txs)
// 	return txs, err
// }

func (_ CommonTx) QueryTxByHash(hash string) (CommonTx, error) {
	dbm := getDb()
	defer dbm.Session.Close()

	var result CommonTx
	query := bson.M{}
	query[Tx_Field_Hash] = hash
	err := dbm.C(CollectionNmCommonTx).Find(query).Sort(desc(Tx_Field_Height)).One(&result)

	return result, err
}

func (_ CommonTx) QueryTxByHeight(height int64) ([]CommonTx, error) {
	dbm := getDb()
	defer dbm.Session.Close()

	var result []CommonTx
	query := bson.M{}
	query[Tx_Field_Height] = height
	err := dbm.C(CollectionNmCommonTx).Find(query).Sort(desc(Tx_Field_Height)).All(&result)

	return result, err
}

func (_ CommonTx) QueryHtlcTx(query bson.M) (CommonTx, error) {
	dbm := getDb()
	defer dbm.Session.Close()

	var result CommonTx
	err := dbm.C(CollectionNmCommonTx).Find(query).Sort(desc(Tx_Field_Time)).One(&result)

	return result, err
}

type Counter []struct {
	Type  string `bson:"_id,omitempty"`
	Count int
}

func (cArr Counter) String() string {
	res := ""
	for k, v := range cArr {
		res += fmt.Sprintf("idx: %v Type  :%v  \t	Count :%v \n", k, v.Type, v.Count)
	}
	return res
}

// func (_ CommonTx) CountByType(query bson.M) (Counter, error) {

// 	var typeArr []string
// 	typeArr = append(typeArr, types.BankList...)
// 	typeArr = append(typeArr, types.DeclarationList...)
// 	typeArr = append(typeArr, types.StakeList...)
// 	typeArr = append(typeArr, types.GovernanceList...)
// 	query[Tx_Field_Type] = bson.M{
// 		"$in": typeArr,
// 	}
// 	query = FilterUnknownTxs(query)

// 	counter := Counter{}

// 	c := getDb().C(CollectionNmCommonTx)
// 	defer c.Database.Session.Close()

// 	pipe := c.Pipe(
// 		[]bson.M{
// 			{"$match": query},
// 			{"$group": bson.M{
// 				"_id":   "$type",
// 				"count": bson.M{"$sum": 1},
// 			}},
// 		},
// 	)

// 	err := pipe.All(&counter)

// 	return counter, err
// }

// func (_ CommonTx) GetTxlistByDuration(startTime, endTime string) ([]TxNumStat, error) {

// 	query := bson.M{}
// 	query[TxNumStat_Field_Date] = bson.M{"$gte": startTime, "$lt": endTime}

// 	var selector = bson.M{
// 		TxNumStat_Field_Date: 1, TxNumStat_Field_Num: 1,
// 		TNSFieldTotalAccNum: 1, TNSFieldDelegatorNum: 1,
// 		TNSFieldTokenStat: 1,
// 	}
// 	var txNumStatList []TxNumStat

// 	q := orm.NewQuery()
// 	q.SetCollection(CollectionTxNumStat)
// 	q.SetCondition(query)
// 	q.SetSelector(selector).SetSort(TxNumStat_Field_Date)
// 	q.SetResult(&txNumStatList)

// 	defer q.Release()

// 	err := q.Exec()
// 	return txNumStatList, err
// }

func (_ CommonTx) GetTxCountByDuration(startTime, endTime time.Time) (int, error) {

	db := orm.GetDatabase()
	defer db.Session.Close()

	txStore := db.C(CollectionNmCommonTx)

	query := bson.M{}
	query = FilterUnknownTxs(query)
	query["time"] = bson.M{"$gte": startTime, "$lt": endTime}

	return txStore.Find(query).Count()
}

// func (_ CommonTx) GetTxsByDurationAddress(startTime, endTime time.Time, address string) ([]CommonTx, error) {

// 	db := orm.GetDatabase()
// 	defer db.Session.Close()

// 	txStore := db.C(CollectionNmCommonTx)
// 	var res []CommonTx

// 	query := bson.M{}
// 	query = FilterUnknownTxs(query)
// 	query[Tx_Field_Time] = bson.M{"$gte": startTime, "$lt": endTime}
// 	query[Tx_Field_Status] = types.Success
// 	if address != "" {
// 		query["signers.addr_bech32"] = address

// 	}

// 	if err := txStore.Find(query).All(&res); err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

func (_ CommonTx) GetTxsByType(txtype string, status string) ([]CommonTx, error) {
	condition := bson.M{Tx_Field_Type: txtype, Tx_Field_Status: status}
	var txs []CommonTx
	err := queryAll(CollectionNmCommonTx, nil, condition, desc(Tx_Field_Time), 0, &txs)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

// func (_ CommonTx) QueryProposalTxFromById(idArr []uint64) (map[uint64]string, error) {

// 	selector := bson.M{Tx_Field_From: 1, Tx_Field_ProposalId: 1}
// 	condition := bson.M{Tx_Field_Type: "SubmitProposal", Tx_Field_Status: "success", Tx_Field_ProposalId: bson.M{"$in": idArr}}
// 	var txs []CommonTx
// 	condition = FilterUnknownTxs(condition)

// 	err := queryAll(CollectionNmCommonTx, selector, condition, desc(Tx_Field_Time), 0, &txs)

// 	proposerById := map[uint64]string{}

// 	for _, v := range txs {
// 		proposerById[v.ProposalId] = v.From
// 	}

// 	return proposerById, err
// }

func (_ CommonTx) QueryProposalTxListById(idArr []uint64) ([]CommonTx, error) {

	selector := bson.M{Tx_Field_Amount: 1, Tx_Field_ProposalId: 1}
	condition := bson.M{Tx_Field_Type: "SubmitProposal", Tx_Field_Status: "success", Tx_Field_ProposalId: bson.M{"$in": idArr}}
	var txs []CommonTx
	condition = FilterUnknownTxs(condition)

	err := queryAll(CollectionNmCommonTx, selector, condition, desc(Tx_Field_Time), 0, &txs)

	return txs, err
}

func (_ CommonTx) QueryProposalDeposit(before int, size int, id int) ([]CommonTx, error) {
	condition := bson.M{}
	if before != 0 {
		condition[Tx_Field_Height] = bson.M{
			"$lt": before,
		}
	}
	condition[Tx_Field_Event_Key] = "proposal_id"
	condition[Tx_Field_Event_Value] = strconv.Itoa(id)
	condition[Tx_Field_Event_Type] = TypeForDeposit
	var txs []CommonTx

	err := queryAll(CollectionNmCommonTx, nil, condition, desc(Tx_Field_Time), 0, &txs)

	return txs, err
}

func (_ CommonTx) QueryProposalVote(before int, size int, id int) ([]CommonTx, error) {
	condition := bson.M{}
	if before != 0 {
		condition[Tx_Field_Height] = bson.M{
			"$lt": before,
		}
	}
	condition[Tx_Field_Event_Key] = "proposal_id"
	condition[Tx_Field_Event_Value] = strconv.Itoa(id)
	condition[Tx_Field_Event_Type] = TypeForVote
	var txs []CommonTx

	err := queryAll(CollectionNmCommonTx, nil, condition, desc(Tx_Field_Time), 0, &txs)

	return txs, err
}

func (_ CommonTx) GetValueOfLog(logs []Log, eventType string, key string) (amount string) {
	log := logs[0]
	for _, event := range log.Events {
		if event.Type == eventType {
			fmt.Println(event)
			for _, attribute := range event.Attributes {
				if attribute.Key == key {
					return attribute.Value
				}
			}
		}
	}
	return amount
}

// func (_ CommonTx) QueryProposalTxById(proposalId int64, page, size int, total bool, iaaAddrs []string) (int, []CommonTx, error) {

// 	txs := []CommonTx{}

// 	selector := bson.M{
// 		Tx_Field_Height: 1,
// 		Tx_Field_Time:   1,
// 		Tx_Field_Hash:   1,
// 		Tx_Field_From:   1,
// 	}
// 	condition := bson.M{
// 		Tx_Field_Status:     types.TxTypeStatus,
// 		Tx_Field_ProposalId: proposalId,
// 		Tx_Field_Type:       types.TxTypeVote,
// 		Tx_Field_From:       bson.M{"$in": iaaAddrs},
// 	}
// 	condition = FilterUnknownTxs(condition)
// 	sort := fmt.Sprintf("-%v", Tx_Field_Height)

// 	num, err := pageQuery(CollectionNmCommonTx, selector, condition, sort, page, size, total, &txs)

// 	return num, txs, err
// }

// func (_ CommonTx) QueryDepositedProposalTxByValidatorWithSubmitOrDepositType(validatorAddrAcc string, page, size int, total bool) (int, []CommonTx, error) {

// 	txs := []CommonTx{}
// 	selector := bson.M{
// 		Tx_Field_Hash:       1,
// 		Tx_Field_From:       1,
// 		Tx_Field_Amount:     1,
// 		Tx_Field_Type:       1,
// 		Tx_Field_ProposalId: 1,
// 	}
// 	condition := bson.M{
// 		Tx_Field_Status: types.TxTypeStatus,
// 		Tx_Field_From:   validatorAddrAcc,
// 		Tx_Field_Type: bson.M{
// 			"$in": []string{types.TxTypeSubmitProposal, types.TxTypeDeposit},
// 		},
// 	}
// 	condition = FilterUnknownTxs(condition)
// 	sort := fmt.Sprintf("-%v", Tx_Field_Height)
// 	num, err := pageQuery(CollectionNmCommonTx, selector, condition, sort, page, size, total, &txs)

// 	return num, txs, err
// }

// func (_ CommonTx) QueryProposalTxByIdWithSubmitOrDepositType(proposalId int64, page, size int, total bool) (int, []CommonTx, error) {

// 	txs := []CommonTx{}
// 	selector := bson.M{
// 		Tx_Field_Hash:   1,
// 		Tx_Field_From:   1,
// 		Tx_Field_Amount: 1,
// 		Tx_Field_Type:   1,
// 		Tx_Field_Time:   1,
// 	}
// 	condition := bson.M{
// 		Tx_Field_Status:     types.TxTypeStatus,
// 		Tx_Field_ProposalId: proposalId,
// 		Tx_Field_Type: bson.M{
// 			"$in": []string{types.TxTypeSubmitProposal, types.TxTypeDeposit},
// 		},
// 	}
// 	condition = FilterUnknownTxs(condition)
// 	sort := fmt.Sprintf("-%v", Tx_Field_Height)
// 	num, err := pageQuery(CollectionNmCommonTx, selector, condition, sort, page, size, total, &txs)

// 	return num, txs, err
// }

// func (_ CommonTx) QueryTxAsset(assetType, tokenType, symbol, gateway string, page, size int, total bool) (int, []CommonTx, error) {
// 	txs := []CommonTx{}
// 	selector := bson.M{
// 		Tx_Field_Hash:      1,
// 		Tx_Field_Height:    1,
// 		Tx_Field_From:      1,
// 		Tx_Field_To:        1,
// 		Tx_Field_Amount:    1,
// 		Tx_Field_Type:      1,
// 		Tx_Field_Status:    1,
// 		Tx_Field_ActualFee: 1,
// 		Tx_Field_Tags:      1,
// 		Tx_Field_Msgs:      1,
// 		Tx_Field_Time:      1,
// 	}
// 	condition := bson.M{
// 		Tx_Field_Msgs_UdInfo: assetType,
// 	}
// 	condition = FilterUnknownTxs(condition)
// 	if tokenType != "" {
// 		condition[Tx_Field_Type] = tokenType
// 	} else {
// 		condition[Tx_Field_Type] = bson.M{
// 			"$in": []string{types.TxTypeIssueToken, types.TxTypeEditToken, types.TxTypeMintToken, types.TxTypeTransferTokenOwner},
// 		}
// 	}
// 	if symbol != "" {
// 		condition[Tx_Field_Msgs_UdInfo_Symbol] = bson.M{
// 			"$regex":   symbol,
// 			"$options": "$i",
// 		}
// 	}
// 	if gateway != "" {
// 		condition[Tx_Field_Msgs_UdInfo_Gateway] = bson.M{
// 			"$regex":   gateway,
// 			"$options": "$i",
// 		}
// 	}
// 	sort := fmt.Sprintf("-%v", Tx_Field_Height)
// 	num, err := pageQuery(CollectionNmCommonTx, selector, condition, sort, page, size, total, &txs)
// 	return num, txs, err
// }

func (_ CommonTx) QueryTxTransferGatewayOwner(moniker string, page, size int, total bool) (int, []CommonTx, error) {
	txs := []CommonTx{}
	selector := bson.M{
		Tx_Field_Hash:      1,
		Tx_Field_Height:    1,
		Tx_Field_From:      1,
		Tx_Field_To:        1,
		Tx_Field_Amount:    1,
		Tx_Field_Type:      1,
		Tx_Field_Status:    1,
		Tx_Field_ActualFee: 1,
		Tx_Field_Tags:      1,
		Tx_Field_Msgs:      1,
		Tx_Field_Time:      1,
	}

	condition := bson.M{
		Tx_Field_Type: Tx_Asset_TxType_TransferGatewayOwner,
	}
	condition = FilterUnknownTxs(condition)
	if moniker != "" {
		condition[Tx_Field_Msgs_Moniker] = moniker
	}

	sort := fmt.Sprintf("-%v", Tx_Field_Height)
	num, err := pageQuery(CollectionNmCommonTx, selector, condition, sort, page, size, total, &txs)
	return num, txs, err
}

func FilterUnknownTxs(query bson.M) bson.M {

	if status, ok := query["status"]; ok {
		query["status"] = status.(string)
	} else {
		query["status"] = bson.M{
			"$in": []string{TxStatusSuccess, TxStatusFail},
		}
	}
	return query
}

func (_ CommonTx) GetCountTxs(condition bson.M) (int, error) {
	result := []bson.M{}
	var query = orm.NewQuery()
	defer query.Release()
	query.SetResult(&result).
		SetCollection(CollectionNmCommonTx).
		PipeQuery(
			[]bson.M{
				{
					"$match": condition,
				},
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
