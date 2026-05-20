package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HtetAungKhant23/shortix/api"
	"github.com/HtetAungKhant23/shortix/shortener"
	"github.com/HtetAungKhant23/shortix/store"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)

	memStore := store.NewInMemoryStore()
	svc := shortener.New(memStore)
	router := api.GetRootRouter(svc)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("shortix server listening on %s", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
