package routes

import (
	"ltse_api/api/apinosql"
	"ltse_api/api/apisql"
	"ltse_api/swagger"

	// "net/http"

	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {

	// fmt.Println("join CORS")
	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	// // Second, we handle the OPTIONS problem
	// if c.Request.Method != "OPTIONS" {
	// 	c.Next()
	// } else {

	// 	// Everytime we receive an OPTIONS request,
	// 	// we just return an HTTP 200 Status Code
	// 	// Like this, Angular can now do the real
	// 	// request using any other method than OPTIONS
	// 	c.AbortWithStatus(http.StatusOK)
	// }
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORS)

	apisql.Routing(r.Group("/sql"))
	apinosql.Routing(r.Group("/nosql"))
	swagger.Routing(r.Group("/"))

	return r
}
