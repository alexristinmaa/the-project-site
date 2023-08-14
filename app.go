package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
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
	AuthorName  string
	Date        time.Time
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	// If it has extension -> Handle it as a file, else return index.html
	cleanUrl := filepath.Clean(r.URL.Path)
	extension := filepath.Ext(cleanUrl)

	var file []byte
	var err error

	if extension == "" {
		file, err = os.ReadFile("./frontend/dist/index.html")
	} else {
		file, err = os.ReadFile(filepath.Join("./frontend/dist/" + cleanUrl))
	}

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Error reading file %s", cleanUrl)
		return
	}

	w.Header().Add("Content-Type", mime.TypeByExtension(extension))
	w.Write(file)
}

func individualPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("postId"))
}

func allPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-type", "text/json")

	rows, err := conn.Query(context.Background(), "select post_id, title, body, author, author_name, description from posts")
	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	ret := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author, &post.AuthorName, &post.Description)

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

	port := os.Getenv("PORT")
	loggingLocation := os.Getenv("LOG_LOC")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePass := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv(("DATABASE_NAME"))
	databaseIP := os.Getenv("DATABASE_IP")

	if port == "" {
		panic("No port supplied. Add to environment variables")
	}

	if databaseUser == "" || databasePass == "" || databaseName == "" || databaseIP == "" {
		panic("One of: DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME or DATABASE_IP not supplied: add into environment variables (.bashrc)")
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
	databaseUser = url.QueryEscape(databaseUser)
	databasePass = url.QueryEscape(databasePass)
	databaseName = url.QueryEscape(databaseName)
	databaseIP = url.QueryEscape(databaseIP)

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", databaseUser, databasePass, databaseIP, databaseName)

	conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	router := httprouter.New()

	router.GET("/api/posts/:postId", individualPost)
	router.GET("/api/posts", allPosts)

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)

	mux.Handle("/api/", router)

	fmt.Println("Serving at http://localhost:" + port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
