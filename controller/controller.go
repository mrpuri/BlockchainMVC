package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	http "net/http"
	"project/connection"
	"project/model"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

//handlers

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/fetch", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/enter", handleWriteBlock).Methods("POST")
	return muxRouter

}

//get handler

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(model.Blockchain, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("error 1")
		return
	}
	io.WriteString(w, string(bytes))
}

//function for post request

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m model.Message //contains the user given value in post request

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	var (
		index int
		hash  string
		phash string
	)
	var i uint8
	var s string
	stmt, err := connection.Db.Prepare("SELECT 'index' FROM Block WHERE phash = ?")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(s).Scan(&i)

	if err != nil {
		log.Fatal(err)
	}
	rows, err := connection.Db.Query("SELECT hash,phash,'index' FROM block WHERE index = ?", i)
	if err != nil {

		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&hash, &phash, &index)
		if err != nil {
			log.Fatal(err)
		}

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	newBlock, err := model.GenerateBlock(index, phash, m.BPM) //fetching previous block and calling function to generate new block
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if isBlockValid(newBlock, model.Blockchain[len(model.Blockchain)-1]) {
		newBlockchain := append(model.Blockchain, newBlock)
		replaceChain(newBlockchain)
		spew.Dump(model.Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// To validate the block i.e to check if it is tampered with
func isBlockValid(newBlock, oldBlock model.Blockd) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if model.CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

//keep the longest chain in case two blocks are being added simultaneously

func replaceChain(newBlocks []model.Blockd) {
	if len(newBlocks) > len(model.Blockchain) {
		model.Blockchain = newBlocks
	}
}
