package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

//CHAIN
var NossaChain []Block

type Block struct {
	Index     int
	TimeStamp string
	Data      string
	Hash      string
	PrevHash  string
}

type Message struct {
	Data string
}

func addBlock(oldBlock Block, data string) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = t.String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
}
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.TimeStamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(NossaChain) {
		NossaChain = newBlocks
	}
}
