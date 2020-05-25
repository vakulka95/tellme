package app

import (
	"path"
	"strings"

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

	s.engine.GET("/admin", authentication, s.webIndexPage)
	s.engine.POST("/admin/login", s.webAdminLogin)
	s.engine.GET("/admin/login", s.webAdminGetLoginPage)

	s.engine.GET("/admin/logout", s.webAdminLogout)

	s.engine.GET("/admin/expert", authentication, adminAuthorization, s.webAdminExpertList)
	s.engine.POST("/admin/expert", authentication, adminAuthorization, s.webAdminExpertList)

	s.engine.PUT("/admin/expert/:expertId/block", authentication, adminAuthorization, s.webAdminExpertBlock)
	s.engine.PUT("/admin/expert/:expertId/activate", authentication, adminAuthorization, s.webAdminExpertActivate)
	s.engine.PUT("/admin/expert/:expertId/password", authentication, adminAuthorization, s.webAdminUpdateExpertPassword)
	s.engine.DELETE("/admin/expert/:expertId", authentication, adminAuthorization, s.webAdminExpertDelete)

	s.engine.GET("/admin/requisition", authentication, allAuthorization, s.webAdminRequisitionList)
	s.engine.POST("/admin/requisition", authentication, allAuthorization, s.webAdminRequisitionList)
	s.engine.PUT("/admin/requisition/:requisitionId/take", authentication, allAuthorization, s.webAdminRequisitionTake)
	s.engine.PUT("/admin/requisition/:requisitionId/complete", authentication, allAuthorization, s.webAdminRequisitionComplete)
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

	//
	// Service API
	//
	s.engine.GET("/service/v1/pg_stat", s.serviceDatabaseStat)
}

func (s *apiserver) cors(c *gin.Context) {
	methods := strings.Join([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH"}, ",")
	headers := strings.Join([]string{"Content-Type", "Accept", "Authorization", "Language", "Set-Cookie"}, ",")

	c.Header("Allow", methods)
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", headers)
	c.Header("Access-Control-Allow-Methods", methods)
}
