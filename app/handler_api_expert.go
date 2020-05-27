package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

func (s *apiserver) apiExpertRegister(c *gin.Context) {
	const logPref = "apiExpertRegister"

	var (
		err       error
		expertReq = &representation.CreateExpertRequest{}
	)

	if err = c.ShouldBindJSON(expertReq); err != nil {
		log.Printf("(ERR) %s: failed to read body: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, representation.ErrInvalidRequest)
		return
	}

	expert, err := representation.ExpertAPItoPersistence(expertReq)
	if err != nil {
		log.Printf("(ERR) %s: failed to read validate expert [%s]: %v", logPref, expertReq.Email, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return
	}

	expertRes, err := s.repository.CreateExpert(expert)
	if err != nil {
		if strings.Contains(err.Error(), "constraint") {
			log.Printf("(ERR) %s: failed to create expert [%s]: %v", logPref, expertReq.Email, err)
			c.AbortWithStatusJSON(http.StatusConflict, representation.ErrConflict)
			return
		}
		log.Printf("(ERR) %s: failed to create expert [%s]: %v", logPref, expertReq.Email, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	log.Printf("(INFO) %s: expert [%s] successfully created", logPref, expertReq.Email)
	c.JSON(http.StatusCreated, representation.CreateExpertResponse{ID: expertRes.ID})
}

func (s *apiserver) apiUploadExpertDoc(c *gin.Context) {
	const logPref = "apiUploadExpertDoc"

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

	if expert.Status != model.ExpertStatusOnReview {
		log.Printf("(ERR) %s: expert [%s] already processed", logPref, expertID)
		c.AbortWithStatusJSON(http.StatusConflict, representation.ErrExpertAlreadyProcessed)
		return
	}

	lendoc := len(expert.DocumentURLs)
	if lendoc >= maxExpertDocumentCount {
		log.Printf("(ERR) %s: expert [%s] already processed", logPref, expertID)
		c.AbortWithStatusJSON(http.StatusConflict, representation.ErrExpertAlreadyProcessed)
		return
	}

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

	log.Printf("(INFO) %s: expert [%s] docs successfully uploaded", logPref, expertID)
	c.Status(http.StatusAccepted)
}

func (s *apiserver) apiValidateExpertPhone(c *gin.Context) {
	const logPref = "apiValidateExpertPhone"

	var (
		err       error
		expertReq = &representation.ValidateExpertPhoneRequest{}
	)

	if err = c.ShouldBindJSON(expertReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, representation.ErrInvalidRequest)
		return
	}

	_, err = s.repository.GetExpertByPhone(strings.TrimPrefix(expertReq.Phone, "+"))
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.Status(http.StatusOK)
			return
		}
		log.Printf("(ERR) %s: ailed to find expert by phone: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	c.Status(http.StatusConflict)
}

func (s *apiserver) apiValidateExpertEmail(c *gin.Context) {
	const logPref = "apiValidateExpertEmail"

	var (
		err       error
		expertReq = &representation.ValidateExpertEmailRequest{}
	)

	if err = c.ShouldBindJSON(expertReq); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, representation.ErrInvalidRequest)
		return
	}

	_, err = s.repository.GetExpertByEmail(strings.ToLower(expertReq.Email))
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			c.Status(http.StatusOK)
			return
		}
		log.Printf("(ERR) %s: failed to find expert by email: %v", logPref, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, representation.ErrInternal)
		return
	}

	c.Status(http.StatusConflict)
}

func (s *apiserver) apiGetDiagnosisList(c *gin.Context) {
	c.JSON(http.StatusOK, representation.Diagnoses)
}

func (s *apiserver) saveUploadedFile(c *gin.Context, expertID string, lendoc int) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return "", fmt.Errorf("failed to read from file: %v", err)
	}

	srcFile, err := file.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer srcFile.Close()

	srcBytes, err := ioutil.ReadAll(srcFile)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	mimeType := http.DetectContentType(srcBytes)
	ext, ok := allowedExtensions[mimeType]
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return "", fmt.Errorf("image mime type %s is not supported", mimeType)
	}

	dir := filepath.Join(s.documentHostDir, expertID)
	_ = os.MkdirAll(dir, os.ModePerm)

	filename := fmt.Sprintf("%d.%s", lendoc+1, ext)
	fullFilename := filepath.Join(dir, filename)

	out, err := os.Create(fullFilename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	err = ioutil.WriteFile(fullFilename, srcBytes, os.ModePerm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, representation.ErrInvalidRequest)
		return "", fmt.Errorf("failed to write file file: %v", err)
	}

	return filename, nil
}
