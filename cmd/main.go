package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type transaction struct {
	sender         string
	receiver       string
	amount         float64
	transactionFee float64
	signature      string
	timestamp      string
}

type block struct {
	index        int
	timestamp    string
	transaction  []transaction
	previousHash string
	hash         string
}

var blockchain []block

// Create the genesis block (first block in the blockchain)
func createGenesisBlock() {
	genesisBlock := block{
		index:        0,
		timestamp:    time.Now().String(),
		transaction:  []transaction{},
		previousHash: "",
	}
	genesisBlock.hash = calculateHash(&genesisBlock)

	blockchain = append(blockchain, genesisBlock)
}

// Calculate the SHA-256 hash of a block
func calculateHash(block *block) string {
	record := fmt.Sprintf("%d%s%v%s", block.index, block.timestamp, block.transaction, block.previousHash)
	hashInBytes := sha256.Sum256([]byte(record))

	return hex.EncodeToString(hashInBytes[:])
}

// Create a new transaction
func createTransaction() transaction {
	return transaction{
		sender:         "Alice",
		receiver:       "Bob",
		amount:         10.0,
		transactionFee: 0.2,
		signature:      "abc123", // Placeholder for the transaction signature
		timestamp:      time.Now().String(),
	}
}

// Create a new block with the given data
func generateBlock(previousBlock *block, transactions *[]transaction) block {
	newBlock := block{
		index:        previousBlock.index + 1,
		timestamp:    time.Now().String(),
		transaction:  *transactions,
		previousHash: previousBlock.hash,
	}
	newBlock.hash = calculateHash(&newBlock)

	return newBlock
}

// Check if a block is valid by verifying its hash and index
func isBlockValid(newBlock *block, previousBlock *block) bool {
	if previousBlock.index+1 != newBlock.index || previousBlock.hash != newBlock.previousHash || calculateHash(newBlock) != newBlock.hash {
		return false
	}

	return true
}

func printBlockchain() {
	for _, block := range blockchain {
		fmt.Printf("Index: %d\n", block.index)
		fmt.Printf("Timestamp: %s\n", block.timestamp)
		fmt.Printf("Transactions: %v\n", block.transaction)
		fmt.Printf("Previous Hash: %s\n", block.previousHash)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Println()
	}
}

func main() {
	createGenesisBlock()

	newTransaction := createTransaction()

	newBlock := generateBlock(&blockchain[len(blockchain)-1], &[]transaction{newTransaction})

	if isBlockValid(&newBlock, &blockchain[len(blockchain)-1]) {
		blockchain = append(blockchain, newBlock)
	}

	printBlockchain()
}
