package representation

const (
	GenderMale     = "male"
	GenderFemale   = "female"
	GenderNoMatter = "no_matter"

	FeedbackTypePhone    = "phone"
	FeedbackTypeViber    = "viber"
	FeedbackTypeSkype    = "skype"
	FeedbackTypeTelegram = "telegram"
	FeedbackTypeZoom     = "zoom"
)

var allowedGender = map[string]bool{
	GenderMale:     true,
	GenderFemale:   true,
	GenderNoMatter: true,
}

var allowedFeedbackType = map[string]bool{
	FeedbackTypePhone:    true,
	FeedbackTypeViber:    true,
	FeedbackTypeSkype:    true,
	FeedbackTypeTelegram: true,
	FeedbackTypeZoom:     true,
}

type CreateExpertRequest struct {
	Name            string   `json:"name" description:"Name and surname in one field. Validation: not be empty" maximum:"256"`
	Gender          string   `json:"gender" description:"Validation: 'male' or 'female'"`
	Phone           string   `json:"phone" description:"Validation: required 10 digits"`
	Email           string   `json:"email" description:"Validation: required one @ symbol" maximum:"128"`
	Password        string   `json:"password" description:"Validation: minimum 7 symbols, maximum: no limit, at least 1 digit, at least 1 upper case"`
	Specializations []string `json:"specializations" description:"Validation: only allowed keys from /api/v1/diagnosis endpoint"`
	Education       string   `json:"education" description:"Validation: not be empty" maximum:"256"`
}

type CreateExpertResponse struct {
	ID string `json:"id" description:"Just expert's identifier"`
}

type ValidateExpertPhoneRequest struct {
	Phone string `json:"phone" description:"Validation: required 10 digits"`
}

type ValidateExpertEmailRequest struct {
	Email string `json:"email" description:"Validation: required one @ symbol" maximum:"128"`
}

type CreateRequisitionRequest struct {
	Name                 string `json:"name" description:"Name and surname in one field. Validation: not be empty" maximum:"256"`
	Gender               string `json:"gender" description:"Validation: 'male' or 'female'"`
	Phone                string `json:"phone" description:"Validation: required 10 digits"`
	Diagnosis            string `json:"diagnosis" description:"Validation: only allowed keys from /api/v1/diagnosis endpoint"`
	DiagnosisDescription string `json:"diagnosisDescription" description:"Just diagnosis description text"`
	ExpertGender         string `json:"expertGender" description:"Validation: 'male', 'female', 'no_matter'"`
	FeedbackType         string `json:"feedbackType" description:"Validation: 'skype', 'zoom', 'phone', 'viber', 'telegram' allowed only"`
	FeedbackContact      string `json:"feedbackContact" description:"Validation: could be empty only if 'phone' feedback type specified" maximum:"128"`
	FeedbackTime         string `json:"feedbackTime" description:"Validation: 8:00, 13:00, 16:00"`
	FeedbackWeekDay      string `json:"feedbackWeekDay" description:"Validation: Day of week: mon, tue, wed, thu, fri, sat, sun"`
	IsAdult              bool   `json:"isAdult" description:"Validation: must be 'true'"`
}

type CreateRequisitionResponse struct {
	ID string `json:"id" description:"Just requisition's identifier"`
}

type ConfirmReviewRequest struct {
	PlatformReview     string `json:"platform_review" description:"Platform review"`
	ConsultationCount  int    `json:"consultation_count" description:"Consultation count"`
	ConsultationReview string `json:"consultation_review" description:"Consultation satisfaction"`
	ExpertPoint        int    `json:"expert_point" description:"Expert point: min - 1, max - 5"`
	ExpertReview       string `json:"expert_review" description:"Expert work review"`
	Token              string `json:"token" description:"Token to authorize review request"`
}

type ConfirmReviewResponse struct {
	ID string `json:"id" description:"Just review's identifier"`
}
