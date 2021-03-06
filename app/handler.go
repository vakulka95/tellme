package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/pretty"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/pkg/swagger"
)

func (s *apiserver) registerHandlers() {
	gin.SetMode(gin.ReleaseMode)

	s.engine = gin.New()
	s.engine.MaxMultipartMemory = maxExpertDocumentSize
	s.engine.Use(gin.Logger(), gin.Recovery())

	//
	// Web admin config
	//
	s.engine.LoadHTMLGlob(path.Join(s.config.StaticFilesDir, "/templates/*"))
	s.engine.Static("/static/styles", path.Join(s.config.StaticFilesDir, "/styles"))
	s.engine.Static("/document/expert", s.config.ExpertDocumentsStoreDir)

	authentication := s.authenticationInterceptor()
	allAuthorization := s.authorizationInterceptor(UserRoleExpert, UserRoleAdmin)
	adminAuthorization := s.authorizationInterceptor(UserRoleAdmin)
	expertAuthorization := s.authorizationInterceptor(UserRoleExpert)

	statusActiveAuthentication := s.checkStatusInterceptor(model.ExpertStatusActive)
	statusAllAuthentication := s.checkStatusInterceptor(model.ExpertStatusActive, model.ExpertStatusBlocked, model.ExpertStatusOnReview)

	//
	// Admin Auth
	//
	s.engine.GET("/admin/", authentication, s.webIndexPage)
	s.engine.POST("/admin/login", s.webAdminLogin)
	s.engine.GET("/admin/login", s.webAdminGetLoginPage)
	s.engine.GET("/admin/logout", s.webAdminLogout)

	//
	// Admin Expert
	//
	s.engine.GET("/admin/expert", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertList)
	s.engine.POST("/admin/expert", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertList)
	s.engine.GET("/admin/expert/:expertId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertItem)
	s.engine.POST("/admin/expert/:expertId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminUpdateExpertItem)
	s.engine.POST("/admin/expert/:expertId/document", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminUploadExpertDocument)
	s.engine.GET("/admin/profile", authentication, expertAuthorization, statusAllAuthentication, s.webAdminExpertProfile)

	s.engine.PUT("/admin/expert/:expertId/block", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertBlock)
	s.engine.PUT("/admin/expert/:expertId/activate", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertActivate)
	s.engine.PUT("/admin/expert/:expertId/password", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminUpdateExpertPassword)
	s.engine.POST("/admin/expert/:expertId/comment", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertCommentCreate)
	s.engine.DELETE("/admin/expert/:expertId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertDelete)
	s.engine.DELETE("/admin/expert/:expertId/document/:documentId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertDocumentDelete)

	s.engine.GET("/admin/expert_rating", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertRatingList)
	s.engine.POST("/admin/expert_rating", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertRatingList)

	s.engine.GET("/admin/expert_rating/pdf", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertRatingPDF)
	s.engine.GET("/admin/expert_rating/excel", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertRatingExcel)

	//
	// Admin Requisition
	//
	s.engine.GET("/admin/requisition", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionList)
	s.engine.GET("/admin/requisition/:requisitionId", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionItem)
	s.engine.PUT("/admin/requisition/:requisitionId/take", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionTake)
	s.engine.PUT("/admin/requisition/:requisitionId/discard", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionDiscard)
	s.engine.PUT("/admin/requisition/:requisitionId/complete", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionComplete)
	s.engine.PUT("/admin/requisition/:requisitionId/no_answer", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionNoAnswer)
	s.engine.DELETE("/admin/requisition/:requisitionId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminRequisitionDelete)

	//
	// Admin Review
	//
	s.engine.GET("/admin/review", authentication, allAuthorization, s.webAdminReviewList)
	s.engine.GET("/admin/review/:reviewId", authentication, allAuthorization, s.webAdminReviewItem)

	//
	// Admin Session
	//
	s.engine.POST("/admin/requisition/:requisitionId/session/apply", authentication, expertAuthorization, statusActiveAuthentication, s.webAdminSessionCreate)
	s.engine.DELETE("/admin/requisition/:requisitionId/session/discard", authentication, expertAuthorization, statusActiveAuthentication, s.webAdminSessionDiscard)

	//
	// Main API
	//
	s.engine.GET("/api/v1/diagnosis", s.cors, s.apiGetDiagnosisList)
	s.engine.OPTIONS("/api/v1/diagnosis", s.cors)

	s.engine.POST("/api/v1/expert", s.cors, s.apiExpertRegister)
	s.engine.OPTIONS("/api/v1/expert", s.cors)

	s.engine.POST("/api/v1/expert/:expertId/document", s.cors, s.apiUploadExpertDoc)
	s.engine.OPTIONS("/api/v1/expert/:expertId/document", s.cors)

	s.engine.POST("/api/v1/validate/expert/phone", s.cors, s.apiValidateExpertPhone)
	s.engine.OPTIONS("/api/v1/validate/expert/phone", s.cors)

	s.engine.POST("/api/v1/validate/expert/email", s.cors, s.apiValidateExpertEmail)
	s.engine.OPTIONS("/api/v1/validate/expert/email", s.cors)

	s.engine.POST("/api/v1/requisition", s.cors, s.apiRequisitionCreate)
	s.engine.OPTIONS("/api/v1/requisition", s.cors)

	s.engine.POST("/api/v1/review", s.cors, s.apiReviewConfirm)
	s.engine.OPTIONS("/api/v1/review", s.cors)

	//
	// Service API
	//
	s.engine.GET("/service/v1/pg_stat", s.serviceDatabaseStat)
	s.engine.POST("/service/v1/send_not_reviewed_requisition", s.serviceSendRequisitionReview)
}

func (s *apiserver) cors(c *gin.Context) {
	methods := strings.Join([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH"}, ",")
	headers := strings.Join([]string{"Content-Type", "Accept", "Authorization", "Language", "Set-Cookie"}, ",")

	c.Header("Allow", methods)
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", headers)
	c.Header("Access-Control-Allow-Methods", methods)
}

func (s *apiserver) registerSwaggerApidocs() error {
	const swaggerDir = "/etc/tellme.api/swagger"

	bb, err := swagger.GetApidocsJSON()
	if err != nil {
		return err
	}

	_ = os.MkdirAll(swaggerDir, os.ModePerm)
	filepath := path.Join(swaggerDir, "/apidocs.json")

	err = ioutil.WriteFile(filepath, pretty.Pretty(bb), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	s.engine.StaticFile("/api/v1/apidocs.json", filepath)
	return nil
}
