package services

import (
	"encoding/json"
	"fmt"
	"log"
	"main/internals/repository"
	"main/pkg/rabbitmq"
	"net/http"
	"time"

	"github.com/go-playground/validator"
)

type ChangeLimitRequest struct {
	EndPoint  string `json:"endpoint" validate:"required"`
	Ratelimit int    `json:"ratelimit" validate:"required,min=1"`
}

func ChangeLimit(w http.ResponseWriter, r *http.Request) {
	// Check the request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var changeLimitReq ChangeLimitRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&changeLimitReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate the request body
	validate := validator.New()
	if err := validate.Struct(changeLimitReq); err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	client := repository.Getclient()
	redisKey := fmt.Sprintf("ratelimit:%s", changeLimitReq.EndPoint)

	err := client.Set(redisKey, changeLimitReq.Ratelimit, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to update Redis for endpoint %s: %v", changeLimitReq.EndPoint, err)
		http.Error(w, "Failed to update rate limit in Redis", http.StatusInternalServerError)
		return
	}

	messageBody, err := json.Marshal(changeLimitReq)
	if err != nil {
		log.Printf("Failed to marshal JSON for endpoint %s: %v", changeLimitReq.EndPoint, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ch := rabbitmq.Channel()
	err = rabbitmq.PublishMessage("setlimit", ch, messageBody)

	if err != nil {
		log.Printf("Failed to push message to RabbitMQ for endpoint %s: %v", changeLimitReq.EndPoint, err)
		http.Error(w, "Failed to push message to RabbitMQ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Rate limit updated successfully",
	})
}
