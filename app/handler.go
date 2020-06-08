package app

import (
	"path"
	"strings"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"

	"github.com/gin-gonic/gin"
)

func (s *apiserver) registerHandlers() {
	gin.SetMode(gin.ReleaseMode)

	s.engine = gin.New()
	s.engine.MaxMultipartMemory = maxExpertDocumentSize
	s.engine.Use(gin.Logger(), gin.Recovery())

	//
	// Web admin config
	//
	s.engine.LoadHTMLGlob(path.Join(s.config.TemplatesStaticFilesDir, "/*"))

	authentication := s.authenticationInterceptor()
	allAuthorization := s.authorizationInterceptor(UserRoleExpert, UserRoleAdmin)
	adminAuthorization := s.authorizationInterceptor(UserRoleAdmin)
	expertAuthorization := s.authorizationInterceptor(UserRoleExpert)

	statusActiveAuthentication := s.checkStatusInterceptor(model.ExpertStatusActive)
	statusAllAuthentication := s.checkStatusInterceptor(model.ExpertStatusActive, model.ExpertStatusBlocked, model.ExpertStatusOnReview)

	//
	// Admin Auth
	//
	s.engine.GET("/admin", authentication, s.webIndexPage)
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
	s.engine.DELETE("/admin/expert/:expertId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminExpertDelete)

	//
	// Admin Requisition
	//
	s.engine.GET("/admin/requisition", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionList)
	s.engine.GET("/admin/requisition/:requisitionId", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionItem)
	s.engine.PUT("/admin/requisition/:requisitionId/take", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionTake)
	s.engine.PUT("/admin/requisition/:requisitionId/discard", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionDiscard)
	s.engine.PUT("/admin/requisition/:requisitionId/complete", authentication, allAuthorization, statusActiveAuthentication, s.webAdminRequisitionComplete)

	//
	// Admin Review
	//
	s.engine.GET("/admin/review", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminReviewList)
	s.engine.GET("/admin/review/:reviewId", authentication, adminAuthorization, statusActiveAuthentication, s.webAdminReviewItem)

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
