package main

import "github.com/syndtr/goleveldb/leveldb"

func (dm *DBManager) Batch(batchlist [][]string) error {
	batch := new(leveldb.Batch)
	for _, v := range batchlist {
		key := dm.BatchKey(v)
		if key != "" {
			value := v[len(v)-1]
			batch.Put([]byte(key), []byte(value))
		}
	}
	err := dm.DB.Write(batch, nil)
	return err
}

func (dm *DBManager) BatchKey(args []string) string {
	key := ""
	switch args[0] {
	case "hset":
		if len(args) == 4 {
			key = EncodeHash(args[1], args[2])
		}
	case "set":
		if len(args) == 3 {
			key = EncodeKV(args[1])
		}
	}
	return key
}
