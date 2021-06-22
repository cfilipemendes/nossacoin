package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func bootstrapServer() error {
	mux := mux.NewRouter()

	mux.HandleFunc("/", getNossaChain).Methods("GET")
	mux.HandleFunc("/", writeNewBlock).Methods("POST")

	httpPort := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println("Listening on port ", httpPort)

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

	newBlock, err := addNewBlock(m.Data)

	if err != nil {
		respondWithJson(w, r, http.StatusInternalServerError, m)
		return
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
