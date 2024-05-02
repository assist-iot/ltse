package apinosql

import (
	"fmt"
	"io/ioutil"
	"log"
	"ltse_api/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ElasticProxy godoc
// @summary For use API ElasticSearch, initial path
// @Tags nosql API ElasticSearch
// @description RestApi ElasticSearch 7.17 guide
// @description https://www.elastic.co/guide/en/elasticsearch/reference/7.17/rest-apis.html
// @Produce json
// @Param proxyPath path string false "Optional PostgRest Options" format:"omitempty"
// @Success 200
// @Failure 409
// @Router /nosql/api/ [get]
func elasticProxy(g *gin.Context) {
	proxyPath := g.Param("proxyPath")

	err := common.ReverseProxy(getESUrl(), g, proxyPath)

	if err != nil {
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}
}

// Create Index godoc
// @summary Create index in elastic
// @Tags noSQL API
// @Create Index
// @Produce json
// @Param indexName path string true "Index name"
// @Success 200
// @Failure 409
// @Router /nosql/index/{indexName} [put]
func createIndex(g *gin.Context) {
	indexName := g.Param("indexName")

	err := createIndexService(indexName)
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Index %s created", indexName),
	})
}

// Create Index godoc
// @summary Get index info
// @description Returns information about indexname index from the ltse nosql cluster
// @Tags noSQL API
// @Produce json
// @Param indexName path string true "index name"
// @Success 200
// @Failure 500
// @Router /nosql/index/{indexName} [get]
func getIndexInfo(g *gin.Context) {
	indexName := g.Param("indexName")

	result, err := getIndexInfoService(indexName)
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, result)
}

// Create Index godoc
// @summary Add a json document to index and makes it searchable with an id
// @description Add a json document to the specified <indexname> index of the ltse nosql cluster and makes it searchable with an <_id>
// @Tags noSQL API
// @Produce json
// @Param indexName path string true "Index name"
// @Success 201
// @Failure 500
// @Router /nosql/index/{indexName}/document/{id} [put]
func addDocumentWithIdToIndex(g *gin.Context) {
	indexName := g.Param("indexName")
	id := g.Param("id")
	document, err := ioutil.ReadAll(g.Request.Body)

	if err != nil {
		if err != nil {
			log.Println(err)
			g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	err = addDocumentWithIdToIndexService(indexName, document, id)
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	g.Status(http.StatusCreated)
}

// Create Index godoc
// @summary Add a json document to index without id
// @description Adds a json document to the specified <indexname > index of the ltse nosql cluster and return id to makes it searchable
// @Tags noSQL API
// @Produce json
// @Param indexName path string true "Index name"
// @Success 201
// @Failure 500
// @Router /nosql/index/{indexName}/document [put]
func addDocumentToIndex(g *gin.Context) {
	indexName := g.Param("indexName")
	document, err := ioutil.ReadAll(g.Request.Body)

	if err != nil {
		if err != nil {
			log.Println(err)
			g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}
	var result map[string]interface{}
	result, err = addDocumentToIndexService(indexName, document)
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, result)
}

// Create Index godoc
// @summary Get json document by id
// @description Retrieves the specified json document <_id> from the indexname of the ltse nosql cluster.
// @Tags noSQL API
// @Produce json
// @Param indexName path string true "Index name"
// @Success 200
// @Failure 500
// @Router /nosql/index/{indexName}/document/{id} [get]
func getDocumentByIdFromIndex(g *gin.Context) {
	indexName := g.Param("indexName")
	id := g.Param("id")

	result, err := getDocumentByIdFromIndexService(indexName, id)
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, result)
}
