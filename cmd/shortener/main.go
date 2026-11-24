package main

import (
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"sync"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength = 8
	baseURL  = "http://localhost:8080"
)

var (
	urls = make(map[string]string)
	mu   sync.RWMutex
)

func generateShortID() string {
	buf := make([]byte, idLength)
	for i := 0; i < idLength; i++ {
		ind := rand.IntN(len(alphabet))
		buf[i] = alphabet[ind]
	}
	return string(buf)

}

func handlePost(res http.ResponseWriter, req *http.Request) {
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	body := string(rawBody)
	if _, err := url.ParseRequestURI(body); err != nil {
		http.Error(res, "invalid url", http.StatusBadRequest)
		return
	}

	id := generateShortID()

	mu.Lock()
	urls[id] = body
	mu.Unlock()

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(baseURL + "/" + id))

}

func handleGet(res http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		http.Error(res, "id is empty", http.StatusBadRequest)
		return
	}

	mu.RLock()
	url, ok := urls[id]
	mu.RUnlock()
	if !ok {
		http.Error(res, "id is not found", http.StatusBadRequest)
		return
	}

	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{$}", handlePost)
	mux.HandleFunc("GET /{id}", handleGet)
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.Error(res, "bad request", http.StatusBadRequest)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
