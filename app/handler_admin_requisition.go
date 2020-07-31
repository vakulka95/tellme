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
	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusFound, "admin/login")
		return
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "admin/login")
		return
	}

	if iRole.(string) == UserRoleAdmin || iStatus.(string) == model.ExpertStatusActive {
		c.Redirect(http.StatusFound, "/admin/requisition?limit=10&offset=0")
		return
	}
	c.Redirect(http.StatusFound, "/admin/profile")
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
		c.Redirect(http.StatusFound, "/admin/login")
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

	var (
		renderedList gin.H
		username     string
		role         = iRole.(string)
	)

	if !isExpert {
		renderedList = representation.RequisitionListPersistenceToAPIForAdmin(list)
	} else {
		iUserID, ok := c.Get("userID")
		if !ok {
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}
		userID := iUserID.(string)
		renderedList = representation.RequisitionListPersistenceToAPIForExpert(list, userID)

		expert, err := s.repository.GetExpert(userID)
		if err != nil {
			log.Printf("(ERR) Failed to fetch expeert details: %v", err)
		}
		username = expert.Username
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	datetimeFrom, datetimeTo := qlp.GetDatetimeRange()
	c.HTML(http.StatusOK, "requisition_list.html",
		gin.H{
			"metadata": gin.H{
				"title":     "Заявки",
				"username":  username,
				"logged_in": true,
				"role":      role,
				"status":    iStatus.(string),
			},
			"data":       renderedList,
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
		c.Redirect(http.StatusFound, "/admin/login")
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

	var (
		item gin.H
		role = iRole.(string)
	)

	if role == UserRoleAdmin {
		item = representation.RequisitionItemPersistenceToAPIForAdmin(requisitionRes)
	} else {
		iUserID, ok := c.Get("userID")
		if !ok {
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}
		item = representation.RequisitionItemPersistenceToAPIForExpert(requisitionRes, iUserID.(string))
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "requisition_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      role,
				"status":    iStatus.(string),
			},
			"requisition": item,
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
	const logPref = "webAdminRequisitionDiscard"

	var requisitionID = c.Param("requisitionId")

	err := s.repository.DeleteRequisitionSessions(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to delete requisition sessions: %v", logPref, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = s.repository.UpdateRequisitionStatus(&model.Requisition{ID: requisitionID, Status: model.RequisitionStatusCreated})
	if err != nil {
		log.Printf("(ERR) %s: failed to update requisition: %v", logPref, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminRequisitionNoAnswer(c *gin.Context) {
	const logPref = "webAdminRequisitionNoAnswer"

	var requisitionID = c.Param("requisitionId")

	err := s.repository.DeleteRequisitionSessions(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to delete requisition sessions: %v", logPref, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	requisitionRes, err := s.repository.GetRequisition(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to get requisition [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	_, err = s.repository.UpdateRequisitionStatus(&model.Requisition{
		ID: requisitionID, ExpertID: requisitionRes.ExpertID, Status: model.RequisitionStatusNoAnswer,
	})
	if err != nil {
		log.Printf("(ERR) %s: failed to update requisition: %v", logPref, err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminRequisitionComplete(c *gin.Context) {
	const logPref = "webAdminRequisitionComplete"

	var requisitionID = c.Param("requisitionId")

	requisitionSessions, err := s.repository.GetRequisitionSessionList(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to get requisition sessions [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	requisitionRes, err := s.repository.GetRequisition(requisitionID)
	if err != nil {
		log.Printf("(ERR) %s: failed to get expert view [%s]: %v", logPref, requisitionID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}
	expertID := requisitionRes.ExpertID

	if len(requisitionSessions) == 0 {
		_, err = s.repository.CreateSession(&model.Session{
			ID:            uuid.NewV4().String(),
			ExpertID:      expertID,
			RequisitionID: requisitionID,
		})
		if err != nil {
			log.Printf("(ERR) %s: failed to create session on requisition [%s]: %v", logPref, requisitionID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
			return
		}
	}

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

func (s *apiserver) webAdminRequisitionDelete(c *gin.Context) {
	requisitionID := c.Param("requisitionId")

	err := s.repository.DeleteRequisition(&model.Requisition{ID: requisitionID})
	if err != nil {
		log.Printf("(ERR) Failed to delete requisitions: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}
