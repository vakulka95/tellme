package app

const (
	UserRoleAdmin  = "admin"
	UserRoleExpert = "expert"

	authCookieKey          = "fight_covid_auth"
	maxExpertDocumentCount = 10
	maxExpertDocumentSize  = 8 << 20 // 8 MiB

	maxRequisitionSessionCount = 3
	minRequisitionSessionCount = 0
)

var allowedExtensions = map[string]string{
	"image/jpg":          "jpg",
	"image/jpeg":         "jpeg",
	"image/png":          "png",
	"application/pdf":    "pdf",
	"application/msword": "doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "docx",
}
