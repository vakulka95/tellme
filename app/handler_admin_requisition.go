package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"
	"gitlab.com/tellmecomua/tellme.api/pkg/util/generator"
)

func (s *apiserver) webIndexPage(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/admin/requisition?limit=10&offset=0")
}

func (s *apiserver) webAdminRequisitionList(c *gin.Context) {
	qlp := &representation.QueryListParams{}

	if err := c.BindQuery(qlp); err != nil {
		log.Printf("(ERR) Bind query error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		return
	}

	isExpert := iRole.(string) == UserRoleExpert
	if qlp.Status == "" && isExpert {
		qlp.Status = model.RequisitionStatusCreated
	}

	requisitionQuery := representation.QueryRequisitionAPItoPersistence(qlp)
	if qlp.Status != model.RequisitionStatusCreated && isExpert {
		iExpertID, _ := c.Get("userID")
		requisitionQuery.ExpertID = iExpertID.(string)
	}

	list, err := s.repository.GetRequisitionList(requisitionQuery)
	if err != nil {
		log.Printf("(ERR) Failed to fetch requisition list: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	datetimeFrom, datetimeTo := qlp.GetDatetimeRange()
	c.HTML(http.StatusOK, "requisition_list.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
			},
			"data":       representation.RequisitionListPersistenceToAPI(list),
			"pagination": qlp.GeneratePagination(list.Total),
			"queries": gin.H{
				"search":            qlp.Search,
				"status":            qlp.Status,
				"specializations":   representation.GenerateDiagnosesOptions(qlp.Specializations),
				"feedback_time":     qlp.FeedbackTime,
				"feedback_week_day": qlp.FeedbackWeekDay,
				"datetime_from":     datetimeFrom,
				"datetime_to":       datetimeTo,
			},
		},
	)
}

func (s *apiserver) webAdminRequisitionItem(c *gin.Context) {
	const logPref = "webAdminRequisitionItem"

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		return
	}

	requisitionID := c.Param("requisitionId")
	requisitionRes, err := s.repository.GetRequisition(requisitionID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: requisition [%s] not found", logPref, requisitionID)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get expert view [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	c.HTML(http.StatusOK, "requisition_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
			},
			"requisition": representation.RequisitionItemPersistenceToAPI(requisitionRes),
		},
	)
}

func (s *apiserver) webAdminRequisitionTake(c *gin.Context) {
	const logPref = "webAdminRequisitionTake"

	var (
		role, _       = c.Get("role")
		iExpertID, _  = c.Get("userID")
		requisitionID = c.Param("requisitionId")
	)

	expertID := iExpertID.(string)

	if role == UserRoleAdmin {
		requisitionRes, err := s.repository.GetRequisition(requisitionID)
		if err != nil {
			log.Printf("(ERR) %s: failed to get expert view [%s]: %v", logPref, requisitionID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
			return
		}
		expertID = requisitionRes.ExpertID
	}

	_, err := s.repository.UpdateRequisitionStatus(&model.Requisition{ID: requisitionID, ExpertID: expertID, Status: model.RequisitionStatusProcessing})
	if err != nil {
		log.Printf("(ERR) Failed to update requisition: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminRequisitionDiscard(c *gin.Context) {
	requisitionID := c.Param("requisitionId")

	_, err := s.repository.UpdateRequisitionStatus(&model.Requisition{ID: requisitionID, Status: model.RequisitionStatusCreated})
	if err != nil {
		log.Printf("(ERR) Failed to update requisition: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminRequisitionComplete(c *gin.Context) {
	const logPref = "webAdminRequisitionComplete"

	var requisitionID = c.Param("requisitionId")

	requisitionRes, err := s.repository.GetRequisition(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to get expert view [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}
	expertID := requisitionRes.ExpertID

	_, err = s.repository.UpdateRequisitionStatus(&model.Requisition{ID: requisitionID, ExpertID: expertID, Status: model.RequisitionStatusCompleted})
	if err != nil {
		log.Printf("(ERR) Failed to update requisition: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	go func() {
		err := retry.Do(
			func() error {
				return s.requestRequisitionReview(requisitionRes.Phone, requisitionID, expertID)
			},
		)
		if err != nil {
			log.Printf("(WARN) Failed to send request for requisition review sms: %v", err)
		}
	}()

	c.Status(http.StatusAccepted)
}

func (s *apiserver) requestRequisitionReview(requisitionPhone, requisitionID, expertID string) error {
	token := generator.NewReviewToken()
	_, err := s.repository.CreateReview(
		&model.Review{
			ID:            uuid.NewV4().String(),
			ExpertID:      expertID,
			RequisitionID: requisitionID,
			Token:         token,
			Status:        model.ReviewStatusRequested,
		},
	)
	if err != nil {
		return err
	}

	return s.smscli.SendRequisitionReview(
		requisitionPhone,
		fmt.Sprintf("https://%s/review?token=%s", s.config.DomainName, token),
	)
}
