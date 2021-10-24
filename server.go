package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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
			server.getRouter()))
}

func (server *Server) getRouter() *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)

	jwtMiddleware := server.getJwtMiddleware()

	myRouter.HandleFunc("/", homePage)
	addRouteWithMiddleware("/items", returnAllItems, myRouter, jwtMiddleware.HandlerWithNext)
	addRouteWithMiddleware("/items/{id}", returnSingleItem, myRouter, jwtMiddleware.HandlerWithNext)

	return myRouter
}

func addRouteWithMiddleware(path string, function http.HandlerFunc, router *mux.Router, middleware negroni.HandlerFunc) {
	router.Handle(path, negroni.New(
		negroni.HandlerFunc(middleware),
		negroni.Wrap(http.HandlerFunc(function)),
	))
}

func (server *Server) getJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(
		jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(server.config.Authentication.Jwt.Secret), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})
}

func homePage(responseWriter http.ResponseWriter, request *http.Request) {
	fmt.Fprint(responseWriter, "This is the home Page bro")
}

func returnAllItems(responseWriter http.ResponseWriter, request *http.Request) {
	writeJson(responseWriter, Items)
}

func returnSingleItem(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	key, _ := strconv.Atoi(vars["id"])

	item := linq.From(Items).Where(func(c interface{}) bool {
		return c.(Item).Id == key
	}).First()

	if item == nil {
		writeError(responseWriter, "Item not found")
		return
	}

	writeJson(responseWriter, item)
}

func writeError(responseWriter http.ResponseWriter, message string) {
	responseWriter.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(responseWriter, message)
}

func writeJson(responseWriter http.ResponseWriter, input interface{}) {
	json.NewEncoder(responseWriter).Encode(input)
}
