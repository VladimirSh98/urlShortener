package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

var globalUrlStorage = map[string]string{}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func run() error {
	router := http.NewServeMux()
	router.HandleFunc("/", createShortUrl)
	router.HandleFunc("GET /{id}", returnFullUrl)
	return http.ListenAndServe(":8080", router)
}

func createShortUrl(res http.ResponseWriter, req *http.Request) {
	if !(req.Method == http.MethodPost && req.URL.Path == "/") {
		res.WriteHeader(http.StatusBadRequest)
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	// check url in body
	urlMask := createRandomMask()
	globalUrlStorage[urlMask] = string(body)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	responseUrl := fmt.Sprintf("localhost:8080/%s", urlMask)
	res.Write([]byte(responseUrl))
}

func createRandomMask() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func returnFullUrl(res http.ResponseWriter, req *http.Request) {
	urlId := req.PathValue("id")
	resultUrl, ok := globalUrlStorage[urlId]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", resultUrl)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
