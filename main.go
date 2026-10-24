package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type codeRequest struct{ Phone, Code string }

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	service := LoginService{SMS: client}
	http.HandleFunc("/login/code", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var in codeRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := service.RequestCode(in.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "code_sent"})
	})
	http.HandleFunc("/login/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var in codeRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		session, err := service.VerifyCode(in.Phone, in.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"session": session})
	})
	log.Println("property login listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
