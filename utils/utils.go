package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func SendJSONResponse(w http.ResponseWriter, data any) {
	jsonBytes, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Error in SendJSONResponse", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func GenerateID() string {
	return uuid.NewString()
}
