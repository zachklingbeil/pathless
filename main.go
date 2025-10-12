package main

import (
	_ "embed"
	"html/template"
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

	tmpl, err := template.New("page").Parse(zero)
	if err != nil {
		panic("template parse error: " + err.Error())
	}

	var b strings.Builder
	if err := tmpl.Execute(&b, map[string]string{
		"Title":  title,
		"APIURL": apiURL,
	}); err != nil {
		panic("template execute error: " + err.Error())
	}
	one = []byte(b.String())
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
