package config

import "ltse_api/common"

var CFG ApiConfig

type ApiConfig struct {
	Port             string
	PostgresDbConfig PostgresDbConfig
	ElasticConfig    ElasticConfig
	PostgRestConfig  PostgRestConfig
}

type PostgresDbConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type ElasticConfig struct {
	Host string
	Port string
}

type PostgRestConfig struct {
	Host string
	Port string
}

func LoadConfig() ApiConfig {
	CFG.Port = common.GetEnvironment("API_PORT", "8080")
	CFG.PostgresDbConfig = *loadConfigPostgres()
	CFG.ElasticConfig = *loadConfigElastic()
	CFG.PostgRestConfig = *loadConfigPostgREST()

	return CFG
}

func loadConfigPostgres() *PostgresDbConfig {
	return &PostgresDbConfig{
		Host:     common.GetEnvironment("POSTGRESQL_HOST", "localhost"),
		Port:     common.GetEnvironment("POSTGRESQL_PORT", "5432"),
		User:     common.GetEnvironment("POSTGRESQL_USER", "postgres"),
		Password: common.GetEnvironment("POSTGRESQL_PASS", "4ss1st10t"),
		DBName:   common.GetEnvironment("POSTGRESQL_DB", "assist"),
	}
}

func loadConfigElastic() *ElasticConfig {
	return &ElasticConfig{
		Host: common.GetEnvironment("ELASTIC_HOST", "localhost"),
		Port: common.GetEnvironment("ELASTIC_PORT", "9200"),
	}
}

func loadConfigPostgREST() *PostgRestConfig {
	return &PostgRestConfig{
		Host: common.GetEnvironment("POSTGREST_HOST", "localhost"),
		Port: common.GetEnvironment("POSTGREST_PORT", "3000"),
	}
}
