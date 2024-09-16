# Blockchain Parser

This is a simple blockchain parser that can be used to parse transactions from a eth blockchain and subscribe to new transactions for a given address.


## Usage

```
go run cmd/server/main.go
```

Interface is exposed at http://localhost:8081/ and has the following endpoints:

http://localhost:8081/getcurrentblock

http://localhost:8081/subscribe?address=0xb1ce613ca397b76a74ad86d5a1e064d1d59abd0d

http://localhost:8081/gettransactions?address=0xb1ce613ca397b76a74ad86d5a1e064d1d59abd0d

### Running tests

```
go test ./...
```

## Architecture

Subscribe enpoint creates a new subscription and store transaction for last 100 blocks in the database for a given address.
In production, we would use polling to check for new transactions and store them in the database.

GetTransactions enpoint fetches all transactions for a given address from the database.

GetCurrentBlock enpoint fetches the current block number from the blockchain.

