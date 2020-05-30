package app

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

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
				"search":          qlp.Search,
				"status":          qlp.Status,
				"specializations": representation.GenerateDiagnosesOptions(qlp.Specializations),
			},
		},
	)
}

func (s *apiserver) webAdminExpertItem(c *gin.Context) {
	const logPref = "webAdminExpertItem"

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/admin/login")
		return
	}

	expertID := c.Param("expertId")
	expertRes, err := s.repository.GetExpertView(expertID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: expert [%s] not found", logPref, expertID)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get expert view [%s]: %v", logPref, expertID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	c.HTML(http.StatusOK, "expert_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
			},
			"expert": representation.ExpertViewItemPersistenceToAPI(expertRes),
		},
	)
}

func (s *apiserver) webAdminExpertBlock(c *gin.Context) {
	expertID := c.Param("expertId")

	_, err := s.repository.UpdateExpertStatus(&model.Expert{ID: expertID, Status: model.ExpertStatusBlocked})
	if err != nil {
		log.Printf("(ERR) Failed to update expert: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminExpertActivate(c *gin.Context) {
	expertID := c.Param("expertId")

	_, err := s.repository.UpdateExpertStatus(&model.Expert{ID: expertID, Status: model.ExpertStatusActive})
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
