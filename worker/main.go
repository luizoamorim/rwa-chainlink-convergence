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
// WEBSOCKET MANAGEMENT
//////////////////////////////////////////////////////////////

// Allow connections from any origin (development only)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Active WebSocket clients
var clients = make(map[*websocket.Conn]bool)
var wsMutex sync.Mutex

//////////////////////////////////////////////////////////////
// EXECUTION LOCK (Prevents multiple simultaneous executions)
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
// MAIN ENTRY POINT
//////////////////////////////////////////////////////////////

func main() {
	http.HandleFunc("/tokenize", handleTokenize)
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("🚀 Worker running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

//////////////////////////////////////////////////////////////
// WEBSOCKET HANDLER
//////////////////////////////////////////////////////////////

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	wsMutex.Lock()
	clients[conn] = true
	wsMutex.Unlock()

	fmt.Println("🔌 WebSocket client connected")

	defer func() {
		wsMutex.Lock()
		delete(clients, conn)
		wsMutex.Unlock()
		conn.Close()
		fmt.Println("❌ WebSocket client disconnected")
	}()

	// Keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

//////////////////////////////////////////////////////////////
// BROADCAST FUNCTION
//////////////////////////////////////////////////////////////

// Sends stage updates to all connected WebSocket clients
func broadcast(message interface{}) {
	wsMutex.Lock()
	defer wsMutex.Unlock()

	data, _ := json.Marshal(message)

	for client := range clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
}

//////////////////////////////////////////////////////////////
// TOKENIZATION HANDLER
//////////////////////////////////////////////////////////////

func handleTokenize(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Prevent concurrent executions
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

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	broadcast(map[string]string{"stage": "received"})

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode payload", 500)
		return
	}

	broadcast(map[string]string{"stage": "executing_cre"})

	start := time.Now()

	fmt.Println("📥 Received payload:")
	fmt.Println(string(bodyBytes))

	// Execute CRE workflow via CLI
	cmd := exec.Command(
		"cre",
		"workflow",
		"simulate",
		"./auto-lock-defi",
		"--target", "staging-settings",
		"--broadcast", // REAL TRANSACTION ENABLED
		"--trigger-index", "0",
		"--non-interactive",
		"--http-payload", string(bodyBytes),
	)

	cmd.Dir = ".."

	output, err := cmd.CombinedOutput()

	fmt.Println("📡 CRE output:")
	fmt.Println(string(output))

	if err != nil {
		broadcast(map[string]string{
			"stage": "error",
			"error": string(output),
		})
		http.Error(w, string(output), 500)
		return
	}

	elapsed := time.Since(start)
	fmt.Println("⏱ Execution time:", elapsed)

	// Extract JSON result from CLI output
	cleanJSON := extractSimulationJSON(output)
	if cleanJSON == nil {
		http.Error(w, "No simulation result found", 500)
		return
	}

	var result CREOutput
	if err := json.Unmarshal(cleanJSON, &result); err != nil {
		http.Error(w, "Failed to parse CRE result", 500)
		return
	}

	broadcast(map[string]string{
		"stage":  "success",
		"txHash": result.TxHash,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

//////////////////////////////////////////////////////////////
// CLI OUTPUT PARSER
//////////////////////////////////////////////////////////////

// Extracts JSON object from CRE CLI output
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
