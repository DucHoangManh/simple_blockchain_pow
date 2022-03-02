package blockchain

import "go.etcd.io/bbolt"

// ChainIterator use to iterator over blocks
type ChainIterator struct {
	currentHash []byte
	db          *bbolt.DB
}

func (i *ChainIterator) Next() *Block {
	var block *Block
	_ = i.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		serializedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(serializedBlock)

		return nil
	})
	i.currentHash = block.PrevBlockHash
	return block
}
