package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Blockd struct { //Block structure
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

type Message struct {
	BPM int //This is the payload i.e the values we can give in the post request in JSON format
}

var Blockchain []Blockd

// Generating a new block using the has from calculateHash function and current time and index incremented from previous block

func GenerateBlock(index int, hash string, BPM int) (Blockd, error) {
	var newBlock Blockd
	var oldIndex = index
	var oldHash = hash

	t := time.Now()

	newBlock.Index = oldIndex + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldHash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock, nil
}

//takes our block data and creates a SHA256 hash of it
func CalculateHash(block Blockd) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	println(h)
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
