package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	// TimeStamp time when this block is created
	TimeStamp int64
	// PrevBlockHash link this block with previous block in chain by a hash
	PrevBlockHash []byte
	// Hash contain hash of this block, acquire by mining
	Hash []byte
	// Nonce is used to verify the block
	Nonce int

	// Data contains block data
	Data []byte
}

// Serialize block can be serialize into bytes for storing
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	_ = encoder.Encode(b)

	return result.Bytes()
}

func DeserializeBlock(b []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(b))
	_ = decoder.Decode(&block)

	return &block

}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := &Block{
		TimeStamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Nonce:         0,
	}

	pow := NewProofOfWork(b)
	nonce, hash := pow.Run()

	b.Hash = hash
	b.Nonce = nonce

	return b
}
