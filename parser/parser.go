package parser

import (
	"errors"
	"strconv"

	parser "github.com/junaidk/eth-parser"
)

type BlockchainParser struct {
	url        string
	client     *client
	repo       parser.TransactionRepository
	blockLimit int
}

func New(url string, repo parser.TransactionRepository) *BlockchainParser {
	return &BlockchainParser{
		url:        url,
		client:     newClient(url),
		repo:       repo,
		blockLimit: 100,
	}
}

func (p *BlockchainParser) GetCurrentBlock() int {
	result, err := p.client.call("eth_blockNumber", nil)
	if err != nil {
		return 0
	}

	blockNumber, ok := result.(string)
	if !ok {
		return 0
	}

	blockInt, err := strconv.ParseInt(blockNumber[2:], 16, 64)
	if err != nil {
		return 0
	}

	return int(blockInt)
}

func (p *BlockchainParser) Subscribe(address string) bool {
	sub, err := p.repo.GetSubscriptionByAddress(address)
	if !errors.Is(err, parser.ErrSubscriptionNotFound) {
		return false
	}

	if sub != nil && sub.Address == address {
		return true
	}

	_, err = p.repo.AddSubscription(address)
	if err != nil {
		return false
	}

	tx := p.getTransactions(address)

	err = p.repo.SaveTransactions(tx)
	return err == nil
}

func (p *BlockchainParser) getTransactions(address string) []parser.Transaction {
	transactions := []parser.Transaction{}

	latestBlock := p.GetCurrentBlock()

	// Iterate through the last blockLimit blocks or up to the latest block, whichever is smaller
	for i := 0; i < p.blockLimit && latestBlock-i >= 0; i++ {
		blockNumber := latestBlock - i

		blockNumberHex := "0x" + strconv.FormatInt(int64(blockNumber), 16)

		result, err := p.client.call("eth_getBlockByNumber", []interface{}{blockNumberHex, true})
		if err != nil {
			continue
		}

		block, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		blockTransactions, ok := block["transactions"].([]interface{})
		if !ok {
			continue
		}

		// Iterate through transactions in the block
		for _, tx := range blockTransactions {
			transaction, ok := tx.(map[string]interface{})
			if !ok {
				continue
			}

			// Check if the transaction involves the given address
			if transaction["from"] == address || transaction["to"] == address {
				transactions = append(transactions, parser.Transaction{
					From:  transaction["from"].(string),
					To:    transaction["to"].(string),
					Value: transaction["value"].(string),
					Hash:  transaction["hash"].(string),
					Block: blockNumber,
				})
			}
		}
	}

	return transactions
}

func (p *BlockchainParser) GetTransactions(address string) []parser.Transaction {
	txs, err := p.repo.GetTransactionsByAddress(address)
	if err != nil {
		return []parser.Transaction{}
	}

	return txs
}
