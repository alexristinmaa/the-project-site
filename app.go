package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Post struct {
	Id          string
	Title       string
	Body        string
	Description string
	Author_id   int
	Date        time.Time
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func individualPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("post"))
}

func allPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-type", "text/json")

	rows, err := db.Query("select Post_ID, Title, Body, Author_ID, Date, Description from posts")
	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	ret := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author_id, &post.Date, &post.Description)

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

	port := "3000"
	dbPath := "all.db"
	loggingLocation := "application_logs/error.log"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if os.Getenv("DB_PATH") != "" {
		dbPath = os.Getenv("DB_PATH")
	}

	if os.Getenv("LOG_LOC") != "" {
		loggingLocation = os.Getenv("LOG_LOC")
	}

	f, err := os.OpenFile(loggingLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router.NotFound = http.FileServer(http.Dir("frontend/dist")) // ugly solution

	// Overrides all file paths
	router.GET("/posts/:post", individualPost)
	router.GET("/posts", allPosts)

	fmt.Println("Serving at http://localhost:" + port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
