package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
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
	genesisTransaction := transaction.Transaction{
		Inputs: []transaction.TransactionInput{},
		Outputs: []transaction.TransactionOutput{
			{
				Value:     1000, // Initial supply
				PublicKey: nil,  // Placeholder for recipient's public key
			},
		},
	}

	privateKey := generatePrivateKey()

	signTransaction(&genesisTransaction, privateKey)

	return &genesisTransaction
}

// Sign a transaction input using the provided private key
func signTransaction(tx *transaction.Transaction, privateKey *ecdsa.PrivateKey) {
	for _, input := range tx.Inputs {
		data := sha256.Sum256(*serializeTransactionData(tx))

		r, s, err := ecdsa.Sign(rand.Reader, privateKey, data[:])
		if err != nil {
			slog.With(slog.String("error", err.Error())).Error("Failed to sign transaction: %v", err)
			os.Exit(1)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		input.Signature = append(signature, byte(0)) // Append recovery ID (0 or 1)
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

func (bc *blockchain) createNewTransaction(previousBlock *block) *transaction.Transaction {
	transactionValue := 123456 // Sending 1.23456 BTC
	transactionFee := 10000    // Transaction fee (e.g., 0.0001 BTC)

	newTx := transaction.Transaction{
		Inputs: []transaction.TransactionInput{
			{
				TransactionID: calculateHash(previousBlock),
				OutputIndex:   0,
				Signature:     nil,
			},
		},
		Outputs: []transaction.TransactionOutput{
			{
				Value:     transactionValue,
				PublicKey: nil,
			},
			{
				Value:     previousBlock.transactions[0].Outputs[0].Value - transactionFee, // Return change to the sender
				PublicKey: nil,
			},
		},
	}

	privateKey := generatePrivateKey()

	signTransaction(&newTx, privateKey)

	return &newTx
}

func generatePrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		slog.With(slog.String("error", err.Error())).Error("Failed to generate new ECDSA key")
		os.Exit(1)
	}

	return privateKey
}

// Create a new block with the given data
func (bc *blockchain) generateBlock() {
	previousBlock := bc.chain[len(bc.chain)-1]
	tx := bc.createNewTransaction(previousBlock)

	newBlock := block{
		index:        previousBlock.index + 1,
		timestamp:    time.Now().Unix(),
		transactions: []*transaction.Transaction{tx},
		previousHash: previousBlock.hash,
	}
	newBlock.hash = calculateHash(&newBlock)

	bc.chain = append(bc.chain, &newBlock)
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

	bc.generateBlock()

	bc.printBlockchain()
}
