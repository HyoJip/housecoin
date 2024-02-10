package db

import (
	"github.com/Hyojip/housecoin/utils"
	"github.com/boltdb/bolt"
)

const (
	DbName        = "blockchain.db"
	BucketData    = "data"
	BucketBlocks  = "blocks"
	KeyCheckpoint = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db != nil {
		return db
	}

	dbPointer, err := bolt.Open(DbName, 0600, nil)
	utils.HandleError(err)
	db = dbPointer
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketData))
		utils.HandleError(err)
		_, err = tx.CreateBucketIfNotExists([]byte(BucketBlocks))
		utils.HandleError(err)
		return err
	})
	return db
}

func Close() {
	db.Close()
}

func SaveBlock(hash string, data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketBlocks))
		return bucket.Put([]byte(hash), data)
	})
	utils.HandleError(err)
}

func SaveCheckpoint(data []byte) {
	err := DB().Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketData))
		return bucket.Put([]byte(KeyCheckpoint), data)
	})
	utils.HandleError(err)
}

func Checkpoint() []byte {
	var blockchainInDB []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketData))
		blockchainInDB = bucket.Get([]byte(KeyCheckpoint))
		return nil
	})
	return blockchainInDB
}

func Block(hash string) []byte {
	var block []byte
	DB().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketBlocks))
		block = bucket.Get([]byte(hash))
		return nil
	})
	return block
}
