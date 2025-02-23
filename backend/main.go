package main

import (
	"encoding/json"
	"log"
	"main/internals/repository"
	"main/internals/services"
	"main/pkg/rabbitmq"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type Limiter struct {
	EndPoint         string   `bson:"endpoint"`
	Ratelimit        int      `bson:"ratelimit"`
	Exclusives       []string `bson:"exclusives"`
	IsCustom         bool     `bson:"is_custom"`
	DefaultRateLimit int      `bson:"defaultratelimit"`
}

type RouteInfo struct {
	Path      string `json:"path"`
	RateLimit int    `json:"rate_limit"`
	Source    string `json:"source"` // "Redis", "MongoDB", or "Default"
}

var registeredRoutes []string

func registerRoute(mux *http.ServeMux, pattern string, handler http.Handler) {
	registeredRoutes = append(registeredRoutes, pattern)
	mux.Handle(pattern, handler)
}

func listRoutes(w http.ResponseWriter, r *http.Request) {
	var routes []RouteInfo

	for _, route := range registeredRoutes {
		defaultRateLimit := 1000
		routeInfo := RouteInfo{
			Path:      route,
			RateLimit: defaultRateLimit,
			Source:    "Default",
		}

		// First try Redis
		rateLimitValue, err := repository.GetPath("ratelimit:" + route)
		if err == nil && rateLimitValue != "" {
			rateLimitCache, _ := strconv.Atoi(rateLimitValue)
			routeInfo.RateLimit = rateLimitCache
			routeInfo.Source = "Redis"
		} else {
			// If Redis fails or empty, check MongoDB
			doc, err := repository.GetDocuments(route)
			if err == nil && doc != nil {
				if rateLimit, ok := (*doc)["ratelimit"]; ok {
					switch v := rateLimit.(type) {
					case int32:
						routeInfo.RateLimit = int(v)
						routeInfo.Source = "MongoDB"
					case int64:
						routeInfo.RateLimit = int(v)
						routeInfo.Source = "MongoDB"
					case float64:
						routeInfo.RateLimit = int(v)
						routeInfo.Source = "MongoDB"
					}
				}
			}
		}
		routes = append(routes, routeInfo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"routes": routes,
		"count":  len(routes),
	})
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	repository.InitRedis()
	repository.InitMongo(os.Getenv("URI"), os.Getenv("MONGO_DB"))
	rabbitmq.RabbitMQInit()
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://example.com"}, // Allow specific origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},                 // Allowed HTTP methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"},                // Allowed headers
		AllowCredentials: true,                                                     // Allow credentials (cookies, auth headers)
	}
	mux := http.NewServeMux()

	registerRoute(mux, "/ratelimit/changelimit", RateLimitMiddleware(http.HandlerFunc(services.ChangeLimit)))
	registerRoute(mux, "/ratelimit/setlimit", RateLimitMiddleware(http.HandlerFunc(services.SetLimit)))
	registerRoute(mux, "/listroutes", RateLimitMiddleware(http.HandlerFunc(listRoutes)))
	registerRoute(mux, "/test/test", RateLimitMiddleware(http.HandlerFunc(services.Test)))

	corsHandler := cors.New(corsOptions)
	handler := corsHandler.Handler(mux)

	err = http.ListenAndServe(":3000", handler)
	if err != nil {
		log.Panic("Error creating new server.", err)
	}

}
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defaultRateLimit := 1000
		startTime := time.Now()
		clientIP := r.RemoteAddr
		path := r.URL.Path

		log.Println("➡️ Processing request:", path, "from", clientIP)

		visitCount, err := repository.TrackIPForPath(path, clientIP)
		if err != nil {
			log.Println("❌ Error tracking IP in Redis:", err)
		} else {
			log.Println("✅ Retrieved visit count from Redis:", visitCount)
		}

		rateLimitValue, err := repository.GetPath("ratelimit:" + path)
		if err != nil {
			log.Println("⚠️ Error retrieving rate limit from Redis, falling back to MongoDB:", err)
		} else if rateLimitValue != "" {
			log.Println("✅ Retrieved rate limit from Redis cache:", rateLimitValue)
			rateLimitCache, _ := strconv.Atoi(rateLimitValue)
			defaultRateLimit = rateLimitCache
		} else {
			log.Println("⚠️ No rate limit found in Redis, checking MongoDB...")
		}

		if rateLimitValue == "" || err != nil {
			doc, err := repository.GetDocuments(path)
			if err != nil {
				log.Println("❌ Error retrieving rate limit from MongoDB:", err)
			} else {
				log.Println("✅ Retrieved rate limit document from MongoDB:", doc)
				rateLimit, ok := (*doc)["ratelimit"]
				if ok {
					switch v := rateLimit.(type) {
					case int32:
						defaultRateLimit = int(v)
						log.Println("✅ Rate limit (int32) from MongoDB:", defaultRateLimit)
					case int64:
						defaultRateLimit = int(v)
						log.Println("✅ Rate limit (int64) from MongoDB:", defaultRateLimit)
					case float64:
						defaultRateLimit = int(v)
						log.Println("✅ Rate limit (float64) from MongoDB:", defaultRateLimit)
					default:
						log.Println("⚠️ Warning: Unsupported rateLimit type:", reflect.TypeOf(v))
					}
				} else {
					log.Println("⚠️ No 'ratelimit' field found in MongoDB document")
				}
			}
		}

		log.Printf("📢 Request: %s %s | Visits: %d | Limit: %d", r.Method, r.URL.Path, visitCount, defaultRateLimit)

		if visitCount <= defaultRateLimit {
			log.Println("✅ Request allowed.")
			next.ServeHTTP(w, r)
		} else {
			log.Println("❌ Rate limit exceeded!")
			http.Error(w, "Rate limit exceeded!", http.StatusTooManyRequests)
		}

		log.Printf("⏱ Completed request in %v", time.Since(startTime))
	})
}
