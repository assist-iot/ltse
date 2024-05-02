package apisql

import (
	"github.com/gin-gonic/gin"
)

func Routing(router *gin.RouterGroup) {

	router.Any("/api/*proxyPath", postgrestProxy)
	router.GET("/schemas", getSchemas)
	router.POST("/schemas/:schemaName", createSchema)
	router.PUT("/schemas/:schemaName", schemaInPostRest)
	router.POST("/schemas/:schemaName/tables/:tableName", createTable)
	router.PUT("/schemas/:schemaName/tables/:tableName", truncateTable)
	router.DELETE("/schemas/:schemaName/tables/:tableName", deleteTable)
}
