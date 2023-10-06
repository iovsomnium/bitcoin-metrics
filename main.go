package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Result int         `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

	    url := "http://127.0.0.1:8332"

		payload := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "getblockcount", "params": []}`)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			http.Error(w, "Error making POST request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var responseData Response
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&responseData)
		if err != nil {
			http.Error(w, "Error decoding JSON resp", http.StatusInternalServerError)
			return
		}

		result := responseData.Result
		fmt.Fprintf(w, "bitcoin_latest_block_height %d\n", result)
	})

	fmt.Println("Server listening on port 8991 for GET requests...")
	http.ListenAndServe(":8991", nil)
}
