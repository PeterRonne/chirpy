package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "Welcome til min server")
	// })

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
