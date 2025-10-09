package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

var pathless []byte

func main() {
	socketPath := "/tmp/pathless.sock"
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	os.Chmod(socketPath, 0666)
	http.HandleFunc("/", one)
	http.Serve(listener, nil)
}

// Redirect anything that's not exactly "/" to "/"
func one(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.URL.RawQuery != "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(pathless)
}

func init() {
	apiURL := os.Getenv("API_URL")
	apiURL = "https://" + apiURL

	pathless = []byte(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>hello_universe</title>
		<style>
			*,
			*::before,
			*::after {
				box-sizing: border-box;
				margin: 0;
				scrollbar-width: none;
				-ms-overflow-style: none;
				user-select: none;
			}
			*::-webkit-scrollbar {
				display: none;
			}
			html {
				scroll-behavior: smooth;
			}
			body {
				color: #f1f1f1;
				background-color: #000000;
				height: 100vh;
				font-family: 'Roboto', sans-serif;
				overflow: hidden;
				display: flex;
				border: medium solid blue;
				border-radius: 0.3125em;
			}

			a,
			a:hover,
			a:visited,
			a:active,
			a:focus {
				text-align: center;
				color: inherit;
				text-decoration: underline;
			}
		</style>
		<script>
			const API_URL = '%s';
			const nav = { prev: 0, next: 0 };

			function loadFrame(frameIndex) {
				fetch(API_URL, { headers: { Y: frameIndex } })
					.then((response) => {
						nav.prev = parseInt(response.headers.get('X'));
						nav.next = parseInt(response.headers.get('Z'));
						return response.text();
					})
					.then((html) => {
						document.getElementById('frame').innerHTML = html;
					})
					.catch(console.error);
			}

			document.addEventListener('keydown', (event) => {
				switch (event.key.toLowerCase()) {
					case 'e':
						loadFrame(nav.next);
						break;
					case 'q':
						loadFrame(nav.prev);
						break;
				}
			});

			document.addEventListener('DOMContentLoaded', () => {
				loadFrame(0);
			});
		</script>
	</head>
	<body>
		<div id="frame"></div>
	</body>
</html>`, apiURL))
}
