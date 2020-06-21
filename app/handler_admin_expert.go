package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	uuid "github.com/satori/go.uuid"
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
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "expert_list.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
				"status":    iStatus.(string),
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
		c.Redirect(http.StatusFound, "/admin/login")
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

	expertComments, err := s.repository.GetExpertCommentList(expertID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: expert [%s] not found", logPref, expertID)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get expert comments view [%s]: %v", logPref, expertID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "expert_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
				"status":    iStatus.(string),
			},
			"expert":   representation.ExpertViewItemPersistenceToAPI(expertRes),
			"comments": representation.CommentViewListPersistenceToAPI(expertComments),
		},
	)
}

func (s *apiserver) webAdminUpdateExpertItem(c *gin.Context) {
	const logPref = "webAdminUpdateExpertItem"

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	var (
		expertID   = c.Param("expertId")
		expertRole = iRole.(string)
		req        = &representation.UpdateExpertRequest{}
	)

	err := c.MustBindWith(req, binding.FormPost)
	if err != nil {
		log.Printf("(ERR) %s: failed to bind request form: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}
	req.ID = expertID

	expert, err := representation.ExpertAdminFormToPersistence(req, expertRole == UserRoleAdmin)
	if err != nil {
		log.Printf("(ERR) %s: invalid request: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}

	if expertRole == UserRoleExpert {
		_, err = s.repository.UpdateExpert(expert)
	} else {
		_, err = s.repository.UpdateExpertSpecializations(expert)
	}
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

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "expert_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      expertRole,
				"status":    iStatus,
			},
			"expert": representation.ExpertViewItemPersistenceToAPI(expertRes),
		},
	)
}

func (s *apiserver) webAdminExpertCommentCreate(c *gin.Context) {
	const logPref = "webAdminExpertCommentCreate"

	iUserID, ok := c.Get("userID")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	var (
		expertID = c.Param("expertId")
		req      = &representation.CreateExpertCommentRequest{}
	)

	err := c.MustBindWith(req, binding.FormPost)
	if err != nil {
		log.Printf("(ERR) %s: failed to bind request form: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}
	req.ExpertID = expertID
	req.AdminID = iUserID.(string)

	_, err = s.repository.CreateComment(&model.Comment{
		ID:       uuid.NewV4().String(),
		AdminID:  req.AdminID,
		ExpertID: req.ExpertID,
		Body:     req.Body,
	})
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

	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/expert/%s", expertID))
}

func (s *apiserver) webAdminUploadExpertDocument(c *gin.Context) {
	const logPref = "webAdminUploadExpertDocument"

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	expertID := c.Param("expertId")
	if strings.TrimSpace(expertID) == "" {
		log.Printf("(ERR) %s: empty expertId", logPref)
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}

	expert, err := s.repository.GetExpert(expertID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: expert [%s] not found", logPref, expertID)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get expert [%s]: %v", logPref, expertID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	lendoc := len(expert.DocumentURLs)
	if lendoc >= maxExpertDocumentCount {
		log.Printf("(ERR) %s: expert [%s] already processed", logPref, expertID)
		c.AbortWithStatusJSON(http.StatusConflict, representation.ErrExpertAlreadyProcessed)
		return
	}

	c.Request.Header.Add("Content-Type", "multipart/form-data")
	filename, err := s.saveUploadedFile(c, expertID, lendoc)
	if err != nil {
		log.Printf("(ERR) %s: saveUploadedFile error: %v", logPref, err)
		return
	}

	expert.DocumentURLs = append(expert.DocumentURLs, s.generateDocumentURL(expertID, filename))
	_, err = s.repository.UpdateExpertDocumentURLs(expert)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: expert [%s] failed to upload doc urls: %v", logPref, expertID, err)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to update expert [%s] document urls: %v", logPref, expertID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "expert_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      iRole.(string),
				"status":    iStatus.(string),
			},
			"expert": representation.ExpertViewItemPersistenceToAPI(expert),
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
		log.Printf("(ERR) Failed to delete expert: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func (s *apiserver) webAdminExpertProfile(c *gin.Context) {
	const logPref = "webAdminExpertItem"

	iRole, ok := c.Get("role")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	iUserID, ok := c.Get("userID")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	var (
		role     = iRole.(string)
		expertID = iUserID.(string)
	)

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

	iStatus, ok := c.Get("status")
	if !ok {
		c.Redirect(http.StatusFound, "/admin/login")
		return
	}

	c.HTML(http.StatusOK, "expert_item.html",
		gin.H{
			"metadata": gin.H{
				"logged_in": true,
				"role":      role,
				"status":    iStatus.(string),
			},
			"expert": representation.ExpertViewItemPersistenceToAPI(expertRes),
		},
	)
}

func (s *apiserver) webAdminExpertDocumentDelete(c *gin.Context) {
	const logPref = "webAdminExpertDocumentDelete"

	expertID := c.Param("expertId")
	documentID := c.Param("documentId")

	expert, err := s.repository.GetExpert(expertID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: expert [%s] not found", logPref, expertID)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to get expert [%s]: %v", logPref, expertID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	for i, docURL := range expert.DocumentURLs {
		if strings.HasPrefix(filepath.Base(docURL), documentID) {
			expert.DocumentURLs[i] = expert.DocumentURLs[len(expert.DocumentURLs)-1]
			expert.DocumentURLs[len(expert.DocumentURLs)-1] = ""
			expert.DocumentURLs = expert.DocumentURLs[:len(expert.DocumentURLs)-1]
			break
		}
	}

	_, err = s.repository.UpdateExpertDocumentURLs(expert)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			log.Printf("(ERR) %s: expert [%s] failed to upload doc urls: %v", logPref, expertID, err)
			c.AbortWithStatusJSON(http.StatusNotFound, representation.ErrNotFound)
			return
		}
		log.Printf("(ERR) %s: failed to update expert [%s] document urls: %v", logPref, expertID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	wildcard := filepath.Join(s.documentHostDir, expertID, fmt.Sprintf("%s.%s", documentID, "*"))

	if err := removeFilesWildcard(wildcard); err != nil {
		log.Printf("(ERR) %s: failed to remove expert document wildcard [%s]: %v", logPref, wildcard, err)
	}

	c.Status(http.StatusOK)
}

func removeFilesWildcard(wildcard string) error {
	files, err := filepath.Glob(wildcard)
	if err != nil {
		return err
	}
	for _, f := range files {
		_ = os.Remove(f)
	}
	return nil
}
