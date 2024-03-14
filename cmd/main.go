package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/DimitarL/blockchain-example/transaction"
)

// type transaction struct {
// 	sender         string
// 	receiver       string
// 	amount         float64
// 	transactionFee float64
// 	signature      string
// 	timestamp      string
// }

type block struct {
	index        int
	timestamp    int64
	transactions []*transaction.Transaction
	previousHash []byte
	hash         []byte
}

type blockchain struct {
	chain []*block
}

// Create the genesis block (first block in the blockchain)
func (bc *blockchain) createGenesisBlock() {
	tx := createGenesisTransaction()
	genesisBlock := block{
		index:        0,
		timestamp:    time.Now().Unix(),
		transactions: []*transaction.Transaction{tx},
		previousHash: nil,
	}
	genesisBlock.hash = calculateHash(&genesisBlock)

	bc.chain = append(bc.chain, &genesisBlock)
}

// Create a new transaction
func createGenesisTransaction() *transaction.Transaction {
	return &transaction.Transaction{
		Inputs: []transaction.TransactionInput{},
		Outputs: []transaction.TransactionOutput{
			{
				Value:     1000, // Initial supply
				PublicKey: nil,  // Placeholder for recipient's public key
			},
		},
	}
}

// Calculate the SHA-256 hash of a block
func calculateHash(block *block) []byte {
	var data []byte
	data = append(data, *int64ToBytes(int64(block.index))...)
	data = append(data, block.previousHash...)
	data = append(data, *int64ToBytes(block.timestamp)...)
	for _, tx := range block.transactions {
		data = append(data, *serializeTransactionData(tx)...)
	}
	hashInBytes := sha256.Sum256(data)

	return hashInBytes[:]
}

func int64ToBytes(number int64) *[]byte {
	numberInBytes := make([]byte, 8)

	binary.BigEndian.PutUint64(numberInBytes, uint64(number))

	return &numberInBytes
}

// Create a consistent representation of the transaction's content
// in a format that can be hashed
func serializeTransactionData(tx *transaction.Transaction) *[]byte {
	var data []byte

	for _, input := range tx.Inputs {
		data = append(data, input.TransactionID...)
		data = append(data, byte(input.OutputIndex))
	}

	for _, output := range tx.Outputs {
		data = append(data, output.PublicKey...)
		data = append(data, *int64ToBytes(int64(output.Value))...)
	}

	return &data
}

func (bc *blockchain) printBlockchain() {
	for _, block := range bc.chain {
		fmt.Printf("Index: %d\n", block.index)
		fmt.Printf("Timestamp: %d\n", block.timestamp)
		fmt.Printf("Transactions: %v\n", block.transactions)
		fmt.Printf("Previous Hash: %s\n", block.previousHash)
		fmt.Printf("Hash: %s\n", block.hash)
		fmt.Println()
	}
}

func main() {
	bc := blockchain{chain: []*block{}}

	bc.createGenesisBlock()

	bc.printBlockchain()
}
