package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *apiserver) serviceDatabaseStat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": s.repository.Name(),
		"stat": s.repository.Stat(),
	})
}
