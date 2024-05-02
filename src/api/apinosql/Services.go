package apinosql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func createIndexService(indexName string) error {
	es, err := getClient()
	if err != nil {
		log.Println(err)
		return err
	}

	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		log.Println(res)
		message := fmt.Sprintf("The index: %s exists", indexName)
		return errors.New(message)
	}
	defer res.Body.Close()

	res, err = es.Indices.Create(indexName)
	if err != nil {
		log.Println(err)
		message := fmt.Sprintf("%s", err)
		return errors.New(message)
	}
	defer res.Body.Close()

	return nil
}

func getIndexInfoService(indexName string) (map[string]interface{}, error) {
	es, err := getClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := es.Indices.Get([]string{indexName})
	if err != nil {
		log.Println(err)
		message := fmt.Sprintf("%s", err)
		return nil, errors.New(message)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func addDocumentWithIdToIndexService(indexName string, document []byte, id string) error {
	es, err := getClient()
	if err != nil {
		log.Println(err)
		return err
	}

	var res *esapi.Response
	res, err = es.Index(indexName, bytes.NewReader(document), es.Index.WithDocumentID(id))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		message := fmt.Sprintf("Elasticsearch returned an error: %s", res.String())
		return errors.New(message)
	}

	return nil
}

func addDocumentToIndexService(indexName string, document []byte) (map[string]interface{}, error) {
	es, err := getClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var res *esapi.Response
	res, err = es.Index(indexName, bytes.NewReader(document), es.Index.WithOpType("create"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if res.IsError() {
		message := fmt.Sprintf("Elasticsearch returned an error: %s", res.String())
		return nil, errors.New(message)
	}

	return result, nil
}

func getDocumentByIdFromIndexService(indexName string, id string) (map[string]interface{}, error) {
	es, err := getClient()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var res *esapi.Response
	res, err = es.Get(indexName, id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		message := fmt.Sprintf("Elasticsearch returned an error: %s", res.String())
		return nil, errors.New(message)
	}

	var result map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
