package main

import (
	"fmt"
	"log"

	"go-api/middleware/auth"
	v1 "go-api/routers/v1"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()

	router.NoRoute(noRouteHandler)

	router.Static("/react", "../client/react-client/build")
	router.Static("/static", "../client/react-client/build/static")

	auth.UseJWTMiddleware(router, "/v1", v1.AuthRoutes)
}

func main() {
	config := NewConfig()

	PrintObject(config)

	log.Fatal(
		router.RunTLS(
			fmt.Sprint(":", config.Host.Port),
			config.Host.Ssl.Cert,
			config.Host.Ssl.Key))
}

func noRouteHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Printf("NoRoute claims: %#v\n", claims)
	c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
}
