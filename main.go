package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

//CHAIN
var NossaChain []Block

type Block struct {
	Index     int
	TimeStamp string
	Data      decimal.Decimal //To Represent decimal values where each right value must not be rounded.
	Hash      string
	PrevHash  string
}

type Message struct {
	Data decimal.Decimal
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())

}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.TimeStamp + block.Data.String() + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)

}

func generateBlock(oldBlock Block, data decimal.Decimal) (Block, error) {
	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.TimeStamp = t.String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock, nil
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

func run() error {
	mux := mux.NewRouter()

	mux.HandleFunc("/", getNossaChain).Methods("GET")
	mux.HandleFunc("/", writeNewBlock).Methods("POST")

	httpPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println("Listening on port ", os.Getenv("PORT"))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(httpPort, mux))

	return nil

}

func getNossaChain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(NossaChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func writeNewBlock(w http.ResponseWriter, r *http.Request) {
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJson(w, r, http.StatusBadRequest, m)
		return
	}
	defer r.Body.Close()

	currLastBlock := NossaChain[len(NossaChain)-1]
	newBlock, err := generateBlock(currLastBlock, m.Data)
	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, m)
		return
	}

	if isBlockValid(newBlock, currLastBlock) {
		updatedNossaChain := append(NossaChain, newBlock)
		replaceChain(updatedNossaChain)
		spew.Dump(NossaChain)
	}

	respondWithJson(w, r, http.StatusCreated, newBlock)

}

func respondWithJson(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	//TODO
}
