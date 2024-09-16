package main

import (
	"log/slog"
	"os"

	"github.com/junaidk/eth-parser/http"
	"github.com/junaidk/eth-parser/inmem"
	"github.com/junaidk/eth-parser/parser"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	repo := inmem.NewInMemEthRepository()
	parser := parser.New("https://ethereum-rpc.publicnode.com", repo)
	server := http.NewServer(parser, ":8081")
	server.Start()
}
