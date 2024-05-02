package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestSqlApiSchemas(t *testing.T) {
	var config httpexpect.Config = httpexpect.Config{
		BaseURL:  baseUrl,
		Reporter: httpexpect.NewRequireReporter(t),
	}

	e := httpexpect.WithConfig(config)

	// Crea un squema en la BBDD
	t.Run("Create an empty schema", func(t *testing.T) {
		e.POST("/sql/schemas/{schemaName}").
			WithPath("schemaName", "example").
			Expect().
			Status(http.StatusCreated).
			JSON().Object().
			HasValue("message", "Schema created")
	})

	// Carga un archivo desde el sistema de archivos
	t.Run("Import a schema from sql file", func(t *testing.T) {
		file, err := os.Open("test_schema.sql")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		e.POST("/sql/schemas/{schemaName}").
			WithMultipart().
			WithFile("file", "test_schema.sql", file).
			WithPath("schemaName", "example2").
			Expect().
			Status(http.StatusCreated).
			JSON().Object().
			HasValue("message", "Schema Imported")
	})

	t.Run("Get database schemas", func(t *testing.T) {
		e.GET("/sql/schemas/").
			WithQuery("activesInPostgRest", "false").
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			HasValue("message", "public,example,example2")
	})

	t.Run("Gets the active squemas for PostgREST", func(t *testing.T) {
		e.GET("/sql/schemas/").
			WithQuery("activesInPostgRest", "true").
			Expect().
			Status(http.StatusOK).
			JSON().Object().
			HasValue("message", "example,example2")
	})
}

func TestSqlApiTables(t *testing.T) {
	var config httpexpect.Config = httpexpect.Config{
		BaseURL:  baseUrl,
		Reporter: httpexpect.NewRequireReporter(t),
	}

	e := httpexpect.WithConfig(config)

	t.Run("Create table", func(t *testing.T) {
		jsonBody := map[string]interface{}{
			"columns": []map[string]interface{}{
				{
					"name": "id",
					"type": "SERIAL PRIMARY KEY",
				},
				{
					"name": "name",
					"type": "varchar(255)",
				},
				{
					"name": "description",
					"type": "varchar(255)",
				},
				{
					"name": "components",
					"type": "integer",
				},
			},
		}

		e.POST("/sql/schemas/{schemaName}/tables/{tableName}").
			WithPath("schemaName", "example").
			WithPath("tableName", "enablerlist").
			WithJSON(jsonBody).
			Expect().
			Status(http.StatusCreated).
			JSON().Object().HasValue("message", "Table enablerlist created in Schema example")
	})
}

func TestSqlProxyToApiPostgREST(t *testing.T) {
	var config httpexpect.Config = httpexpect.Config{
		BaseURL:  baseUrl,
		Reporter: httpexpect.NewRequireReporter(t),
	}

	e := httpexpect.WithConfig(config)

	t.Run("Get schema swagger PostgREST", func(t *testing.T) {

		resp := e.GET("/sql/api/").
			WithHeader("Accept-Profile", "example").
			Expect().
			Status(http.StatusOK)

		resp.Header("Content-Profile").IsEqual("example")
		resp.JSON().Object()

	})

	t.Run("Insert data in table", func(t *testing.T) {
		e.POST("/sql/api/{proxyPath}").
			WithPath("proxyPath", "enablerlist").
			WithHeader("Content-Profile", "example").
			WithJSON(jsonBody).
			Expect().
			Status(http.StatusCreated)
	})

	t.Run("Get data in table", func(t *testing.T) {
		resp := e.GET("/sql/api/{proxyPath}").
			WithPath("proxyPath", "enablerlist").
			WithHeader("Content-Profile", "example").
			Expect().
			Status(http.StatusOK).
			JSON().Array()

		resp.Length().IsEqual(1)
		enabler := resp.Value(0)
		enabler.Object().Value("id").IsEqual(1)
		enabler.Object().Value("name").IsEqual("vpn 1")
		enabler.Object().Value("description").IsEqual("VPN 1 enabler based on WireGuard")
		enabler.Object().Value("components").IsEqual(1)
	})

	t.Run("Delete data in table", func(t *testing.T) {
		e.DELETE("/sql/api/{proxyPath}").
			WithPath("proxyPath", "enablerlist").
			WithQuery("id", "eq.1").
			WithHeader("Content-Profile", "example").
			Expect().
			Status(http.StatusNoContent)
	})
}
