package inmem

import (
	"sync"

	parser "github.com/junaidk/eth-parser"
)

// InMemEthRepository implements EthRepository using an in-memory store
type InMemEthRepository struct {
	transactions  map[string][]parser.Transaction
	subscriptions map[int]parser.Subscription
	mutex         sync.RWMutex
	nextSubID     int
}

// NewInMemEthRepository creates a new instance of InMemEthRepository
func NewInMemEthRepository() *InMemEthRepository {
	return &InMemEthRepository{
		transactions:  make(map[string][]parser.Transaction),
		subscriptions: make(map[int]parser.Subscription),
		nextSubID:     1,
	}
}

// SaveTransaction adds a new transaction to the in-memory store
func (r *InMemEthRepository) SaveTransactions(txs []parser.Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, tx := range txs {
		r.transactions[tx.From] = append(r.transactions[tx.From], tx)
		r.transactions[tx.To] = append(r.transactions[tx.To], tx)
	}
	return nil
}

func (r *InMemEthRepository) GetTransactionsByAddress(address string) ([]parser.Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	txs := make([]parser.Transaction, 0, len(r.transactions))
	for _, tx := range r.transactions {
		for _, t := range tx {
			if t.To == address || t.From == address {
				txs = append(txs, t)
			}
		}
	}
	return txs, nil
}

// AddSubscription adds a new address subscription
func (r *InMemEthRepository) AddSubscription(address string) (*parser.Subscription, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	sub := &parser.Subscription{
		ID:      r.nextSubID,
		Address: address,
	}
	r.subscriptions[sub.ID] = *sub
	r.nextSubID++
	return sub, nil
}

// GetSubscriptions retrieves all subscriptions
// func (r *InMemEthRepository) GetSubscriptions() ([]parser.Subscription, error) {
// 	r.mutex.RLock()
// 	defer r.mutex.RUnlock()

// 	subs := make([]parser.Subscription, 0, len(r.subscriptions))
// 	for _, sub := range r.subscriptions {
// 		subs = append(subs, sub)
// 	}
// 	return subs, nil
// }

func (r *InMemEthRepository) GetSubscriptionByAddress(address string) (*parser.Subscription, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, sub := range r.subscriptions {
		if sub.Address == address {
			return &sub, nil
		}
	}
	return nil, parser.ErrSubscriptionNotFound
}
