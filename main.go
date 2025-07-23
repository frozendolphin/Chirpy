package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	server_struct := http.Server {
		Handler: mux,
		Addr: ":8080",
	}

	err := server_struct.ListenAndServe()
	if err != nil {
		log.Fatalf("err occured: %v", err)
	}

	
} 