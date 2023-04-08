package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Alert struct {
	CreatedAt   string `json:"createdAt"`
	Name        string `json:"name"`
	Protocol    string `json:"protocol"`
	FindingType string `json:"findingType"`
	Source      struct {
		TransactionHash string `json:"transactionHash"`
		Block           struct {
			Number  int `json:"number"`
			ChainId int `json:"chainId"`
		} `json:"block"`
		Bot struct {
			Id string `json:"id"`
		} `json:"bot"`
	} `json:"source"`
	Severity string `json:"severity"`
	Hash     string `json:"hash"`
}

type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
	EndCursor   struct {
		AlertId     string `json:"alertId"`
		BlockNumber int    `json:"blockNumber"`
	} `json:"endCursor"`
}

type AlertData struct {
	PageInfo PageInfo `json:"pageInfo"`
	Alerts   []Alert  `json:"alerts"`
}

type AlertResponse struct {
	Alerts AlertData `json:"alerts"`
}

type Response struct {
	Data AlertResponse `json:"data"`
}

func main() {
	http.HandleFunc("/graphql", authMiddleware(dataHandler))
	log.Fatal(http.ListenAndServe(":42000", nil))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implement your authentication logic here
		isAuthenticated := true

		if !isAuthenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// create sample data
	alerts := []Alert{
		{
			CreatedAt:   time.Now().Format(time.RFC3339),
			Name:        "src: 0x3f88c2b3e267e6b8e9dE017cdB47a59aC9Ecb284",
			Hash:        randomHash(),
			Protocol:    "ethereum",
			FindingType: "INFORMATION",
			Source: struct {
				TransactionHash string `json:"transactionHash"`
				Block           struct {
					Number  int `json:"number"`
					ChainId int `json:"chainId"`
				} `json:"block"`
				Bot struct {
					Id string `json:"id"`
				} `json:"bot"`
			}{
				TransactionHash: "",
				Block: struct {
					Number  int `json:"number"`
					ChainId int `json:"chainId"`
				}{
					Number:  16994597,
					ChainId: 1,
				},
				Bot: struct {
					Id string `json:"id"`
				}{
					Id: "0xbe1872858e63b6ed4ef7b84fc453970dc8d89968715797662a4f43c01d598aab",
				},
			},
			Severity: "LOW",
		},
		// ... add more alerts here
	}

	// create response
	response := Response{
		Data: AlertResponse{
			Alerts: AlertData{
				PageInfo: PageInfo{
					HasNextPage: false,
				},
				Alerts: alerts,
			},
		},
	}

	resp, _ := json.Marshal(response)
	_, err := fmt.Fprint(w, string(resp))
	if err != nil {
		return
	}
}

func randomHash() string {
	hashBytes := make([]byte, 32)
	_, err := rand.Read(hashBytes)
	if err != nil {
		panic(err)
	}
	return "0x" + hex.EncodeToString(hashBytes)
}
