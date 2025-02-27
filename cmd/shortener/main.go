package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"math/rand"
	"net/http"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseUrl       string `env:"BASE_URL"`
}

func main() {
	parseFlags()
	err := run()
	if err != nil {
		panic(err)
	}
}

var globalURLStorage = map[string]string{}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func run() error {
	router := chi.NewMux()
	router.Post("/", createShortURL)
	router.Get("/{id}", returnFullURL)

	return http.ListenAndServe(flagRunAddr, router)
}

func createShortURL(res http.ResponseWriter, req *http.Request) {
	if !(req.Method == http.MethodPost && req.URL.Path == "/") {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	urlMask := createRandomMask()
	globalURLStorage[urlMask] = string(body)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	responseURL := fmt.Sprintf("%s/%s", flagResultAddr, urlMask)
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
