package apisql

import (
	"fmt"
	"log"
	"ltse_api/common"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Postgress SQL
// @summary For use API PostgREST, initial path
// @description You must take into account that for some operations it is necessary to use specific headers such as **Accept-Profile** or **Content-Profile** both to indicate the name of the schema when making the request. You can get more information in the api guide.
// @description
// @description **RestApi PostgRest 9.0.0 guide**
// @description https://postgrest.org/en/v9.0/
// @Tags sql PostgREST
// @Get PostgRest Doc of Schema
// @Produce json
// @Param Accept-Profile header string true "Active Schema in PostgRest"
// @Param proxyPath path string false "Optional PostgRest Options" format:"omitempty"
// @Header 200 {string} Accept-Profile "Header description"
// @Success 200
// @Failure 409
// @Router /sql/api/ [get]
func postgrestProxy(g *gin.Context) {
	proxyPath := g.Param("proxyPath")

	err := common.ReverseProxy(getPgRestUrl(), g, proxyPath)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
}

// @summary Get database schema list
// @description You can get the list activated schemas for use in API PostgREST by using the option activesInPostgRest=true
// @Tags SQL API
// @Produce json
// @Param activesInPostgRest query boolean false "If true, show only the activate schemas for to be used in PostgREST" default("false")
// @Success 200 "Schema list message"
// @Failure 500 "Error message"
// @Router /sql/schemas [get]
func getSchemas(g *gin.Context) {
	active, err := strconv.ParseBool(g.DefaultQuery("activesInPostgRest", "false"))
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = cleanDeletedSchemas()
	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if active {
		g.JSON(http.StatusOK, gin.H{
			"message": getActiveSchemasService(),
		})
	} else {
		g.JSON(http.StatusOK, gin.H{
			"message": getSchemasService(),
		})
	}
}

// @summary Create schema in database
// @description You can create an empty schema or import it from file.
// @Tags SQL API
// @Param schemaName path string true "Schema name"
// @Param file formData file false "SQL file to import. The schema name in the script must be the same as the declared schemaName field above."
// @Success 200 "Schema created message"
// @Failure 500 "Error message"
// @Router /sql/schemas/{schemaName} [post]
func createSchema(g *gin.Context) {
	var err error
	schemaName := g.Param("schemaName")
	message := "Schema created"

	// Create schema without file
	if g.Request.Body == nil || g.Request.ContentLength == 0 {
		log.Println(message)
		err = createSchemaService(schemaName)
	} else {
		// Check if file is being uploaded
		message = "Schema Imported"
		var file multipart.File
		var fileHeader *multipart.FileHeader
		file, fileHeader, err = g.Request.FormFile("file")
		if err == nil {
			// Import schema from file
			log.Println("Import Schema from file")
			defer file.Close()
			err = importShemaFromFileService(schemaName, fileHeader)
		}
	}

	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	} else {
		g.JSON(http.StatusCreated, gin.H{
			"message": message,
		})
	}
}

// @summary Enable or disable the schema for use in API PostgREST
// @description You can Enable or disable the schema for use in API PostgREST
// @Tags SQL API
// @Produce json
// @Param schemaName path string true "Schema name"
// @Param isActiveInPostgRest query boolean true "Activate or deactivate the scheme in postrest" default("true")
// @Success 200 "Schema updated message"
// @Failure 500 "Error message"
// @Router /sql/schemas/{schemaName} [put]
func schemaInPostRest(g *gin.Context) {
	var err error
	schemaName := g.Param("schemaName")
	active, err := strconv.ParseBool(g.DefaultQuery("isActiveInPostgRest", "true"))
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = schemaInPostRestService(schemaName, active)

	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	} else {
		g.JSON(http.StatusOK, gin.H{
			"message": "Schema updated",
		})
	}
}

// @summary Create table
// @description You can create a table on schema
// @Tags SQL API
// @Produce json
// @Param schemaName path string true "Schema name"
// @Param tableName path string true "Schema name"
// @Param raw_body body DDLCreateTable false "JSON"
// @Success 200 "Table created messsage"
// @Failure 500 "Error message"
// @Router /sql/schemas/{schemaName}/tables/{tableName} [post]
func createTable(g *gin.Context) {
	schemaName := g.Param("schemaName")
	tableName := g.Param("tableName")
	var ddl DDLCreateTable
	err := g.BindJSON(&ddl)

	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	err = createTableService(ddl, tableName, schemaName)

	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	} else {
		message := fmt.Sprintf("Table %s created in Schema %s", tableName, schemaName)
		g.JSON(http.StatusCreated, gin.H{
			"message": message,
		})
	}
}

// @summary Truncate table
// @description You can truncate a table on schema.
// @Tags SQL API
// @Produce json
// @Param schemaName path string true "Schema name"
// @Param tableName path string true "Schema name"
// @Param identity query string false "RESTART or IDENTITY"
// @Param constraint query string false "CASCADE or RESTRICT"
// @Success 200 "Table truncate message"
// @Failure 500 "Error message"
// @Router /sql/schemas/{schemaName}/tables/{tableName} [put]
func truncateTable(g *gin.Context) {
	schemaName := g.Param("schemaName")
	tableName := g.Param("tableName")
	var ddl DDLTableOptions

	ddl.identity = g.DefaultQuery("identity", "")
	ddl.constraint = g.DefaultQuery("constraint", "")

	err := truncateTableService(ddl, tableName, schemaName)

	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	} else {
		message := fmt.Sprintf("Table %s truncated in Schema %s", tableName, schemaName)
		g.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	}
}

// @summary Delete table
// @description You can delete table on schema. You must first truncate the table.
// @Tags SQL API
// @Produce json
// @Param schemaName path string true "Schema name"
// @Param tableName path string true "Schema name"
// @Param constraint query string false "CASCADE or RESTRICT"
// @Success 200 "Table dropped message"
// @Failure 500 "Error message"
// @Router /sql/schemas/{schemaName}/tables/{tableName} [delete]
func deleteTable(g *gin.Context) {
	tableName := g.Param("tableName")
	schemaName := g.Param("schemaName")
	var ddl DDLTableOptions
	ddl.constraint = g.DefaultQuery("constraint", "")

	err := deleteTableService(ddl, tableName, schemaName)

	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
	} else {
		message := fmt.Sprintf("Table %s dropped in SchemaName %s", tableName, schemaName)
		g.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	}
}
