package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	JsonRpc string    `json:"jsonrpc"`
	Id      int       `json:"id"`
	Result  string    `json:"result"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

	    url := "http://0.0.0.0:8545"

		payload := []byte(`{ "jsonrpc":"2.0", "method":"eth_blockNumber","params":[],"id":1}`)

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

		result, err := strconv.ParseInt(responseData.Result, 0, 0)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		
		fmt.Fprintf(w, "mantle_latest_block_height %d\n", result)
	})

	fmt.Println("Server listening on port 8991 for GET requests...")
	http.ListenAndServe(":8991", nil)
}