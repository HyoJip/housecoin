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
	"strconv"
)

var port string

type url string

func (u url) marshalText() (text []byte, err error) {
	link := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(link), nil
}

type urlDescription struct {
	Url         url    `json:"Url"`
	Method      string `json:"Method"`
	Description string `json:"Description"`
	Payload     string `json:"Payload,omitempty"`
}

type newBlockRequest struct {
	Message string
}

type errorResponse struct {
	Message string `json:"message"`
}

type test http.Handler

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
	handler.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")

	fmt.Printf("Start web server http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})

}

func documentation(writer http.ResponseWriter, request *http.Request) {
	descriptions := []urlDescription{
		{
			Url:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			Url:         url("/blocks"),
			Method:      "POST",
			Description: "Add New Block",
			Payload:     "data:string",
		},
		{
			Url:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See The Block",
		},
	}
	json.NewEncoder(writer).Encode(descriptions)
}

func blocks(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		utils.HandleError(json.NewEncoder(writer).Encode(blockchain.AllBlock()))
	case "POST":
		var newBlockDTO newBlockRequest
		utils.HandleError(json.NewDecoder(request.Body).Decode(&newBlockDTO))
		blockchain.GetBlockchain().AddBlock(newBlockDTO.Message)
		writer.WriteHeader(http.StatusCreated)
	}
}
func block(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	height, err := strconv.Atoi(vars["height"])
	utils.HandleError(err)

	theBlock, err := blockchain.GetBlock(height)
	encoder := json.NewEncoder(writer)
	if errors.Is(err, blockchain.ErrNotFound) {
		writer.WriteHeader(http.StatusNotFound)
		encoder.Encode(errorResponse{fmt.Sprint(err)})
		return
	}
	utils.HandleError(encoder.Encode(theBlock))
}
