package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "wrong method!", http.StatusMethodNotAllowed)
	}
	fmt.Println("got it test!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "successfully",
	})

}
