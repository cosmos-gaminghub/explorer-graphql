package document

import (
	"fmt"
	"log"

	"graphql-go-example/orm"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func desc(field string) string {
	return fmt.Sprintf("-%s", field)
}

func queryAll(collation string, selector, condition bson.M, sort string, size int, result interface{}) error {
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(collation).
		SetCondition(condition).
		SetSelector(selector).
		SetSort(sort).
		SetSize(size).
		SetResult(result)

	err := query.Exec()
	if err != nil {
		log.Fatal("queryAll error")
	}
	return err
}

func queryOne(collation string, selector, condition bson.M, result interface{}) error {
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(collation).
		SetCondition(condition).
		SetSelector(selector).
		SetResult(result)

	err := query.Exec()
	if err != nil {
		return err
	}
	return nil
}

func querylistByOffsetAndSize(collection string, selector, condition bson.M, sort string, offset, size int, result interface{}) error {

	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(collection).
		SetCondition(condition).
		SetSelector(selector).
		SetSort(sort).
		SetOffset(offset).
		SetSize(size).
		SetResult(result)

	err := query.Exec()
	if err != nil {
		return err
	}
	return nil
}

func pageQuery(collation string, selector, condition bson.M, sort string, page, size int, total bool, result interface{}) (int, error) {
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(collation).
		SetCondition(condition).
		SetSelector(selector).
		SetSort(sort).
		SetPage(page).
		SetSize(size).
		SetResult(result)

	cnt, err := query.ExecPage(total)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func getDb() *mgo.Database {
	return orm.GetDatabase()
}

func asc(field string) string {
	return fmt.Sprintf("%s", field)
}
