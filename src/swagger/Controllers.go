package swagger

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func apiExport(g *gin.Context) {
	swaggerPath, err := filepath.Abs("./docs/swagger.json")
	if err != nil {
		log.Println(err)
		g.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	g.File(swaggerPath)
}
