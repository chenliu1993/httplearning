package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

// ClientHeaderInfoHandler get clientID and secret.
func ClientHeaderInfoHandler(r *http.Request) (string, string, error) {
	if len(r.Header) == 0 {
		return "", "", fmt.Errorf("No specific headers are set")
	}
	clientID := r.Header.Get("Client_ID")
	clientSecret := r.Header.Get("Client_Secret")
	return clientID, clientSecret, nil
}

// ValidateToken validates the tooken fields.
func (server *WebServer) ValidateToken(r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return fmt.Errorf("empty field of token")
	}
	_, ok := server.CToken[token]
	if !ok {
		return fmt.Errorf("no such token")
	}
	return nil
}

// Token returns a token for the client.
func (server *WebServer) Token(w http.ResponseWriter, r *http.Request) {
	// gen a token.
	token := uuid.New().String()[:8]
	clientID := r.Header.Get("client_ID")
	clientSecret := r.Header.Get("client_Secret")
	val, ok := server.CStore[clientID]
	if !ok {
		glog.Errorf("no such user")
		return
	}
	if val == clientSecret {
		server.CToken[token] = clientID
	}
	fmt.Println(server.CStore[clientID])
	go server.DelayTokenExisting(10, token)
	fmt.Println(server.CStore[clientID])
	w.Write([]byte(token))
}

// Credentials sets a credentials for client.
func (server *WebServer) Credentials(w http.ResponseWriter, r *http.Request) {
	clientID := uuid.New().String()[:8]
	clientSecret := uuid.New().String()[:8]
	if val, ok := server.CStore[clientID]; ok || val != "" {
		glog.Errorf("ClientID problems or client exists")
		return
	}
	server.CStore[clientID] = clientSecret
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID": clientID, "CLIENT_SECRET": clientSecret})
}

// Validate validates a request for headers token,
func (server *WebServer) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := server.ValidateToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			glog.Errorf("Error: %v", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// DelayTokenExisting defines token life-existence time.
func (server *WebServer) DelayTokenExisting(exists int, token string) {
	time.Sleep(time.Duration(exists) * time.Second)
	delete(server.CStore, server.CToken[token])
	delete(server.CToken, token)
}
