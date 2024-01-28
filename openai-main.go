package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/rs/cors"
)

const apiKey = "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"
const apiEndpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	// HTTP Handlers

	// Establish database connection
	connectionStr := "postgres://postgres:Jaden93014124!@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Println("Error connecting to PostgreSQL:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/generate-research", func(w http.ResponseWriter, r *http.Request) { generateResearchHandler(w, r, db) })

	cor_var := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		Debug:          true,
	})

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}
	handler := cor_var.Handler(mux)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: handler,
	}

	go func() {
		fmt.Println("Server is listening on", serverAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server gracefully stopped")
}

func generateResearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	topic := strings.TrimSpace(r.FormValue("topic"))
	if topic == "" {
		http.Error(w, "Topic parameter is required", http.StatusBadRequest)
		return
	}

	generatedQuestions := generateQuestions(topic)

	var response struct {
		Status    string   `json:"status"`
		Message   string   `json:"message"`
		Questions []string `json:"questions"`
		Articles  []string `json:"articles"`
	}

	for _, question := range generatedQuestions {
		articles := getArticles(topic)

		// Insert userResponse into the PostgreSQL

		err := insertUserAnswer(topic, question, articles, db)
		if err != nil {
			fmt.Println("Error inserting user answer:", err)
			response.Status = "error"
			response.Message = "Internal Server Error"
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response.Questions = append(response.Questions, question)
		response.Articles = append(response.Articles, articles...)
	}

	response.Status = "success"
	response.Message = "Research questions generated and answers inserted into the database"

	responseJSON, err := json.MarshalIndent(response, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(responseJSON))
}

// Create a new function to insert UserResponse into database
// Parameters: The question and userResponse, alongside whatever database used

func insertUserAnswer(topic string, question string, articles []string, db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS research_assistant (
		id SERIAL PRIMARY KEY,
		topic TEXT,
		question TEXT,
		articles TEXT[]
	);
	`
	_, panic := db.Exec(createTableQuery)
	if panic != nil {
		return panic
	}

	query := `INSERT INTO research_assistant (topic, question, articles) VALUES ($1, $2, $3) RETURNING id`
	var userID int
	err := db.QueryRow(query, topic, question, pq.Array(articles)).Scan(&userID)
	if err != nil {
		return err
	}

	return nil
}
