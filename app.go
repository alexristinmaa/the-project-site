package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var sourceMap map[string]string

func initSourceMap() error {
	// Loop through the public folder and add folders to the source map
	sourceMap = make(map[string]string)

	err := filepath.Walk("public", func(pathToFile string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", pathToFile, err)
			return err
		}

		if !info.IsDir() && path.Base(pathToFile) == "index.html" {
			sourceMap[path.Dir(pathToFile[6:])] = "true" // [6:] is to remove public/ from the path
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func AssertAndHandle(reqPath string, handler func(http.ResponseWriter, *http.Request)) {
	// Check if request is a valid path (with an index.html file) or if it has a handler

	http.HandleFunc(reqPath, func(w http.ResponseWriter, r *http.Request) {
		cleanURL := path.Clean(r.URL.Path)
		if _, ok := sourceMap[cleanURL]; ok {
			toServe, err := os.ReadFile(path.Join("public", cleanURL, "index.html"))
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Something Went Wrong: corrupted file?", http.StatusInternalServerError)
			}

			w.Write(toServe)
			return
		} else if cleanURL != reqPath {
			http.Error(w, "Could not find anything here... are you lost?", http.StatusNotFound)
			return
		}

		handler(w, r)
	})
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func base(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "This is our main page!\n")
}

func main() {
	err := initSourceMap()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(sourceMap)

	AssertAndHandle("/hello", hello)
	AssertAndHandle("/headers", headers)
	AssertAndHandle("/", base)

	fmt.Println("Server started at: http://localhost:" + os.Getenv("PORT"))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
