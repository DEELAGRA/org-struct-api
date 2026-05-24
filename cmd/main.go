package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DEELAGRA/org-struct-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration loading error: %v\n", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	addr := ":" + strconv.Itoa(cfg.ServerPort)
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed started server: %v", err)
	}
}
