package rest

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func createServer(ctx ...interface{}) {

	if server != nil {
		return
	}

	server = gin.Default()

	server.Use(gzip.Gzip(gzip.DefaultCompression))
	server.Use(GinLoggerMiddleware(ctx...))

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           50 * time.Second,
	}))

	server.Use(ErrorHandler)
	// server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func getServer() *gin.Engine {
	if server == nil {
		createServer()
	}

	return server
}

func GinCtx(c *gin.Context) []interface{} {
	var ctx []interface{}
	// mock_ctx solo es para mocks en testing
	if mocks, ok := c.Get("mock_ctx"); ok {
		ctx = mocks.([]interface{})
	}

	ctx = append(ctx, ginLogger(c))

	return ctx
}
