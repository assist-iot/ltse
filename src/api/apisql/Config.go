package apisql

import (
	"fmt"
	"log"
	"ltse_api/config"
)

var pgRestConfig *config.PostgRestConfig = &config.CFG.PostgRestConfig
var pgDbConfig *config.PostgresDbConfig = &config.CFG.PostgresDbConfig

func getPgUser() string {
	return pgDbConfig.User
}

func getPgRestUrl() string {
	return "http://" + pgRestConfig.Host + ":" + pgRestConfig.Port
}

func getDbURL() string {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s database=%s sslmode=disable",
		pgDbConfig.Host,
		pgDbConfig.Port,
		pgDbConfig.User,
		pgDbConfig.Password,
		pgDbConfig.DBName,
	)
	log.Println(url)
	return url
}
