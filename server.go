package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ahmetb/go-linq/v3"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Server struct {
	config *Config
}

func NewServer(config *Config) *Server {
	return &Server{
		config: config,
	}
}

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

func (server *Server) Run() {
	ginRouter := gin.Default()
	ginRouter.GET("/", homePage)

	authMiddleware := getJwtMiddleware()

	ginRouter.POST("/login", authMiddleware.LoginHandler)

	ginRouter.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := ginRouter.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
	}

	log.Fatal(
		ginRouter.RunTLS(
			fmt.Sprint(":", server.config.Host.Port),
			server.config.Host.Ssl.Cert,
			server.config.Host.Ssl.Key))
}

func (server *Server) getRouter() *mux.Router {
	myRouter := mux.NewRouter().StrictSlash(true)

	// jwtMiddleware := server.getJwtMiddleware()

	// myRouter.HandleFunc("/", homePage)
	// addRouteWithMiddleware("/items", returnAllItems, myRouter, jwtMiddleware.HandlerWithNext)
	// addRouteWithMiddleware("/items/{id}", returnSingleItem, myRouter, jwtMiddleware.HandlerWithNext)

	return myRouter
}

func addRouteWithMiddleware(path string, function http.HandlerFunc, router *mux.Router, middleware negroni.HandlerFunc) {
	router.Handle(path, negroni.New(
		negroni.HandlerFunc(middleware),
		negroni.Wrap(http.HandlerFunc(function)),
	))
}

// func (server *Server) getJwtMiddleware() *jwtmiddleware.JWTMiddleware {
// 	return jwtmiddleware.New(
// 		jwtmiddleware.Options{
// 			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
// 				return []byte(server.config.Authentication.Jwt.Secret), nil
// 			},
// 			SigningMethod: jwt.SigningMethodHS256,
// 		})
// }

func homePage(context *gin.Context) {
	context.String(http.StatusOK, "Hello from Gin")
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
