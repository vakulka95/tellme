package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/tidwall/pretty"

	"gitlab.com/tellmecomua/tellme.api/app/representation"
)

func main() {
	(&swagger{ws: &restful.WebService{}}).genSwagger()
}

type swagger struct {
	ws        *restful.WebService
	container *restful.Container
}

func (s *swagger) apiExpertRegister(*restful.Request, *restful.Response)      {}
func (s *swagger) apiGetDiagnosisList(*restful.Request, *restful.Response)    {}
func (s *swagger) apiUploadExpertDoc(*restful.Request, *restful.Response)     {}
func (s *swagger) apiRequisitionCreate(*restful.Request, *restful.Response)   {}
func (s *swagger) apiValidateExpertPhone(*restful.Request, *restful.Response) {}
func (s *swagger) apiValidateExpertEmail(*restful.Request, *restful.Response) {}
func (s *swagger) apiConfirmReview(*restful.Request, *restful.Response)       {}

func (s *swagger) genSwagger() error {
	s.init()

	s.container = restful.NewContainer()
	s.container.Router(restful.CurlyRouter{})
	s.container.Add(s.ws)

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		CookiesAllowed: false,
		Container:      s.container,
	}

	s.container.Filter(cors.Filter)
	s.container.Filter(s.container.OPTIONSFilter)

	config := restfulspec.Config{
		WebServices: []*restful.WebService{s.ws},
		PostBuildSwaggerObjectHandler: func(swo *spec.Swagger) {
			swo.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:   "Fight COVID-19 in UA",
					Contact: &spec.ContactInfo{Name: "Fight COVID-19 in UA"},
					Version: "v0.1.0",
				},
			}
		},
	}

	bb, err := restfulspec.BuildSwagger(config).MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal swagger: %v", err)
	}
	return ioutil.WriteFile("./apidocs/swagger.json", pretty.Pretty(bb), os.ModePerm)
}

func (s *swagger) init() {

	s.ws.
		Path("/api/v1/").
		Consumes(restful.MIME_JSON, "multipart/form-data").
		Produces(restful.MIME_JSON)

	s.ws.Route(s.ws.GET("/diagnosis").
		To(s.apiGetDiagnosisList).
		Doc("Get diagnosis list. List is used for expert's specializations and requisition diagnosis. Key: EN tag, Value: UA description").
		Returns(http.StatusOK, "OK", representation.Diagnoses).
		Writes(representation.Diagnoses).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))

	s.ws.Route(s.ws.POST("/expert").
		To(s.apiExpertRegister).
		Doc("Request provide help. Register expert(psychologist)").
		Returns(http.StatusCreated, "Created", representation.CreateExpertResponse{}).
		Returns(http.StatusBadRequest, "BadRequest", representation.ErrInvalidRequest).
		Returns(http.StatusConflict, "Conflict", representation.ErrConflict).
		Returns(http.StatusUnprocessableEntity, "UnprocessableEntity", representation.ErrInvalidRequest).
		Returns(http.StatusInternalServerError, "InternalServerError", representation.ErrInternal).
		Reads(representation.CreateExpertRequest{}).
		Writes(representation.CreateExpertResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))

	s.ws.Route(s.ws.GET("/validate/expert/phone").
		To(s.apiValidateExpertPhone).
		Doc("Validate unique expert phone").
		Returns(http.StatusOK, "OK", nil).
		Returns(http.StatusConflict, "Conflict", representation.ErrConflict).
		Returns(http.StatusUnprocessableEntity, "UnprocessableEntity", representation.ErrInvalidRequest).
		Returns(http.StatusInternalServerError, "InternalServerError", representation.ErrInternal).
		Reads(representation.ValidateExpertPhoneRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))

	s.ws.Route(s.ws.GET("/validate/expert/email").
		To(s.apiValidateExpertEmail).
		Doc("Validate unique expert email").
		Returns(http.StatusOK, "OK", nil).
		Returns(http.StatusConflict, "Conflict", representation.ErrConflict).
		Returns(http.StatusUnprocessableEntity, "UnprocessableEntity", representation.ErrInvalidRequest).
		Returns(http.StatusInternalServerError, "InternalServerError", representation.ErrInternal).
		Reads(representation.ValidateExpertEmailRequest{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))

	s.ws.Route(s.ws.POST("/expert/{expertId}/document").
		To(s.apiUploadExpertDoc).
		Doc("Upload expert's *.png images of education. Restrictions: only png, up to 8Mb by file. Up to 5 files for one expert").
		Param(restful.PathParameter("expertId", "Expert's identifier")).
		Param(restful.FormParameter("image", "Expert's image in png format up to 8Mb size")).
		Returns(http.StatusAccepted, "Accepted", nil).
		Returns(http.StatusBadRequest, "BadRequest", representation.ErrInvalidRequest).
		Returns(http.StatusNotFound, "NotFound", representation.ErrNotFound).
		Returns(http.StatusConflict, "Conflict", representation.ErrExpertAlreadyProcessed).
		Returns(http.StatusInternalServerError, "InternalServerError", representation.ErrInternal).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))

	s.ws.Route(s.ws.POST("/requisition").
		To(s.apiRequisitionCreate).
		Doc("Create user requisition").
		Returns(http.StatusCreated, "Created", representation.CreateRequisitionResponse{}).
		Returns(http.StatusConflict, "Conflict", representation.ErrConflict).
		Returns(http.StatusBadRequest, "BadRequest", representation.ErrInvalidRequest).
		Returns(http.StatusUnprocessableEntity, "UnprocessableEntity", representation.ErrInvalidRequest).
		Returns(http.StatusInternalServerError, "InternalServerError", representation.ErrInternal).
		Reads(representation.CreateRequisitionRequest{}).
		Writes(representation.CreateRequisitionResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))

	s.ws.Route(s.ws.POST("/review").
		To(s.apiConfirmReview).
		Doc("Confirm user review about expert").
		Returns(http.StatusCreated, "Created", representation.ConfirmReviewResponse{}).
		Returns(http.StatusConflict, "Conflict", representation.ErrConflict).
		Returns(http.StatusNotFound, "NotFound", representation.ErrNotFound).
		Returns(http.StatusBadRequest, "BadRequest", representation.ErrInvalidRequest).
		Returns(http.StatusUnprocessableEntity, "UnprocessableEntity", representation.ErrInvalidRequest).
		Returns(http.StatusInternalServerError, "InternalServerError", representation.ErrInternal).
		Reads(representation.ConfirmReviewRequest{}).
		Writes(representation.ConfirmReviewResponse{}).
		Metadata(restfulspec.KeyOpenAPITags, []string{"main"}))
}
