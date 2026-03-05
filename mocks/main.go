package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type CarData struct {
	Plate     string  `json:"plate"`
	Status    string  `json:"status"` // "clear", "blocked", "stolen"
	Fines     float64 `json:"fines"`
	ModelCode string  `json:"model_code"`
	Price     float64 `json:"price"`
}

func main() {
	http.HandleFunc("/detran/", func(w http.ResponseWriter, r *http.Request) {
		plate := strings.TrimPrefix(r.URL.Path, "/detran/")
		// Mock Data
		data := CarData{
			Plate:     plate,
			Status:    "clear",
			Fines:     0.0,
			ModelCode: "005456-9",
			Price:     75000.00,
		}

		if plate == "ABC1234" {
			data.Price = 35000.00
		}

		if plate == "XYZ9999" {
			data.Price = 120000.00
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}
