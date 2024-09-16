package http

import (
	"log/slog"
	"net/http"

	parser "github.com/junaidk/eth-parser"
)

type Server struct {
	port   string
	parser parser.Parser
}

func NewServer(parser parser.Parser, port string) *Server {
	return &Server{parser: parser, port: port}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/subscribe", s.subscribe)
	mux.HandleFunc("/gettransactions", s.getTransactions)
	mux.HandleFunc("/getcurrentblock", s.getCurrentBlock)

	slog.Info("Server is running", "port", s.port)
	err := http.ListenAndServe(s.port, mux)
	if err != nil {
		slog.Error("Error starting server", "error", err)
		panic(err)
	}
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	success := s.parser.Subscribe(address)
	if success {
		s.writeJSON(w, http.StatusOK, envelope{"message": "Subscription successful"}, nil)
	} else {
		s.writeJSON(w, http.StatusOK, envelope{"message": "Subscription failed"}, nil)
	}
}

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	transactions := s.parser.GetTransactions(address)
	s.writeJSON(w, http.StatusOK, envelope{"transactions": transactions}, nil)
}

func (s *Server) getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := s.parser.GetCurrentBlock()
	s.writeJSON(w, http.StatusOK, envelope{"block": block}, nil)
}
