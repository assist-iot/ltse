package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestNoSqlApi(t *testing.T) {
	var config httpexpect.Config = httpexpect.Config{
		BaseURL:  baseUrl,
		Reporter: httpexpect.NewRequireReporter(t),
	}

	e := httpexpect.WithConfig(config)

	// Crea un squema en la BBDD
	t.Run("Create NoSQL index", func(t *testing.T) {
		e.PUT("/nosql/index/{indexName}").
			WithPath("indexName", "example").
			Expect().
			Status(http.StatusCreated).
			JSON().Object().
			HasValue("message", "Index example created")
	})

	t.Run("Get Index Info", func(t *testing.T) {
		resp := e.GET("/nosql/index/{indexName}").
			WithPath("indexName", "example").
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.ContainsKey("example")
	})

	t.Run("Add Document without Id", func(t *testing.T) {
		resp := e.PUT("/nosql/index/{indexName}/document").
			WithPath("indexName", "example").
			WithJSON(jsonBody).
			Expect().
			Status(http.StatusCreated).
			JSON().Object()

		resp.ContainsKey("_id")
		resp.Value("_index").IsEqual("example")
		resp.Value("result").IsEqual("created")
	})

	t.Run("Add Document to index with ID", func(t *testing.T) {
		documentId := "123456789"

		e.PUT("/nosql/index/{indexName}/document/{id}").
			WithPath("indexName", "example").
			WithPath("id", documentId).
			WithJSON(jsonBody2).
			Expect().
			Status(http.StatusCreated)
	})

	t.Run("Get Document by Id from an index", func(t *testing.T) {
		documentId := "123456789"

		resp := e.GET("/nosql/index/{indexName}/document/{id}").
			WithPath("indexName", "example").
			WithPath("id", documentId).
			WithJSON(jsonBody2).
			Expect().
			Status(http.StatusOK).
			JSON().Object()

		resp.ContainsKey("_id")
		resp.Value("_id").IsEqual(documentId)
		resp.Value("_index").IsEqual("example")
		resp.Value("found").IsEqual(true)
	})
}
