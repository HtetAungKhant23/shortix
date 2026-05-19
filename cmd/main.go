package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HtetAungKhant23/shortix/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	router := api.GetRootRouter()

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("shortix server listening on %s", port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
