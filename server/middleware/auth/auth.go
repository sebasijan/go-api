package auth

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const IdentityKey = "id"

// // Use adds middleware to the group, see example code in GitHub.
// func (group *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
// 	group.Handlers = append(group.Handlers, middleware...)
// 	return group.returnObj()
// }

func UseJWTMiddleware(router *gin.Engine, path string, authRoutes func(r *gin.RouterGroup)) {
	authMiddleware := GetJwtMiddleware()

	api := router.Group("/api")
	api.POST("/login", authMiddleware.LoginHandler)

	authRouter := api.Group(path)
	authRouter.Use(authMiddleware.MiddlewareFunc())
	{
		authRouter.GET("/refresh_token", authMiddleware.RefreshHandler)
		authRoutes(authRouter)
	}
}

type GinJWTAuthenticationOptions struct {
	Authenticator func(c *gin.Context) (interface{}, error)
}

func getDefaultGinJWTOptions() GinJWTAuthenticationOptions {
	return GinJWTAuthenticationOptions{
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Bo-Yi",
					FirstName: "Wu",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
	}
}

func GetJwtMiddleware() *jwt.GinJWTMiddleware {
	options := getDefaultGinJWTOptions()

	authMiddleware, error := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					IdentityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims[IdentityKey].(string),
			}
		},
		Authenticator: options.Authenticator,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if error != nil {
		log.Fatal("JTW Error: ", error.Error())
	}

	return authMiddleware
}
