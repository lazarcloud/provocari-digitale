package utils

import (
	"github.com/rs/cors"
)

var CORSHandler = cors.New(cors.Options{
	AllowedOrigins:   []string{"http://localhost:5173"},
	AllowCredentials: true,
	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
})
