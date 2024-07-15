package server

import (
	"log"
	"net/http"
	"video-streaming/internal/handlers"

	"github.com/gorilla/mux"
)

func StartServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/upload", handlers.UploadPageHandler).Methods("GET")
	r.HandleFunc("/upload", handlers.UploadHandler).Methods("POST")
	r.HandleFunc("/stream", handlers.StreamHandler).Methods("GET")
	r.HandleFunc("/ws", handlers.WebSocketHandler).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
