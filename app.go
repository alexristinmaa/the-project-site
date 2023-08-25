package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

type Post struct {
	Id          string
	Title       string
	Body        string
	Description string
	Author      string
	AuthorName  string
	Tags        []string
	Date        time.Time
}

type Tag struct {
	Name  string
	Count int
}

type PostResponse struct {
	Posts []Post
	Pages int
}

type SearchArguments struct {
	Page   int      `json:"page"`
	Tags   []string `json:"tags"`
	Search string   `json:"search"`
}

type PostArguments struct {
	Id string `json:"id"`
}

var conn *pgx.Conn

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
	searchArguments := PostArguments{}
	body, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &searchArguments)

	if err != nil || searchArguments.Id == "" {
		log.Println(err)
		http.Error(w, "Incorrect data supplied", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	post := Post{}
	var row pgx.Row

	row = conn.QueryRow(context.Background(), "SELECT title, body, author, author_name, tags, date FROM posts WHERE post_id = $1", searchArguments.Id)
	err = row.Scan(&post.Title, &post.Body, &post.Author, &post.AuthorName, &post.Tags, &post.Date)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add("Content-type", "application/json")
			fmt.Fprintf(w, "{\"error\":\"No post with id \\\"%s\\\"\"}", searchArguments.Id)
			return
		}

		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(post)

	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(jsonData)
}

func allPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	searchArguments := SearchArguments{}
	body, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(body, &searchArguments)

	if err != nil {
		log.Println(err)
		http.Error(w, "Incorrect data supplied", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	var rows pgx.Rows
	nPosts := 4
	offset := (searchArguments.Page - 1) * 4
	isRaw := len(searchArguments.Tags) == 0 && searchArguments.Search == ""

	if searchArguments.Page == 1 && isRaw {
		nPosts = 3
	} else if isRaw {
		offset--
	}

	if len(searchArguments.Tags) == 0 {
		rows, err = conn.Query(context.Background(), "select post_id, title, body, author, author_name, description, tags, date from posts limit $1 offset $2", nPosts, offset)
	} else {
		rows, err = conn.Query(context.Background(), "SELECT post_id, title, body, author, author_name, description, tags, date FROM posts WHERE posts.tags::text[] && $1 limit $2 offset $3", searchArguments.Tags, nPosts, offset)
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	posts := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author, &post.AuthorName, &post.Description, &post.Tags, &post.Date)

		if err != nil {
			log.Println(err)
			http.Error(w, "Unknown error occured", http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
	}

	var pages int
	var row pgx.Row

	if len(searchArguments.Tags) == 0 {
		row = conn.QueryRow(context.Background(), "select count(*) from posts")
	} else {
		row = conn.QueryRow(context.Background(), "select count(*) from posts where tags && $1", searchArguments.Tags)
	}

	err = row.Scan(&pages)

	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	if isRaw {
		pages = pages/4 + 1
	} else {
		pages = (pages-1)/4 + 1
	}

	res := PostResponse{posts, pages}

	jsonData, err := json.Marshal(res)

	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.Write(jsonData)
}

func getTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-type", "text/json")

	rows, err := conn.Query(context.Background(), "SELECT tag, COUNT(tag) FROM (SELECT UNNEST(tags) AS tag FROM posts) tag GROUP BY tag")

	if err != nil {
		log.Println(err)
		http.Error(w, "Unknown error occured", http.StatusInternalServerError)
		return
	}

	ret := make([]Tag, 0)

	for rows.Next() {
		tag := Tag{}

		err = rows.Scan(&tag.Name, &tag.Count)

		if err != nil {
			log.Println(err)
			http.Error(w, "Unknown error occured", http.StatusInternalServerError)
			return
		}

		ret = append(ret, tag)
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

	router.GET("/api/getTags", getTags)

	router.POST("/api/post", individualPost)
	router.POST("/api/posts", allPosts)

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)

	mux.Handle("/api/", router)

	fmt.Println("Serving at http://localhost:" + port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
