package parser

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
	Hash  string `json:"hash"`
	Block int    `json:"blockNumber"`
}

type Subscription struct {
	ID      int
	Address string
}

type TransactionRepository interface {
	SaveTransactions(tx []Transaction) error
	GetTransactionsByAddress(address string) ([]Transaction, error)
	AddSubscription(address string) (*Subscription, error)
	GetSubscriptionByAddress(address string) (*Subscription, error)
}
