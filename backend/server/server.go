package main

import (
	"fmt"
	"net/http"
	datahandler "social-network/backend/internal/data"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	datahandler.OpenDb()
	defer datahandler.Database.Close()

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request recieved")
		handleWebSocketConnection(w, r)
	})
	http.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("image request recieved")
		handleImageRequest(w, r)
	})
	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", addCorsHeaders(http.DefaultServeMux))
}

func addCorsHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests only from specific origin
		allowedOrigin := "http://127.0.0.1:5173" // Replace with your specific origin
		// Check if the request origin is allowed
		if origin := r.Header.Get("Origin"); origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		// Handle preflight options requests
		if r.Method == http.MethodOptions {
			return
		}
		// Continue processing request
		handler.ServeHTTP(w, r)
	})
}
