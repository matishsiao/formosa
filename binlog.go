package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Binlog struct {
	db  *leveldb.DB
	seq uint64
}

func (bin *Binlog) Construct() {
	db, err := leveldb.OpenFile("meta", nil)
	if err != nil {
		log.Fatal("Open DB failed:", err)
	}
	bin.db = db
	binlog, err := bin.Get("BINLOG_SEQ")
	if err != nil {
		log.Println("Can't find binlog serial number", err)
		bin.seq = 0
		err = bin.Commit()
		if err != nil {
			log.Println("binlog serial number init failed.", err)
		}
	} else {
		if num, err := strconv.ParseUint(string(binlog), 10, 64); err == nil {
			bin.seq = num
		}
	}
	log.Println("binlog serial number", bin.seq)
}

func (bin *Binlog) Close() {
	if bin.db != nil {
		err := bin.db.Close()
		log.Fatal("Close DB failed:", err)
	}
}

func (bin *Binlog) Push(args []string) {
	meta, err := json.Marshal(args)
	if err != nil {
		log.Println("Set Meta Failed:", err)
	} else {
		bin.Set(meta)
	}
}

func (bin *Binlog) Set(value []byte) error {
	binKey := EncodeBinlog(fmt.Sprintf("%d", bin.seq))
	bin.seq++
	bin.Commit()
	return bin.db.Put([]byte(binKey), value, nil)
}

func (bin *Binlog) Scan(from string, to string) ([]string, error) {
	var scanRange *util.Range
	fromKey := EncodeBinlog(from)
	toKey := EncodeBinlog(to)
	if from != "" && to != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: []byte(toKey)}
	} else if from != "" {
		scanRange = &util.Range{Start: []byte(fromKey), Limit: nil}
	} else if to != "" {
		scanRange = &util.Range{Start: nil, Limit: []byte(toKey)}
	} else {
		scanRange = nil
	}

	iter := bin.db.NewIterator(scanRange, nil)
	var result []string
	for iter.Next() {
		result = append(result, string(iter.Key()))
		result = append(result, string(iter.Value()))
	}
	iter.Release()
	err := iter.Error()
	return result, err
}

func (bin *Binlog) Get(key string) ([]byte, error) {
	binKey := EncodeBinlog(key)
	return bin.db.Get([]byte(binKey), nil)
}

func (bin *Binlog) Commit() error {
	idx := fmt.Sprintf("%d", bin.seq)
	return bin.db.Put([]byte("BINLOG_SEQ"), []byte(idx), nil)
}

func (bin *Binlog) Exists(key string) (bool, error) {
	binKey := EncodeBinlog(key)
	return bin.db.Has([]byte(binKey), nil)
}
