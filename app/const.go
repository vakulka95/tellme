package app

const (
	UserRoleAdmin  = "admin"
	UserRoleExpert = "expert"

	authCookieKey          = "fight_covid_auth"
	maxExpertDocumentCount = 5
	maxExpertDocumentSize  = 8 << 20 // 8 MiB

	maxRequisitionSessionCount = 3
)

var allowedExtensions = map[string]string{
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
	"image/png":  "png",
}
