package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gopkg.in/ahmetb/go-linq.v3"
)

type Server struct {
	config *Config
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Run() {
	log.Fatal(
		http.ListenAndServeTLS(
			fmt.Sprint(":", server.config.Host.Port),
			server.config.Host.Ssl.Cert,
			server.config.Host.Ssl.Key,
			getRouter()))
}

func getRouter() *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle)

	return myRouter
}

func homePage(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(responseWriter, "This is the home Page bro")
	fmt.Println("Endpoint hit: homePage")
}

func returnAllArticles(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint hit: returnAllArticles")
	json.NewEncoder(responseWriter).Encode(Items)
}

func returnSingleArticle(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, _ := strconv.Atoi(vars["id"])

	article := linq.From(Items).Where(func(c interface{}) bool {
		return c.(Item).Id == key
	}).First()

	fmt.Fprintf(responseWriter, "\t%+v\n", article)
}
