package main

import (
	"encoding/binary"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type DBManager struct {
	LocalDB *leveldb.DB
	MetaDB  *leveldb.DB
	MetaIdx uint64
}

func (dm *DBManager) Construct() {
	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
		log.Fatal("Open DB failed:", err)
	}
	dm.LocalDB = db

	db, err = leveldb.OpenFile("meta", nil)
	if err != nil {
		log.Fatal("Open DB failed:", err)
	}
	dm.MetaDB = db
	binlog, err := dm.MetaGet([]byte("BINLOG_SER"))
	if err != nil {
		log.Println("Can't find binlog serial number", err)
		dm.MetaIdx = 0
		err = dm.MetaSerialSet()
		if err != nil {
			log.Println("binlog serial number init failed.", err)
		}
	} else {
		num := binary.LittleEndian.Uint64(binlog)
		dm.MetaIdx = num
	}
	log.Println("binlog serial number", dm.MetaIdx)
}

func (dm *DBManager) Close() {
	if dm.LocalDB != nil {
		err := dm.LocalDB.Close()
		log.Fatal("Close DB failed:", err)
	}

	if dm.MetaDB != nil {
		err := dm.MetaDB.Close()
		log.Fatal("Close DB failed:", err)
	}
}

func (dm *DBManager) Set(key []byte, value []byte) error {
	return dm.LocalDB.Put(key, value, nil)
}

func (dm *DBManager) Get(key []byte) ([]byte, error) {
	return dm.LocalDB.Get(key, nil)
}

func (dm *DBManager) Exists(key []byte) (bool, error) {
	return dm.LocalDB.Has(key, nil)
}

func (dm *DBManager) MetaSet(value []byte) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, dm.MetaIdx)
	return dm.MetaDB.Put(b, value, nil)
}

func (dm *DBManager) MetaGet(key []byte) ([]byte, error) {
	return dm.MetaDB.Get(key, nil)
}

func (dm *DBManager) MetaSerialSet() error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, dm.MetaIdx)
	return dm.MetaDB.Put([]byte("BINLOG_SER"), b, nil)
}

func (dm *DBManager) MetaExists(key []byte) (bool, error) {
	return dm.MetaDB.Has(key, nil)
}
