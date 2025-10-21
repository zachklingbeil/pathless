package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"html/template"
	"net/http"
	"os"
	"regexp"
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
	if apiURL == "" {
		apiURL = "http://localhost:1002"
	} else if !strings.HasPrefix(apiURL, "http://") && !strings.HasPrefix(apiURL, "https://") {
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

	minified := minify(b.String())
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(minified)); err != nil {
		panic("gzip write error: " + err.Error())
	}
	if err := gz.Close(); err != nil {
		panic("gzip close error: " + err.Error())
	}
	one = buf.Bytes()
}

func minify(html string) string {
	// Remove HTML comments
	html = regexp.MustCompile(`<!--[\s\S]*?-->`).ReplaceAllString(html, "")
	// Remove whitespace between tags
	html = regexp.MustCompile(`>\s+<`).ReplaceAllString(html, "><")
	// Collapse multiple spaces/tabs/newlines to single space
	html = regexp.MustCompile(`\s+`).ReplaceAllString(html, " ")
	// Remove spaces around specific characters
	html = strings.ReplaceAll(html, " >", ">")
	html = strings.ReplaceAll(html, "< ", "<")
	// Trim leading/trailing whitespace
	return strings.TrimSpace(html)
}

func Pathless(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.URL.RawQuery != "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write(one)
}

func main() {
	http.HandleFunc("/", Pathless)
	http.ListenAndServe(":1001", nil)
}
