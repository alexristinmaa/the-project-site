package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var outFmt *os.File

type Post struct {
	Id        string
	Title     string
	Body      string
	Author_id int
	Date      int
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func individualPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("post"))
}

func allPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-type", "text/json")

	rows, err := db.Query("select Post_ID, Title, Body, Author_ID from posts")
	if err != nil {
		fmt.Fprint(outFmt, err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	ret := make([]Post, 0)

	defer rows.Close()
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author_id, &post.Date)

		if err != nil {
			fmt.Fprint(outFmt, err)
			http.Error(w, "Unknown error occured", http.StatusInternalServerError)
			return
		}

		ret = append(ret, post)
	}

	jsonData, err := json.Marshal(ret)

	if err != nil {
		fmt.Fprint(outFmt, err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func main() {
	router := httprouter.New()

	db, err := sql.Open("sqlite3", "./../foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	outFmt, err := os.Create("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outFmt.Close()

	router.NotFound = http.FileServer(http.Dir("frontend/dist")) // ugly solution

	// Overrides all file paths
	router.GET("/posts/:post", individualPost)
	router.GET("/posts", allPosts)

	fmt.Println("Serving at http://localhost:" + os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
