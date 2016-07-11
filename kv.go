package formosa

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/util"
)

//KV
func (dm *DBManager) Set(key string, value string) error {
	enKey := EncodeKV(key)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	if !exists {
		data, _ := dm.DB.Get([]byte(DATATYPE_KV_END), nil)
		if len(data) == 0 {
			dm.DB.Put([]byte(DATATYPE_KV_END), []byte("1"), nil)
		} else {
			size := ToInt64(string(data))
			size++
			dm.DB.Put([]byte(DATATYPE_KV_END), []byte(fmt.Sprintf("%d", size)), nil)
		}
	}
	return dm.DB.Put([]byte(enKey), []byte(value), nil)
}

func (dm *DBManager) Del(key string) error {
	enKey := EncodeKV(key)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	if !exists {
		return nil
	}
	data, _ := dm.DB.Get([]byte(DATATYPE_KV_END), nil)
	if len(data) == 0 {
		dm.DB.Put([]byte(DATATYPE_KV_END), []byte("0"), nil)
	} else {
		size := ToInt64(string(data))
		size--
		if size <= 0 {
			size = 0
		}
		dm.DB.Put([]byte(DATATYPE_KV_END), []byte(fmt.Sprintf("%d", size)), nil)
	}
	return dm.DB.Delete([]byte(enKey), nil)
}

func (dm *DBManager) Get(key string) ([]byte, error) {
	enKey := EncodeKV(key)
	return dm.DB.Get([]byte(enKey), nil)
}

func (dm *DBManager) Exists(key string) (bool, error) {
	enKey := EncodeKV(key)
	return dm.DB.Has([]byte(enKey), nil)
}

func (dm *DBManager) Incr(key string, value string) (string, error) {
	enKey := EncodeKV(key)
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return "", err
	}
	exists, _ := transaction.Has([]byte(enKey), nil)
	var response string
	var rep_err error
	if !exists {
		data, _ := transaction.Get([]byte(DATATYPE_KV_END), nil)
		if len(data) == 0 {
			transaction.Put([]byte(DATATYPE_KV_END), []byte("1"), nil)
		} else {
			size := ToInt64(string(data))
			size++
			transaction.Put([]byte(DATATYPE_KV_END), []byte(fmt.Sprintf("%d", size)), nil)
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

func (dm *DBManager) Scan(from string, to string, limit int64) ([]string, error) {
	var scanRange *util.Range
	var result []string
	if limit == 0 {
		return result, fmt.Errorf("limit can't set zero")
	}
	fromKey := EncodeKV(from)
	toKey := EncodeKV(to)
	if to != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(toKey)}
	} else {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(DATATYPE_KV_END)}
	}

	iter := dm.DB.NewIterator(scanRange, nil)
	for iter.Next() {
		key := DecodeKV(string(iter.Key()))
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

func (dm *DBManager) Size() (int64, error) {
	size, err := dm.DB.Get([]byte(DATATYPE_KV_END), nil)
	return ToInt64(string(size)), err
}
