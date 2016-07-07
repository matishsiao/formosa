package main

import (
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func (dm *DBManager) HashSet(hash string, key string, value string) error {
	enKey := EncodeHash(hash, key)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	if !exists {
		endKey := EncodeHashEnd(hash)
		data, _ := dm.DB.Get([]byte(endKey), nil)
		if len(data) == 0 {
			dm.DB.Put([]byte(endKey), []byte("1"), nil)
		} else {
			size := ToInt64(string(data))
			size++
			dm.DB.Put([]byte(endKey), []byte(fmt.Sprintf("%d", size)), nil)
		}
	}
	return dm.DB.Put([]byte(enKey), []byte(value), nil)
}

func (dm *DBManager) HashDel(hash string, key string) error {
	enKey := EncodeHash(hash, key)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	if !exists {
		return nil
	}
	endKey := EncodeHashEnd(hash)
	data, _ := dm.DB.Get([]byte(endKey), nil)
	if len(data) == 0 {
		dm.DB.Put([]byte(endKey), []byte("0"), nil)
	} else {
		size := ToInt64(string(data))
		size--
		if size <= 0 {
			size = 0
		}
		dm.DB.Put([]byte(endKey), []byte(fmt.Sprintf("%d", size)), nil)
	}
	return dm.DB.Delete([]byte(enKey), nil)
}

func (dm *DBManager) HashIncr(hash string, key string, value string) (string, error) {
	enKey := EncodeHash(hash, key)
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return "", err
	}
	exists, _ := transaction.Has([]byte(enKey), nil)
	var response string
	var rep_err error
	if !exists {
		endKey := EncodeHashEnd(hash)
		data, _ := transaction.Get([]byte(endKey), nil)
		if len(data) == 0 {
			transaction.Put([]byte(endKey), []byte("1"), nil)
		} else {
			size := ToInt64(string(data))
			size++
			transaction.Put([]byte(endKey), []byte(fmt.Sprintf("%d", size)), nil)
		}
		rep_err = transaction.Put([]byte(enKey), []byte(value), nil)
		response = value
	} else {
		dbValue, err := transaction.Get([]byte(enKey), nil)
		if err != nil {
			rep_err = err
		} else {
			dbIncr := ToInt64(string(dbValue))
			dbIncr += ToInt64(string(value))
			rep_err = transaction.Put([]byte(enKey), []byte(fmt.Sprintf("%d", dbIncr)), nil)
			response = fmt.Sprintf("%d", dbIncr)
		}
	}
	rep_err = transaction.Commit()
	return response, rep_err
}

func (dm *DBManager) HashGet(hash string, key string) ([]byte, error) {
	enKey := EncodeHash(hash, key)
	return dm.DB.Get([]byte(enKey), nil)
}

func (dm *DBManager) HashExists(hash string, key string) (bool, error) {
	enKey := EncodeHash(hash, key)
	return dm.DB.Has([]byte(enKey), nil)
}

func (dm *DBManager) HashSize(hash string) (int64, error) {
	endKey := EncodeHashEnd(hash)
	size, err := dm.DB.Get([]byte(endKey), nil)
	return ToInt64(string(size)), err
}

func (dm *DBManager) HashScan(hash string, from string, to string, limit int64) ([]string, error) {
	var scanRange *util.Range
	var result []string
	if limit == 0 {
		return result, fmt.Errorf("limit can't set zero")
	}
	fromKey := EncodeHash(hash, from)
	endKey := EncodeHashEnd(hash)
	toKey := EncodeHash(hash, to)
	if to != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(toKey)}
	} else {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(endKey)}
	}

	iter := dm.DB.NewIterator(scanRange, nil)

	for iter.Next() {
		key := DecodeHash(string(iter.Key()), hash)
		log.Println(string(iter.Key()), string(iter.Value()))
		result = append(result, key)
		result = append(result, string(iter.Value()))
		if limit != -1 {
			limit--
		}
		if limit == 0 {
			break
		}
	}
	iter.Release()
	err := iter.Error()
	return result, err
}
