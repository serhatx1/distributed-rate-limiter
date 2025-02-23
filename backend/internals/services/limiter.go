package services

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internals/repository"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetLimit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var setLimitReq ChangeLimitRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&setLimitReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	validate := validator.New()
	if err := validate.Struct(setLimitReq); err != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}
	client := repository.Getclient()
	redisKey := fmt.Sprintf("ratelimit:%s", setLimitReq.EndPoint)

	err := client.Set(redisKey, setLimitReq.Ratelimit, 24*time.Hour).Err()
	if err != nil {
		fmt.Printf("Failed to set Redis for endpoint %s: %v", setLimitReq.EndPoint, err)
	}

	filter := bson.M{"endpoint": bson.M{"$eq": setLimitReq.EndPoint}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := repository.GetRateLimiterCollection()

	var existingEndpoint struct{}
	err = collection.FindOne(ctx, filter).Decode(&existingEndpoint)
	if err == mongo.ErrNoDocuments {
		insert := bson.M{
			"endpoint":  setLimitReq.EndPoint,
			"ratelimit": setLimitReq.Ratelimit,
		}

		_, err := collection.InsertOne(ctx, insert)
		if err != nil {
			http.Error(w, "Failed to create new endpoint", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "New rate limit set successfully",
		})

	} else if err != nil {
		http.Error(w, "Error while checking endpoint", http.StatusInternalServerError)
		return
	} else {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Rate limit already exist",
		})
	}
}
