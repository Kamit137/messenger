package main

import (
	"log"
	"net/http"
	"os"
	"messenger/internal/jwt"
	"messenger/internal/auth"
	"messenger/internal/storage"
	"messenger/web"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	addr := getEnv("HTTP_ADDR", ":8080")
	dbPath := getEnv("DATABASE_PATH", "./messenger.db")
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	if err := storage.NewSQLiteStorage(dbPath); err != nil {
		log.Fatalf("DB init: %v", err)
	}
	defer storage.Close()

	jwt.InitJWT(secret)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", web.Index)
	mux.HandleFunc("GET /login", web.LoginPage)
	mux.HandleFunc("GET /register", web.RegisterPage)
	
	mux.HandleFunc("POST /api/register", auth.Register)
	mux.HandleFunc("POST /api/login", auth.Login)
	mux.HandleFunc("POST /api/logout", auth.Logout)
	mux.HandleFunc("GET /api/me", auth.Me)
	
	mux.Handle("GET /static/", http.StripPrefix("/static/", web.StaticHandler()))


	log.Printf("Server running on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}


func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
