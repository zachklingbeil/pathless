package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"strings"
)

//go:embed pathless.html
var zero string
var one []byte

func init() {
	title := os.Getenv("TITLE")
	if title == "" {
		title = "hello_universe"
	}
	apiURL := os.Getenv("API_URL")
	if apiURL != "" && !strings.HasPrefix(apiURL, "http://") && !strings.HasPrefix(apiURL, "https://") {
		apiURL = "https://" + apiURL
	}
	one = fmt.Appendf(nil, zero, title, apiURL)
}

func Pathless(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.URL.RawQuery != "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(one)
}

func main() {
	http.HandleFunc("/", Pathless)
	http.ListenAndServe(":1001", nil)
}
