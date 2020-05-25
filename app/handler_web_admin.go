package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

func (s *apiserver) webIndexPage(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "/admin/requisition?limit=10&offset=0")
}

//
// Experts
//
func (s *apiserver) webAdminExpertList(c *gin.Context) {
	qlp := &representation.QueryListParams{}

	if err := c.BindQuery(qlp); err != nil {
		log.Printf("(ERR) Bind query error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	list, err := s.repository.GetExpertList(representation.QueryExpertAPItoPersistence(qlp))
	if err != nil {
		log.Printf("(ERR) Failed to fetch expert list: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "expert_list.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
			},
			"data":       representation.ExpertListPersistenceToAPI(list),
			"pagination": qlp.GeneratePagination(list.Total),
			"queries": gin.H{
				"status":          qlp.Status,
				"specializations": representation.GenerateDiagnosesOptions(qlp.Specializations),
			},
		},
	)
}

func (s *apiserver) webAdminExpertBlock(c *gin.Context) {
	expertID := c.Param("expertId")

	_, err := s.repository.UpdateExpertStatus(&model.Expert{ID: expertID, Status: model.StatusBlocked})
	if err != nil {
		log.Printf("(ERR) Failed to update expert: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminExpertActivate(c *gin.Context) {
	expertID := c.Param("expertId")

	_, err := s.repository.UpdateExpertStatus(&model.Expert{ID: expertID, Status: model.StatusActive})
	if err != nil {
		log.Printf("(ERR) Failed to update expert: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminUpdateExpertPassword(c *gin.Context) {
	var (
		expertID = c.Param("expertId")
		pwdReq   = &representation.UpdatePasswordRequest{}
	)

	if err := c.ShouldBindJSON(pwdReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, representation.ErrInvalidRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pwdReq.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = s.repository.UpdateExpertPassword(&model.Expert{ID: expertID, Password: string(hashedPass)})
	if err != nil {
		log.Printf("(ERR) Failed to update expert: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminExpertDelete(c *gin.Context) {
	expertID := c.Param("expertId")

	err := s.repository.DeleteExpert(&model.Expert{ID: expertID})
	if err != nil {
		log.Printf("(ERR) Failed to update expert: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

//
// Requisitions
//
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
		qlp.Status = model.StatusCreated
	}

	requisitionQuery := representation.QueryRequisitionAPItoPersistence(qlp)
	if qlp.Status != model.StatusCreated && isExpert {
		iExpertID, _ := c.Get("userID")
		requisitionQuery.ExpertID = iExpertID.(string)
	}

	list, err := s.repository.GetRequisitionList(requisitionQuery)
	if err != nil {
		log.Printf("(ERR) Failed to fetch requisition list: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "requisition_list.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
			},
			"data":       representation.RequisitionListPersistenceToAPI(list),
			"pagination": qlp.GeneratePagination(list.Total),
			"queries": gin.H{
				"status":            qlp.Status,
				"specializations":   representation.GenerateDiagnosesOptions(qlp.Specializations),
				"feedback_time":     qlp.FeedbackTime,
				"feedback_week_day": qlp.FeedbackWeekDay,
			},
		},
	)
}

func (s *apiserver) webAdminRequisitionTake(c *gin.Context) {
	expertID, _ := c.Get("userID")
	requisitionID := c.Param("requisitionId")

	_, err := s.repository.UpdateRequisitionStatus(&model.Requisition{ID: requisitionID, ExpertID: expertID.(string), Status: model.StatusProcessing})
	if err != nil {
		log.Printf("(ERR) Failed to update requisition: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminRequisitionComplete(c *gin.Context) {
	expertID, _ := c.Get("userID")
	requisitionID := c.Param("requisitionId")

	_, err := s.repository.UpdateRequisitionStatus(&model.Requisition{ID: requisitionID, ExpertID: expertID.(string), Status: model.StatusCompleted})
	if err != nil {
		log.Printf("(ERR) Failed to update requisition: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}
