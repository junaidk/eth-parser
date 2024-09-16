package parser

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/junaidk/eth-parser/inmem"
)

type mockRoundTripper struct {
	response []*http.Response
}

func (rt *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := rt.response[0]
	rt.response = rt.response[1:]
	return resp, nil
}

func createMockResponse(t *testing.T, json string) *http.Response {
	recorder := httptest.NewRecorder()
	recorder.Header().Add("Content-Type", "application/json")
	recorder.WriteString(json)
	return recorder.Result()
}

func TestGetCurrentBlock(t *testing.T) {
	repo := inmem.NewInMemEthRepository()
	json := `{"jsonrpc":"2.0","result":"0x13cd6df","id":1}`
	expectedResponse := createMockResponse(t, json)

	parser := New("https://test.com", repo)

	parser.client = &client{
		url:        "https://test.com",
		httpClient: &http.Client{Transport: &mockRoundTripper{[]*http.Response{expectedResponse}}},
	}

	currentBlock := parser.GetCurrentBlock()
	if currentBlock == 0 {
		t.Errorf("GetCurrentBlock returned 0, expected a non-zero value")
	}

	if currentBlock != 20764383 {
		t.Errorf("GetCurrentBlock returned %d, expected 13123135", currentBlock)
	}
}

func TestSubscribe(t *testing.T) {
	repo := inmem.NewInMemEthRepository()
	parser := New("https://test.com", repo)

	parser.blockLimit = 0

	json := `{"jsonrpc":"2.0","result":"0x13cd6df","id":1}`
	expectedResponse1 := createMockResponse(t, json)

	parser.client = &client{
		url:        "https://test.com",
		httpClient: &http.Client{Transport: &mockRoundTripper{[]*http.Response{expectedResponse1}}},
	}

	res := parser.Subscribe("0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	if !res {
		t.Errorf("Subscribe returned false, expected true")
	}

	sub, err := parser.repo.GetSubscriptionByAddress("0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	if err != nil {
		t.Errorf("GetSubscriptionByAddress returned an error: %v", err)
	}
	if sub == nil {
		t.Fatalf("GetSubscriptionByAddress returned nil, expected a subscription")
	}
	if sub.Address != "0x69c93309b1c9a9452e9fe445468461efb2a72dfe" {
		t.Errorf("GetSubscriptionByAddress returned wrong address, expected 0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	}

}

func TestGetTransactions(t *testing.T) {
	repo := inmem.NewInMemEthRepository()
	parser := New("https://test.com", repo)

	parser.blockLimit = 1

	json := `{"jsonrpc":"2.0","result":"0x13cd6df","id":1}`
	expectedResponse1 := createMockResponse(t, json)
	expectedResponse2 := createMockResponse(t, SubscribeBlockdata)

	parser.client = &client{
		url:        "https://test.com",
		httpClient: &http.Client{Transport: &mockRoundTripper{[]*http.Response{expectedResponse1, expectedResponse2}}},
	}

	res := parser.Subscribe("0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	if !res {
		t.Errorf("Subscribe returned false, expected true")
	}

	sub, err := parser.repo.GetSubscriptionByAddress("0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	if err != nil {
		t.Errorf("GetSubscriptionByAddress returned an error: %v", err)
	}
	if sub == nil {
		t.Fatalf("GetSubscriptionByAddress returned nil, expected a subscription")
	}
	if sub.Address != "0x69c93309b1c9a9452e9fe445468461efb2a72dfe" {
		t.Errorf("GetSubscriptionByAddress returned wrong address, expected 0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	}

	tx, err := parser.repo.GetTransactionsByAddress("0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	if err != nil {
		t.Errorf("GetTransactionsByAddress returned an error: %v", err)
	}
	if len(tx) == 0 {
		t.Fatalf("GetTransactionsByAddress returned no transactions, expected at least one")
	}
	if tx[0].From != "0x69c93309b1c9a9452e9fe445468461efb2a72dfe" {
		t.Errorf("GetTransactionsByAddress returned wrong from address, expected 0x69c93309b1c9a9452e9fe445468461efb2a72dfe")
	}
}

var SubscribeBlockdata = `{
    "jsonrpc": "2.0",
    "result": {
        "baseFeePerGas": "0x26aa029e0",
        "blobGasUsed": "0x20000",
        "difficulty": "0x0",
        "excessBlobGas": "0xc0000",
        "extraData": "0x546974616e2028746974616e6275696c6465722e78797a29",
        "gasLimit": "0x1c9c380",
        "gasUsed": "0xd0b6c8",
        "hash": "0xb4ad8fb8b0bd3880bec7981358febc0da3a483fbe0cce6aa50864172f8f1d3e6",
        "logsBloom": "0x15ade3122401011ce368ea60e170a6815b822193e41b010c4bc5018ee1bafda6445f040b010d802822d1130702ca9300d2d3e5f88f47313837a92d00a12c4d144488e0489a43816d7c046db9a14193204c1bc31619f658a891a082899ca41c090649881e06e0c1646609baa2722c2eda256093e519620c03cae24314231fc350f28ab80166898b00d46ae53442a4000e100eb821018191782c2e284250ba303d0b2a45a08111b483e6a750e1ec16814b0787519f2884082465053d0a31002b6a0840102322622f342d8c802540cd491dc230311e0e0470900201364268626a1360502449860682469e0415f100ad02bfa1ee931046a424c05224f833c1cfa2df",
        "miner": "0x4838b106fce9647bdf1e7877bf73ce8b0bad5f97",
        "mixHash": "0x22ba162cd561ff21da85dde568c32bd92210320b383882ed8d951d29ce0b78ac",
        "nonce": "0x0000000000000000",
        "number": "0x13cd322",
        "parentBeaconBlockRoot": "0x326f51ffca829403766277de12c1695a142e6458d15537c1c81aab9f36c25beb",
        "parentHash": "0x272b6b626101adf46ab7a3d7b13ab073d18f2713182f20b38aaa9fde51bd4079",
        "receiptsRoot": "0xd66ad04ed29f6b61c19b5351986aa9798d7c9402139ba90691b00db385808f72",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0x119cc",
        "stateRoot": "0x805cf154f4cb14240d50d2435017e277c9ac4b406c3f088b4bdac18a3ab7edf2",
        "timestamp": "0x66e83013",
        "totalDifficulty": "0xc70d815d562d3cfa955",
        "transactions": [
            {
                "blockHash": "0xb4ad8fb8b0bd3880bec7981358febc0da3a483fbe0cce6aa50864172f8f1d3e6",
                "blockNumber": "0x13cd322",
                "chainId": "0x1",
                "from": "0x69c93309b1c9a9452e9fe445468461efb2a72dfe",
                "gas": "0xf4240",
                "gasPrice": "0x2a1090936",
                "hash": "0x8c500d7e6c2412bb6c519e668b6e56162395c18def97d9888282a9929a1fed52",
                "input": "0x8803dbee0000000000000000000000000000000000000000000a1c13d6e81e68bde0000000000000000000000000000000000000000000000000000044e1a04f55ccb9ac00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000069c93309b1c9a9452e9fe445468461efb2a72dfe0000000000000000000000000000000000000000000000000000000066e837140000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000088ce174c655b6d11210a069b2c106632dabdb068",
                "nonce": "0x9",
                "r": "0x9ef032b25200f1f19e001272e0007f90dcf176f2161536b233726185817ec778",
                "s": "0x63e56d931c7d460e37d7037110406df6bf7de01f4911652520768fce147c21b6",
                "to": "0x7a250d5630b4cf539739df2c5dacb4c659f2488d",
                "transactionIndex": "0x0",
                "type": "0x0",
                "v": "0x25",
                "value": "0x0"
            }
        ]
    },
    "id": 1
}`
