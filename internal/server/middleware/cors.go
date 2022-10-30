package middleware

import (
	"net/http"

	"github.com/rs/cors"
	"github.com/stevenhansel/csm-ending-prediction-be/internal/config"
)

func Cors(environment config.Environment) func(http.Handler) http.Handler {
	origins := []string{"*"}
	// if environment == config.PRODUCTION {
	// 	origins = append(origins, "https://dennis.dog")
	// } else {
	// 	origins = append(origins, "*")
	// }

	cors := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Access-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return cors.Handler
}
