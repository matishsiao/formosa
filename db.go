package main

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DBManager struct {
	db *leveldb.DB
}

func (dm *DBManager) Construct() {
	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
		log.Fatal("Open DB failed:", err)
	}
	dm.db = db
}

func (dm *DBManager) Close() {
	if dm.db != nil {
		err := dm.db.Close()
		log.Fatal("Close DB failed:", err)
	}
}

func (dm *DBManager) Set(key string, value string) error {
	enKey := EncodeKV(key)
	return dm.db.Put([]byte(enKey), []byte(value), nil)
}

func (dm *DBManager) Get(key string) ([]byte, error) {
	enKey := EncodeKV(key)
	return dm.db.Get([]byte(enKey), nil)
}

func (dm *DBManager) Exists(key string) (bool, error) {
	enKey := EncodeKV(key)
	return dm.db.Has([]byte(enKey), nil)
}

func (dm *DBManager) Scan(from string, to string) ([]string, error) {
	var scanRange *util.Range
	fromKey := EncodeKV(from)
	toKey := EncodeKV(to)
	if from != "" && to != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(toKey)}
	} else if from != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: nil}
	} else if to != "" {
		scanRange = &util.Range{Start: nil, Limit: []byte(toKey)}
	} else {
		scanRange = nil
	}

	iter := dm.db.NewIterator(scanRange, nil)
	var result []string
	for iter.Next() {
		result = append(result, DecodeKV(string(iter.Key())))
		result = append(result, string(iter.Value()))
	}
	iter.Release()
	err := iter.Error()
	return result, err
}
