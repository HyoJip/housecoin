package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hyojip/housecoin/blockchain"
	"github.com/Hyojip/housecoin/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var port string

type url string

type urlDescription struct {
	Url         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type balanceResponse struct {
	Address string `json:"address,omitempty"`
	Amount  int    `json:"amount"`
}

type addTransactionPayload struct {
	To     string `json:"to,omitempty"`
	Amount int    `json:"amount"`
}

func (u url) marshalText() (text []byte, err error) {
	link := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(link), nil
}

func Start(aPort string) {
	port = aPort

	handler := mux.NewRouter()
	handler.Use(jsonContentTypeMiddleware) // go에서 만든 http.Handler의 SPI를 만족하기 위해서
	// HandlerFunc라는 함수 프로토타입을 지정 후 덕타이핑에 필요한 함수(ServeHTTP) 구현
	// HandlerFunc(func) 타입으로 생성할 경우, SPI를 만족하는 덕타이핑 함수가 상속(prototype)
	// 사용자는 인자로 넘기는 func만 구현하면 알아서 http.Handler를 만족하는 어댑터가 만들어짐
	// 약간 FunctionalInterface? 혹은 일급컬렉션과 비슷한 느낌? 하나의 합성필드만 가지는 wrap class
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	handler.HandleFunc("/status", status).Methods("GET")
	handler.HandleFunc("/balance/{address}", balanceAddress).Methods("GET")
	handler.HandleFunc("/mempool", mempool).Methods("GET")
	handler.HandleFunc("/transactions", transactions).Methods("POST")

	fmt.Printf("Start REST server http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})

}

func documentation(writer http.ResponseWriter, _ *http.Request) {
	descriptions := []urlDescription{
		{
			Url:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			Url:         url("/blocks"),
			Method:      "GET",
			Description: "Show All blocks",
		},
		{
			Url:         url("/blocks"),
			Method:      "POST",
			Description: "Add New Block",
			Payload:     "data:string",
		},
		{
			Url:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See The Block",
		},
		{
			Url:         url("status"),
			Method:      "GET",
			Description: "Show Blockchain Status",
		},
		{
			Url:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Show Address's Balance",
		},
		{
			Url:         url("/mempool"),
			Method:      "GET",
			Description: "Show Mempool Transactions",
		},
	}
	utils.HandleError(json.NewEncoder(writer).Encode(descriptions))
}

func blocks(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		utils.HandleError(json.NewEncoder(writer).Encode(blockchain.FindBlocks()))
	case "POST":
		blockchain.GetBlockchain().AddBlock()
		writer.WriteHeader(http.StatusCreated)
	}
}

func block(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	hash := vars["hash"]

	theBlock, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(writer)
	if errors.Is(err, blockchain.ErrNotFound) {
		writer.WriteHeader(http.StatusNotFound)
		utils.HandleError(encoder.Encode(errorResponse{fmt.Sprint(err)}))
		return
	}
	utils.HandleError(encoder.Encode(theBlock))
}
func status(writer http.ResponseWriter, _ *http.Request) {
	utils.HandleError(json.NewEncoder(writer).Encode(blockchain.GetBlockchain()))
}

func balanceAddress(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	address := vars["address"]
	isTotal := request.URL.Query().Get("total")
	var body interface{}
	switch isTotal {
	case "true":
		totalAmount := blockchain.FindBalanceByAddress(address)
		body = balanceResponse{address, totalAmount}
	default:
		body = blockchain.FindUTxOutsByAddress(address)
	}
	utils.HandleError(json.NewEncoder(writer).Encode(body))
}

func mempool(writer http.ResponseWriter, _ *http.Request) {
	utils.HandleError(json.NewEncoder(writer).Encode(blockchain.Mempool.Txs))
}

func transactions(writer http.ResponseWriter, request *http.Request) {
	var payload addTransactionPayload
	utils.HandleError(json.NewDecoder(request.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		utils.HandleError(json.NewEncoder(writer).Encode(errorResponse{"Not enough funds"}))
	}
	writer.WriteHeader(http.StatusCreated)

}
