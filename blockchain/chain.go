package blockchain

import (
	"go.etcd.io/bbolt"
	"log"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const lastBlockKey = "l"

type BlockChain struct {
	// tip store the latest block hash
	tip []byte
	// db connection instance
	db *bbolt.DB
}

func NewBlockchain() (bc *BlockChain, err error) {
	var tip []byte
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Printf("Error when open db for block persistence %s", err)
		return
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		// get the block bucket, if not exist, create a new blockchain with genesis block
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			genesisBlock := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				return err
			}
			err = b.Put([]byte(lastBlockKey), genesisBlock.Hash)
			if err != nil {
				return err
			}
			tip = genesisBlock.Hash
		} else {
			// else get the tip of the blockchain from db
			tip = b.Get([]byte(lastBlockKey))
		}
		return nil
	})
	if err != nil {
		log.Printf("Error while persisting block")
	}
	return &BlockChain{tip, db}, nil
}

// NewGenesisBlock genesis block is the first block in the blockchain
func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

func (bc *BlockChain) AddBlock(data string) (err error) {
	if bc.db == nil {
		err = bc.InitDBConn()
		if err != nil {
			return err
		}
	}
	var lastHash []byte
	err = bc.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Printf("Error when get last block hash from db %s", err)
		return
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err = b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		err = b.Put([]byte(lastBlockKey), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})

	return err
}

func (bc *BlockChain) Iterator() *ChainIterator {
	return &ChainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

func (bc *BlockChain) InitDBConn() (err error) {
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Printf("Error when open db connection %s", err)
		return
	}
	bc.db = db
	return
}

func (bc *BlockChain) CloseDbConn() error {
	return bc.db.Close()
}