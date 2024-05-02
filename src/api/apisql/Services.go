package apisql

import (
	"errors"
	"fmt"
	"io/ioutil"
	"ltse_api/common"
	"mime/multipart"
	"strings"
)

func getActiveSchemasService() string {
	db, err := GetDB()
	if err != nil {
		return err.Error()
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT rolconfig FROM pg_roles WHERE rolname = '%s'", getPgUser())

	row := db.QueryRow(query)

	var rolConfig string
	err = row.Scan(&rolConfig)

	if rolConfig == "" {
		return ""
	}

	if err != nil {
		panic(err)
	}

	rolConfig = strings.Trim(rolConfig, "{}")
	rolConfig = strings.ReplaceAll(rolConfig, "\"", "")
	schemas := strings.Split(rolConfig, "=")[1]

	return schemas
}

func getSchemasService() string {
	db, err := GetDB()
	if err != nil {
		return err.Error()
	}
	defer db.Close()

	query := "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT LIKE 'pg_%' AND schema_name != 'information_schema'"
	rows, _ := db.Query(query)

	schemas := ""

	for rows.Next() {
		var schema string
		err := rows.Scan(&schema)

		if err != nil {
			panic(err)
		}
		schemas += schema + ","
	}
	return schemas[:len(schemas)-1]
}

func createSchemaService(schemaName string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s AUTHORIZATION %s", schemaName, getPgUser()))
	if err != nil {
		return err
	}

	// Execute DDL statement
	_, err = db.Exec(fmt.Sprintf("SET search_path = %s", schemaName))
	if err != nil {
		return err
	}

	schemas := getActiveSchemasService()
	if schemas == "" {
		schemas += schemaName
	} else {
		schemas += "," + schemaName
	}

	return updatePostRestSchemas(schemas)
}

func schemaInPostRestService(schemaName string, active bool) error {
	err := cleanDeletedSchemas()
	if err != nil {
		return err
	}

	currentSchemasInDataBase := getSchemasService()

	if !strings.Contains(currentSchemasInDataBase, schemaName) {
		message := fmt.Sprintf("Schema '%s' not found in the database", schemaName)
		return errors.New(message)
	}

	activesSchemasInPostgRest := getActiveSchemasService()

	if active {
		if activesSchemasInPostgRest == "" {
			activesSchemasInPostgRest += schemaName
		} else {
			activesSchemasInPostgRest += "," + schemaName
		}
	} else {
		// Replace the string we want to delete with an empty string
		activesSchemasInPostgRest = strings.Replace(activesSchemasInPostgRest, schemaName+",", "", -1)
		activesSchemasInPostgRest = strings.Replace(activesSchemasInPostgRest, ","+schemaName, "", -1)
	}

	return updatePostRestSchemas(activesSchemasInPostgRest)
}

func importShemaFromFileService(schemaName string, file *multipart.FileHeader) error {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the contents of the file into a byte slice
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}

	// Convert the byte slice to a string
	bodyToString := string(data)

	// Call the existing service function to import the schema
	return importShema(schemaName, bodyToString)
}

func createTableService(ddl DDLCreateTable, tableName string, schemaName string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s AUTHORIZATION %s", schemaName, getPgUser()))
	if err != nil {
		return err
	}

	// Execute DDL statement
	_, err = db.Exec(fmt.Sprintf("SET search_path = %s", schemaName))
	if err != nil {
		return err
	}

	// check if table already exists
	_, table_check := db.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
	if table_check == nil {
		message := fmt.Sprintf("The table %s already exists in schema %s", tableName, schemaName)
		err := errors.New(message)
		return err
	}

	// create query
	var sql = "CREATE TABLE " + tableName + " ("

	for index, element := range ddl.Columns {
		sql += element.Name + " " + element.Type
		if index < len(ddl.Columns)-1 || len(ddl.Constraints) > 0 {
			sql += ", "
		}
	}
	for index, element := range ddl.Constraints {
		sql += fmt.Sprintf("CONSTRAINT %s %s", element.Name, element.Type)
		if index < len(ddl.Constraints)-1 {
			sql += ", "
		}
	}
	sql += ");"

	// execute query
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func truncateTableService(ddl DDLTableOptions, tableName string, schemaName string) error {
	if ddl.constraint != "CASCADE" && ddl.constraint != "RESTRICT" && ddl.constraint != "" {
		return errors.New("Constraint not valid. Allowed contraints are: CASCADE or RESTRICT")
	}

	if ddl.identity != "CONTINUE" && ddl.identity != "RESTART" && ddl.identity != "" {
		return errors.New("Identity policy not valid. Allowed identity policies are: CONTINUE or RESTART")
	}

	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute DDL statement
	_, err = db.Exec(fmt.Sprintf("SET search_path = %s", schemaName))
	if err != nil {
		return err
	}

	if ddl.identity != "" {
		ddl.identity += " IDENTITY"
	}

	_, err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %s %s %s;", tableName, ddl.identity, ddl.constraint))
	if err != nil {
		return err
	}
	return nil
}

func deleteTableService(ddl DDLTableOptions, tableName string, schemaName string) error {
	if ddl.constraint != "CASCADE" && ddl.constraint != "RESTRICT" && ddl.constraint != "" {
		return errors.New("Constraint not valid. Allowed contraints are: CASCADE or RESTRICT")
	}

	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute DDL statement
	_, err = db.Exec(fmt.Sprintf("SET search_path = %s", schemaName))
	if err != nil {
		return err
	}

	// _, err := db.Exec(`DROP TABLE $1 $2;`, tableName, constraint)
	_, err = db.Exec(fmt.Sprintf("DROP TABLE %s %s;", tableName, ddl.constraint))

	if err != nil {
		return err
	}
	return nil

}

func importShema(schemaName string, body string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !strings.Contains(body, "CREATE SCHEMA IF NOT EXISTS "+schemaName) &&
		!strings.Contains(body, "SET search_path = "+schemaName) {
		message := fmt.Sprintf("The name '%s' does not match with the schema name of the uploaded file.", schemaName)
		return errors.New(message)
	}

	// Execute DDL statement
	_, err = db.Exec(body)

	if err != nil {
		return err
	}

	schemas := getActiveSchemasService()
	if schemas == "" {
		schemas += schemaName
	} else {
		schemas += "," + schemaName
	}

	sql := fmt.Sprintf("ALTER ROLE %s SET pgrst.db_schemas = '%s'", getPgUser(), schemas)
	// Execute DDL statement
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	// Reload schema cache for postgREST
	sql = "NOTIFY pgrst, 'reload config'"
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func cleanDeletedSchemas() error {
	currentSchemas := getSchemasService()
	currentSchemasActivesSlice := strings.Split(getActiveSchemasService(), ",")

	currentSchemasActivesSlice = common.RemoveDuplicates(currentSchemasActivesSlice)

	var chemasActivesSlice []string

	for _, schema := range currentSchemasActivesSlice {
		if strings.Contains(currentSchemas, schema) {
			chemasActivesSlice = append(chemasActivesSlice, schema)
		}
	}

	schemasActives := strings.Join(chemasActivesSlice, ",")

	return updatePostRestSchemas(schemasActives)
}

func updatePostRestSchemas(schemas string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute DDL statement
	_, err = db.Exec(fmt.Sprintf("ALTER ROLE %s SET pgrst.db_schemas = '%s'", getPgUser(), schemas))
	if err != nil {
		return err
	}

	// Reload schema cache for postgREST
	_, err = db.Exec("NOTIFY pgrst, 'reload config'")
	if err != nil {
		return err
	}
	return nil
}
