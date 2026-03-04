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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

type CREOutput struct {
	Status string `json:"Status"`
	TxHash string `json:"TxHash"`
}

func main() {
	http.HandleFunc("/tokenize", handleTokenize)
	http.HandleFunc("/ws", handleWS)

	fmt.Println("🚀 Worker running on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleWS(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	fmt.Println("🔌 WebSocket client connected")

	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
		fmt.Println("❌ WebSocket client disconnected")
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func broadcast(message interface{}) {
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.Marshal(message)

	for client := range clients {
		client.WriteMessage(websocket.TextMessage, data)
	}
}

func handleTokenize(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("📦 Received payload")

	broadcast(map[string]string{"stage": "received"})

	// Marshal inline JSON
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode payload", 500)
		return
	}

	broadcast(map[string]string{"stage": "executing_cre"})

	start := time.Now()

	cmd := exec.Command(
		"cre",
		"workflow",
		"simulate",
		"./auto-lock-defi",
		"--target", "staging-settings",
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
	w.Write(output)
}

func extractSimulationJSON(output []byte) []byte {
	str := string(output)

	marker := "Workflow Simulation Result:"
	index := strings.Index(str, marker)
	if index == -1 {
		return nil
	}

	// Slice from marker onward
	sub := str[index:]

	start := strings.Index(sub, "{")
	end := strings.Index(sub, "}")

	if start == -1 || end == -1 {
		return nil
	}

	return []byte(sub[start : end+1])
}
