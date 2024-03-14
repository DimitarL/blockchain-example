package transaction

type Transaction struct {
	// ID      []byte
	Inputs  []TransactionInput
	Outputs []TransactionOutput
}

type TransactionInput struct {
	TransactionID []byte // stores id of previous output
	OutputIndex   int
	// Signature is a script which provides data to be used in an output’s script PublicKey. If the data is correct, the output can be unlocked, and its value can be used to generate new outputs; if it’s not correct, the output cannot be referenced in the input.
	Signature []byte
}

type TransactionOutput struct {
	Value     int    // amount of coins
	PublicKey []byte // script	// wallet address for now
}
