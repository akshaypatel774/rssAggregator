package main

import (
	// "fmt"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akshaypatel774/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// abc, err := urlToFeed("https://wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(abc)

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found in env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to the database:", err)
	}
	
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	go startScrapping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1router := chi.NewRouter()
	v1router.Get("/healthz", handlerReadiness)
	v1router.Get("/err", handlerErr)

	v1router.Post("/users", apiCfg.handlerCreateUser)
	v1router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	
	v1router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1router.Get("/feeds", apiCfg.handlerGetFeeds)
	
	v1router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
	
	v1router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	router.Mount("/v1", v1router)

	server := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	log.Printf("Server started on port: %v", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}