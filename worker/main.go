package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	fmt.Println("📦 Received payload:", payload)

	broadcast(map[string]string{
		"stage": "received",
	})

	filePath := "worker_payload.json"

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create payload file", 500)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(payload); err != nil {
		http.Error(w, "Failed to write payload file", 500)
		return
	}

	broadcast(map[string]string{
		"stage": "executing_cre",
	})

	start := time.Now()

	cmd := exec.Command(
		"cre",
		"workflow",
		"simulate",
		"./auto-lock-defi",
		"--target", "staging-settings",
		"--broadcast",
		"--trigger-index", "0",
		"--non-interactive",
		"--http-payload", filePath,
	)

	// executa a partir da raiz do projeto
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

	// Extrair JSON do output
	cleanJSON := extractJSON(output)

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
	w.Write(cleanJSON)
}

func extractJSON(output []byte) []byte {
	str := string(output)

	start := strings.Index(str, "{")
	end := strings.LastIndex(str, "}")

	if start == -1 || end == -1 {
		return output
	}

	return []byte(str[start : end+1])
}
