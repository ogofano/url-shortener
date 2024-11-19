package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

var data map[string]string

func shortURL() string {
	var short string
	alph := "0123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	for i := 0; i != 8; i++ {
		short += string(alph[rand.Intn(len(alph))])
	}
	
	return short
}

func handler(res http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	if url == "/" {
		createShortURL(res, req)
		return
	} else if len(url) == 9 {
		getURL(res, req, url)
		return
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func createShortURL(res http.ResponseWriter, req *http.Request) {
	pattern := regexp.MustCompile(`^https?:\/\/[\w.-]+(?:\/[\w\/_.?=%&-]+)?$`)
	if req.Method != http.MethodPost || !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
        http.Error(res, "Expected a POST request with Content-Type: text/plain", http.StatusBadRequest)
        return
    }
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	url := string(body)
	if !pattern.MatchString(url) {
		http.Error(res, "Invalid URL", http.StatusBadRequest)
		return
	}
	
	localhost := "http://localhost:8080/"
	shorten := shortURL()
	data["/" + shorten] = url
	res.WriteHeader(http.StatusCreated)
	res.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(res, localhost + shorten)
}

func getURL(res http.ResponseWriter, req *http.Request, url string) {
	if req.Method != http.MethodGet {
		http.Error(res, "Expected a GET request with Content-Type: text/plain", http.StatusBadRequest)
		return
	}

	val, ok := data[url]
	if !ok {
		http.Error(res, "URL not find", http.StatusNotFound)
		return
	}

	res.Header().Set("Location", val)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	data = make(map[string]string)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}