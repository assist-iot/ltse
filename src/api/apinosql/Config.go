package apinosql

import (
	"ltse_api/config"
)

var elastic *config.ElasticConfig = &config.CFG.ElasticConfig

func getESUrl() string {
	return "http://" + elastic.Host + ":" + elastic.Port
}
