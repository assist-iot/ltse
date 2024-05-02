package apinosql

import (
	"errors"
	"fmt"
	"log"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

var clientES *elasticsearch7.Client

func ConnectToElasticsearch() error {
	var err error
	cfg := elasticsearch7.Config{
		Addresses: []string{getESUrl()},
	}
	clientES, err = elasticsearch7.NewClient(cfg)
	if err != nil || clientES == nil {
		message := fmt.Sprintf("Error creating the client: %s", err)
		return errors.New(message)
	}

	res, err := clientES.Ping()
	if err != nil {
		message := fmt.Sprintf("error creating Elasticsearch client: %s", err)
		return errors.New(message)
	}
	defer res.Body.Close()

	if res.IsError() {
		message := fmt.Sprintf("Elasticsearch returned an error: %s", res.String())
		return errors.New(message)
	}

	log.Println("Successfully created connection to elasticsearch")

	return nil
}

func getClient() (*elasticsearch7.Client, error) {
	err := ConnectToElasticsearch()
	if err != nil {
		return nil, err
	}

	return clientES, nil
}
