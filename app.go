package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
)

var conn *pgx.Conn

type Post struct {
	Id          string
	Title       string
	Body        string
	Description string
	Author      string
	Date        time.Time
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func individualPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("Thats not goods!?")
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("post"))
}

func allPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-type", "text/json")

	rows, err := conn.Query(context.Background(), "select post_id, title, body, author, description from posts")
	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	ret := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author, &post.Description)

		if err != nil {
			log.Println(err)
			http.Error(w, "Unknown error occured", http.StatusInternalServerError)
			return
		}

		ret = append(ret, post)
	}

	jsonData, err := json.Marshal(ret)

	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func main() {
	var err error
	router := httprouter.New()

	port := os.Getenv("PORT")
	loggingLocation := os.Getenv("LOG_LOC")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePass := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv(("DATABASE_NAME"))

	if port == "" {
		panic("No port supplied. Add to environment variables")
	}

	if databaseUser == "" || databasePass == "" || databaseName == "" {
		panic("One of: DATABASE_USER, DATABASE_PASSWORD or DATABASE_NAME not supplied: add into environment variables (.bashrc)")
	}

	if loggingLocation == "" {
		fmt.Println("No logging location supplied: Defaulting to logs/error.log")
	}

	dirPath := filepath.Dir(loggingLocation)
	err = os.MkdirAll(dirPath, os.ModePerm)

	if err != nil {
		log.Fatalf("error creating folder for file: %v", err)
	}

	f, err := os.OpenFile(loggingLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s", databaseUser, databasePass, databaseName)

	conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	router.NotFound = http.FileServer(http.Dir("frontend/dist")) // ugly solution

	// Overrides all file paths
	router.GET("/posts/:post", individualPost)
	router.GET("/posts", allPosts)

	fmt.Println("Serving at http://localhost:" + port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
