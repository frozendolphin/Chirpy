package main

import (
	"log"
	"net/http"
)

func main() {
	server := http.NewServeMux()

	server_struct := http.Server {
		Handler: server,
		Addr: ":8080",
	}

	err := server_struct.ListenAndServe()
	if err != nil {
		log.Fatalf("err occured: %v", err)
	}

	
} 