package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	svr := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Serving on port %v", port)
	log.Fatal(svr.ListenAndServe())
}
