package apinosql

import (
	"github.com/gin-gonic/gin"
)

func Routing(router *gin.RouterGroup) {
	router.Any("/api/*proxyPath", elasticProxy)
	router.PUT("/index/:indexName", createIndex)
	router.GET("/index/:indexName", getIndexInfo)
	router.PUT("/index/:indexName/document", addDocumentToIndex)
	router.PUT("/index/:indexName/document/:id", addDocumentWithIdToIndex)
	router.GET("/index/:indexName/document/:id", getDocumentByIdFromIndex)
}
