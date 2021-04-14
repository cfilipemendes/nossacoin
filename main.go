package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

//CHAIN

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Initialize Chain.
	go func() {
		genesis := Block{0, time.Now().String(), "Genesis.", "", ""}
		spew.Dump(genesis)
		NossaChain = append(NossaChain, genesis)
	}()

	log.Fatal(run())

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
	newBlock, err := addBlock(currLastBlock, m.Data)
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
	res, err := json.MarshalIndent(payload, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500 - Internal Server Error"))
		return
	}
	w.WriteHeader(status)
	w.Write(res)
}
