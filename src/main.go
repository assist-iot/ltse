package main

import (
	"ltse_api/config"
	_ "ltse_api/docs"
	"ltse_api/routes"
)

// @title LTSE
// @version 2.0.0
// @host localhost:8080
// @BasePath /
func main() {
	config.LoadConfig()
	r := routes.SetupRouter()
	r.Run(":" + config.CFG.Port)
}
