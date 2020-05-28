package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

func (s *apiserver) apiReviewConfirm(c *gin.Context) {
	const logPref = "apiReviewConfirm"

	var (
		err       error
		reviewReq = &representation.ConfirmReviewRequest{}
	)

	if err = c.ShouldBindJSON(reviewReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, representation.ErrInvalidRequest)
		return
	}

	review, err := representation.ReviewAPItoPersistence(reviewReq)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}

	reviewRes, err := s.repository.GetReviewByToken(review.Token)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get review: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	if reviewRes.Status == model.ReviewStatusCompleted {
		log.Printf("(ERR) %s: review [%s] already completed", logPref, reviewRes.ID)
		c.AbortWithStatusJSON(http.StatusConflict, representation.ErrConflict)
		return
	}

	review.ID = reviewRes.ID
	reviewRes, err = s.repository.UpdateReviewBodyStatus(review)
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			c.AbortWithStatusJSON(http.StatusConflict, representation.ErrConflict)
			return
		}
		log.Printf("(ERR) %s: failed to create requisition: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	c.JSON(http.StatusOK, representation.ConfirmReviewResponse{ID: reviewRes.ID})
}
