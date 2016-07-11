package formosa

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type DBManager struct {
	DB *leveldb.DB
}

func (dm *DBManager) Construct(dir string) {
	option := &opt.Options{}
	option.BlockSize = 64 * opt.KiB
	option.WriteBuffer = 64 * opt.MiB
	option.CompactionTableSize = 1000 * opt.MiB
	option.BlockCacheCapacity = 500 * opt.MiB
	DB, err := leveldb.OpenFile(dir, option)
	if err != nil {
		log.Fatal("Open DB failed:", err)
	}
	dm.DB = DB
}

func (dm *DBManager) Close() {
	if dm.DB != nil {
		err := dm.DB.Close()
		log.Fatal("Close DB failed:", err)
	}
}

//Global
func (dm *DBManager) GlobalGetAll() ([]string, error) {
	scanRange := &util.Range{Start: nil, Limit: nil}
	iter := dm.DB.NewIterator(scanRange, nil)
	var result []string
	for iter.Next() {
		result = append(result, string(iter.Key()))
		result = append(result, string(iter.Value()))
	}
	iter.Release()
	err := iter.Error()
	return result, err
}
