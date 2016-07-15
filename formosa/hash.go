package formosa

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/util"
)

func (dm *DBManager) HashSet(hash string, key string, value string) (bool, error) {
	enKey := EncodeHash(hash, key)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return false, err
	}
	commit := false
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
		commit = true
	} else {
		data, _ := dm.HashGet(hash, key)
		if string(data) != value {
			commit = true
		}
	}

	if commit {
		transaction.Put([]byte(enKey), []byte(value), nil)
		err = transaction.Commit()
		dm.HashNameSet(hash)
		return true, err
	} else {
		transaction.Discard()
		return false, nil
	}
}

func (dm *DBManager) HashNameSet(hash string) {
	enKey := EncodeHashName(hash)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	if !exists {
		dm.DB.Put([]byte(enKey), []byte("1"), nil)
	}
}

func (dm *DBManager) HashDel(hash string, key string) (bool, error) {
	enKey := EncodeHash(hash, key)
	exists, _ := dm.DB.Has([]byte(enKey), nil)
	if !exists {
		return false, nil
	}
	endKey := EncodeHashEnd(hash)
	transaction, err := dm.DB.OpenTransaction()
	if err != nil {
		return false, err
	}
	data, _ := transaction.Get([]byte(endKey), nil)
	if len(data) == 0 {
		transaction.Put([]byte(endKey), []byte("0"), nil)
	} else {
		size := ToInt64(string(data))
		size--
		if size <= 0 {
			size = 0
		}
		transaction.Put([]byte(endKey), []byte(fmt.Sprintf("%d", size)), nil)
	}
	transaction.Delete([]byte(enKey), nil)
	err = transaction.Commit()
	return true, err
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

func (dm *DBManager) HashReverseScan(hash string, from string, to string, limit int64) ([]string, error) {
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
	if iter.Last() {
		key := DecodeHash(string(iter.Key()), hash)
		result = append(result, key)
		result = append(result, string(iter.Value()))
		if limit != -1 {
			limit--
		}
		if limit != 0 {
			for iter.Prev() {
				key := DecodeHash(string(iter.Key()), hash)
				result = append(result, key)
				result = append(result, string(iter.Value()))
				if limit != -1 {
					limit--
				}
				if limit == 0 {
					break
				}
			}
		}

	}
	iter.Release()
	err := iter.Error()
	return result, err
}

func (dm *DBManager) HashList(from string, to string, limit int64) ([]string, error) {
	var scanRange *util.Range
	var result []string
	if limit == 0 {
		return result, fmt.Errorf("limit can't set zero")
	}
	fromKey := EncodeHashName(from)
	if to != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(EncodeHashName(to))}
	} else {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(DATATYPE_HASH_LIST_END)}
	}

	iter := dm.DB.NewIterator(scanRange, nil)

	for iter.Next() {
		key := DecodeHashName(string(iter.Key()))
		result = append(result, key)
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

func (dm *DBManager) HashReverseList(from string, to string, limit int64) ([]string, error) {
	var scanRange *util.Range
	var result []string
	if limit == 0 {
		return result, fmt.Errorf("limit can't set zero")
	}
	fromKey := EncodeHashName(from)
	if to != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(EncodeHashName(to))}
	} else {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(DATATYPE_HASH_LIST_END)}
	}

	iter := dm.DB.NewIterator(scanRange, nil)
	if iter.Last() {
		key := DecodeHashName(string(iter.Key()))
		result = append(result, key)
		if limit != -1 {
			limit--
		}
		if limit != 0 {
			for iter.Prev() {
				key := DecodeHashName(string(iter.Key()))
				result = append(result, key)
				if limit != -1 {
					limit--
				}
				if limit == 0 {
					break
				}
			}
		}
	}

	iter.Release()
	err := iter.Error()
	return result, err
}
