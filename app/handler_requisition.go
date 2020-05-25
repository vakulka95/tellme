package app

import (
	"log"
	"net/http"
	"strings"

	"gitlab.com/tellmecomua/tellme.api/app/representation"

	retry "github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
)

func (s *apiserver) apiRequisitionCreate(c *gin.Context) {
	var (
		err            error
		requisitionReq = &representation.CreateRequisitionRequest{}
	)

	if err = c.ShouldBindJSON(requisitionReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, representation.ErrInvalidRequest)
		return
	}

	requisition, err := representation.RequisitionAPItoPersistence(requisitionReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}

	requisitionRes, err := s.repository.CreateRequisition(requisition)
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			c.AbortWithStatusJSON(http.StatusConflict, representation.ErrConflict)
			return
		}
		log.Printf("(ERR) Failed to create requisition: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	go func() {
		err := retry.Do(
			func() error {
				return s.smscli.SendRequisitionApply(requisitionReq.Phone)
			},
		)
		if err != nil {
			log.Printf("(WARN) Failed to send requisition apply sms: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, representation.CreateRequisitionResponse{ID: requisitionRes.ID})
}
