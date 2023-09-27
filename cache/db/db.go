package db

import (
	"SQLGuardian/consts"
	"fmt"
	"github.com/dgraph-io/badger/v4"
)

// GetBadgerDBInMemory 获取BadgerDB实例
func GetBadgerDBInMemory() *badger.DB {
	// 打开或创建一个Badger数据库
	opts := badger.DefaultOptions(consts.DBPath).WithInMemory(false)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		fmt.Printf("badger.Open() error: %v\n", err)
	}
	return db
}

// Get key from BadgerDB
func Get(key []byte) (value []byte, err error) {
	db := GetBadgerDBInMemory()
	defer func(db *badger.DB) {
		err = db.Close()
		if err != nil {
			fmt.Printf("db.Close() error: %v\n", err)
		}
	}(db)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return nil
		}
		err = item.Value(func(val []byte) error {
			value = append([]byte{}, val...)
			return nil
		})
		return err
	})
	return
}

// Set key-value to BadgerDB
func Set(key, value []byte) error {
	db := GetBadgerDBInMemory()
	defer func(db *badger.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("db.Close() error: %v\n", err)
		}
	}(db)
	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})
	return err
}
