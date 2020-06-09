package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

func (s *apiserver) webAdminSessionCreate(c *gin.Context) {
	const logPref = "webAdminSessionCreate"

	var (
		iExpertID, _  = c.Get("userID")
		requisitionID = c.Param("requisitionId")
	)

	requisition, err := s.repository.GetRequisition(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to get requisition [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	if requisition.Status == model.RequisitionStatusCompleted ||
		requisition.SessionCount >= maxRequisitionSessionCount {
		log.Printf("(WARN) %s: requisition [%s] already processed: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusConflict, representation.ErrConflict)
		return
	}

	_, err = s.repository.CreateSession(&model.Session{
		ID:            uuid.NewV4().String(),
		ExpertID:      iExpertID.(string),
		RequisitionID: requisitionID,
	})
	if err != nil {
		log.Printf("(ERR) %s: failed to create session on requisition [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	c.Status(http.StatusAccepted)
}
