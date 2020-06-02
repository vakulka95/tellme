package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *apiserver) serviceDatabaseStat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": s.repository.Name(),
		"stat": s.repository.Stat(),
	})
}

func (s *apiserver) serviceSendRequisitionReview(c *gin.Context) {
	list, err := s.repository.GetNotReviewedRequisition()
	if err != nil {
		log.Printf("(ERR) Failed to fetch requisition list: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, item := range list {
		log.Printf("Phone: %s, ID: %s, ExpertID: %s", item.Phone, item.ID, item.ExpertID)
		err := s.requestRequisitionReview(item.Phone, item.ID, item.ExpertID)
		if err != nil {
			log.Printf("(WARN) Failed to send request for requisition review sms: %v", err)
		}
	}

	c.Status(http.StatusAccepted)
}
