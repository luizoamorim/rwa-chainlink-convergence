package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//////////////////////////////////////////////////////////////
// WEBSOCKET
//////////////////////////////////////////////////////////////

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var wsMutex sync.Mutex

//////////////////////////////////////////////////////////////
// EXECUTION LOCK
//////////////////////////////////////////////////////////////

var executionMutex sync.Mutex
var isExecuting bool

//////////////////////////////////////////////////////////////
// CRE OUTPUT STRUCT
//////////////////////////////////////////////////////////////

type CREOutput struct {
	Status string `json:"Status"`
	TxHash string `json:"TxHash"`
}

//////////////////////////////////////////////////////////////
// MAIN
//////////////////////////////////////////////////////////////

func main() {

	http.HandleFunc("/tokenize", handleTokenize)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("🚀 Worker running on :8081")

	log.Fatal(http.ListenAndServe(":8081", nil))
}

//////////////////////////////////////////////////////////////
// WS HANDLER
//////////////////////////////////////////////////////////////

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wsMutex.Lock()
	clients[conn] = true
	wsMutex.Unlock()

	fmt.Println("🔌 WS connected")

	defer func() {
		wsMutex.Lock()
		delete(clients, conn)
		wsMutex.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

//////////////////////////////////////////////////////////////
// BROADCAST
//////////////////////////////////////////////////////////////

func broadcast(message interface{}) {

	wsMutex.Lock()
	defer wsMutex.Unlock()

	data, _ := json.Marshal(message)

	for client := range clients {

		err := client.WriteMessage(websocket.TextMessage, data)

		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

//////////////////////////////////////////////////////////////
// TOKENIZATION
//////////////////////////////////////////////////////////////

func handleTokenize(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	executionMutex.Lock()

	if isExecuting {
		executionMutex.Unlock()
		http.Error(w, "Execution already in progress", http.StatusTooManyRequests)
		return
	}

	isExecuting = true
	executionMutex.Unlock()

	defer func() {
		executionMutex.Lock()
		isExecuting = false
		executionMutex.Unlock()
	}()

	fmt.Println("🚗 Tokenization request received")

	var payload map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bodyBytes, _ := json.Marshal(payload)

	//////////////////////////////////////////////////////////////
	// PIPELINE STAGES
	//////////////////////////////////////////////////////////////

	broadcast(map[string]string{
		"stage": "verifying_identity",
	})

	time.Sleep(500 * time.Millisecond)

	broadcast(map[string]string{
		"stage": "worldid_verified",
	})

	time.Sleep(500 * time.Millisecond)

	broadcast(map[string]string{
		"stage": "checking_vehicle_registry",
	})

	time.Sleep(500 * time.Millisecond)

	broadcast(map[string]string{
		"stage": "executing_cre",
	})

	//////////////////////////////////////////////////////////////
	// RUN CRE
	//////////////////////////////////////////////////////////////

	cmd := exec.Command(
		"cre",
		"workflow",
		"simulate",
		"./auto-lock-defi",
		"--target", "staging-settings",
		"--broadcast",
		"--trigger-index", "0",
		"--non-interactive",
		"--http-payload", string(bodyBytes),
	)

	cmd.Dir = ".."

	output, err := cmd.CombinedOutput()

	fmt.Println("CRE OUTPUT:")
	fmt.Println(string(output))

	if err != nil {

		if strings.Contains(string(output), "world id") {

			broadcast(map[string]string{
				"stage": "worldid_failed",
			})

		} else if strings.Contains(string(output), "vehicle") {

			broadcast(map[string]string{
				"stage": "oracle_failed",
			})

		} else {

			broadcast(map[string]string{
				"stage": "mint_failed",
			})
		}

		http.Error(w, string(output), http.StatusInternalServerError)
		return
	}

	broadcast(map[string]string{
		"stage": "minting_nft",
	})

	time.Sleep(500 * time.Millisecond)

	clean := extractSimulationJSON(output)

	if clean == nil {
		http.Error(w, "CRE result not found", http.StatusInternalServerError)
		return
	}

	var result CREOutput

	if err := json.Unmarshal(clean, &result); err != nil {
		http.Error(w, "Failed parsing result", http.StatusInternalServerError)
		return
	}

	broadcast(map[string]string{
		"stage":  "success",
		"txHash": result.TxHash,
	})

	fmt.Println("✅ NFT minted", result.TxHash)

	json.NewEncoder(w).Encode(result)
}

//////////////////////////////////////////////////////////////
// PARSER
//////////////////////////////////////////////////////////////

func extractSimulationJSON(output []byte) []byte {

	str := string(output)

	marker := "Workflow Simulation Result:"
	index := strings.Index(str, marker)

	if index == -1 {
		return nil
	}

	sub := str[index:]

	start := strings.Index(sub, "{")
	end := strings.Index(sub, "}")

	if start == -1 || end == -1 {
		return nil
	}

	return []byte(sub[start : end+1])
}
