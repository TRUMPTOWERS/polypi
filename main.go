package main

import (
	"log"
	"net/http"
	"time"

	"gopkg.in/redis.v4"

	"github.com/TRUMPTOWERS/polypi/pie/piehandler"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	index := http.FileServer(NewOneFile("index.html"))
	r.Methods("GET").Path("/").Handler(index)

	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	r.Methods("GET").Path("/static/").Handler(static)

	pie := piehandler.Handler{DS: client}
	r.Methods("GET").Path("/pie/{id:[0-9]+}.json").Handler(pie).Name("pie")

	var purchase http.Handler
	r.Methods("POST").Path("/pie/{id:[0-9]+}/purchase").Handler(purchase).Name("purchase")

	var recommend http.Handler
	r.Methods("GET").Path("/pies/recommend").Handler(recommend).Name("recommend")

	log.Fatal(s.ListenAndServe())
}

// OneFile is like http.Dir, but always serves the same file
type OneFile struct {
	file string
	dir  http.FileSystem
}

// Open impliments the http.FileSystem interface
func (of OneFile) Open(_ string) (http.File, error) {
	return of.dir.Open(of.file)
}

// NewOneFile creates a OneFile
func NewOneFile(name string) OneFile {
	return OneFile{
		file: name,
		dir:  http.Dir(""),
	}
}
