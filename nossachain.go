package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
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

func chainInit() {
	genesis := Block{0, time.Now().String(), "Genesis.", "", ""}
	NossaChain = append(NossaChain, genesis)
	if isDev() {
		spew.Dump(genesis)
	}
}

func addNewBlock(data string) (Block, error) {
	newBlock, err := createBlock(NossaChain[len(NossaChain)-1], data)

	if isBlockValid(newBlock, NossaChain[len(NossaChain)-1]) {
		updatedNossaChain := append(NossaChain, newBlock)
		replaceChain(updatedNossaChain)

		if isDev() {
			spew.Dump(NossaChain)
		}
	}

	return newBlock, err
}

func createBlock(oldBlock Block, data string) (Block, error) {

	newBlock := Block{
		oldBlock.Index + 1,
		time.Now().String(),
		data,
		oldBlock.Hash,
		""}
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
