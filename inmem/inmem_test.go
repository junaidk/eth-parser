package inmem

import (
	"testing"

	parser "github.com/junaidk/eth-parser"
)

func TestAddSubscription(t *testing.T) {
	repo := NewInMemEthRepository()
	sub, err := repo.AddSubscription("0x123")
	if err != nil {
		t.Errorf("AddSubscription returned an error: %v", err)
	}
	if sub.ID != 1 {
		t.Errorf("AddSubscription returned an incorrect ID: %v", sub.ID)
	}
}

func TestGetSubscriptionByAddress(t *testing.T) {
	repo := NewInMemEthRepository()
	repo.AddSubscription("0x123")
	sub, err := repo.GetSubscriptionByAddress("0x123")
	if err != nil {
		t.Errorf("GetSubscriptionByAddress returned an error: %v", err)
	}
	if sub.ID != 1 {
		t.Errorf("GetSubscriptionByAddress returned an incorrect ID: %v", sub.ID)
	}
	if sub.Address != "0x123" {
		t.Errorf("GetSubscriptionByAddress returned an incorrect address: %v", sub.Address)
	}
}

func TestSaveTransactions(t *testing.T) {
	repo := NewInMemEthRepository()
	tx := []parser.Transaction{
		{
			From:  "0x123",
			To:    "0x456",
			Value: "100",
			Hash:  "0x789",
			Block: 1,
		},
	}
	err := repo.SaveTransactions(tx)
	if err != nil {
		t.Errorf("SaveTransactions returned an error: %v", err)
	}
	txs, err := repo.GetTransactionsByAddress("0x123")
	if err != nil {
		t.Errorf("GetTransactionsByAddress returned an error: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("GetTransactionsByAddress returned an incorrect number of transactions: %v", len(txs))
	}
}

func TestGetTransactionsByAddress(t *testing.T) {
	repo := NewInMemEthRepository()
	tx := []parser.Transaction{
		{
			From:  "0x123",
			To:    "0x456",
			Value: "100",
			Hash:  "0x789",
			Block: 1,
		},
	}
	err := repo.SaveTransactions(tx)
	if err != nil {
		t.Errorf("SaveTransactions returned an error: %v", err)
	}
	txs, err := repo.GetTransactionsByAddress("0x123")
	if err != nil {
		t.Errorf("GetTransactionsByAddress returned an error: %v", err)
	}
	if len(txs) != 2 {
		t.Errorf("GetTransactionsByAddress returned an incorrect number of transactions: %v", len(txs))
	}
	if txs[0].From != "0x123" {
		t.Errorf("GetTransactionsByAddress returned an incorrect address: %v", txs[0].From)
	}
	if txs[0].To != "0x456" {
		t.Errorf("GetTransactionsByAddress returned an incorrect address: %v", txs[0].To)
	}
	if txs[0].Value != "100" {
		t.Errorf("GetTransactionsByAddress returned an incorrect value: %v", txs[0].Value)
	}
	if txs[0].Hash != "0x789" {
		t.Errorf("GetTransactionsByAddress returned an incorrect hash: %v", txs[0].Hash)
	}
}
