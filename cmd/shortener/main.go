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

var globalURLStorage = map[string]string{}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func run() error {
	router := http.NewServeMux()
	router.HandleFunc("/", createShortURL)
	router.HandleFunc("GET /{id}", returnFullURL)
	return http.ListenAndServe(":8080", router)
}

func createShortURL(res http.ResponseWriter, req *http.Request) {
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
	globalURLStorage[urlMask] = string(body)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	responseURL := fmt.Sprintf("localhost:8080/%s", urlMask)
	res.Write([]byte(responseURL))
}

func createRandomMask() string {
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func returnFullURL(res http.ResponseWriter, req *http.Request) {
	urlID := req.PathValue("id")
	resultURL, ok := globalURLStorage[urlID]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", resultURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
