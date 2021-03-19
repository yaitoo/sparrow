package mgo

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FieldSort int
type FieldSort int

const (
	//FieldAsc Asc
	FieldAsc FieldSort = 1
	//FieldDesc Desc
	FieldDesc FieldSort = -1
)

// FieldDesribe field of index
type FieldDesribe struct {
	FieldName string
	Sort      FieldSort
}

// IdxKey IdxKey
type IdxKey struct {
	Keys []FieldDesribe
	Opt  options.IndexOptions
}

const indexPrefix = "PARA-INDEX_V"
const delimiter = "_"

func genIndexName(version int, baseName string) string {
	var s string
	s = indexPrefix
	s += strconv.Itoa(version)
	s += delimiter
	s += baseName
	return s
}

func getVersionByIdxName(idxName string) (int, bool) {

	if !strings.HasPrefix(idxName, indexPrefix) {
		return 0, false
	}

	su := strings.TrimLeft(idxName, indexPrefix)
	vers := strings.Split(su, delimiter)
	version, _ := strconv.Atoi(vers[0])

	return version, true
}

func genIndexModelByIdxKeys(keys *IdxKey) mongo.IndexModel {
	var ks []primitive.E

	for _, v := range keys.Keys {
		e := primitive.E{v.FieldName, v.Sort}
		ks = append(ks, e)
	}
	return mongo.IndexModel{
		Keys:    ks,
		Options: &keys.Opt,
	}

}

func validateIdxParam(idxs []IdxKey) error {
	if idxs == nil {
		return errors.New("idxs is nil")
	}

	for i, v := range idxs {
		var err string
		keys := v.Keys
		if keys == nil {
			err = "Keys[" + strconv.Itoa(i) + "] is nil"
			return errors.New(err)
		}

		if len(keys) == 0 {
			err = "len of Keys[" + strconv.Itoa(i) + "] is 0"
			return errors.New(err)
		}

		opt := v.Opt
		if opt.Name == nil {
			err = "indexName of option[" + strconv.Itoa(i) + "] is nil"
			return errors.New(err)
		}
	}
	return nil
}

type index struct {
	Key  map[string]int
	NS   string
	Name string
}

func deleteOldIdx(ctx context.Context, collection *mongo.Collection, version int) error {
	indexView := collection.Indexes()
	cursor, err := indexView.List(ctx)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		var idx index
		err = cursor.Decode(&idx)
		if err != nil {
			return err
		}

		oldVer, needMig := getVersionByIdxName(idx.Name)
		if needMig && oldVer < version {
			_, err := indexView.DropOne(ctx, idx.Name)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

func createIdx(ctx context.Context, collection *mongo.Collection, keys *IdxKey, version int) (string, error) {
	indexView := collection.Indexes()
	indexName := genIndexName(version, *keys.Opt.Name)
	keys.Opt.SetName(indexName)
	im := genIndexModelByIdxKeys(keys)
	return indexView.CreateOne(ctx, im)
}

// MigrateIndex create new version index and delete old version index
func MigrateIndex(ctx context.Context, collection *mongo.Collection, version int, idxs []IdxKey) ([]string, error) {

	err := validateIdxParam(idxs)
	if err != nil {
		return nil, err
	}

	// 1. Delete old index
	err = deleteOldIdx(ctx, collection, version)
	if err != nil {
		return nil, err
	}

	// 2. create new index
	var idxName []string
	for _, v := range idxs {
		name, err := createIdx(ctx, collection, &v, version)
		if err != nil {
			return idxName, err
		}
		idxName = append(idxName, name)
	}

	return idxName, nil
}
