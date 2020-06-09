package representation

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

const (
	timestampLayout = "2006-01-02 15:04"
)

var week = map[string]string{
	"mon": "Понеділок",
	"tue": "Вівторок",
	"wed": "Середа",
	"thu": "Четвер",
	"fri": "П'ятниця",
	"sat": "Субота",
	"sun": "Неділя",
}

var feedbackTime = map[string]string{
	"8:00":  "8:00-13:00",
	"13:00": "13:00-18:00",
	"18:00": "18:00-21:00",
}

type Expert struct {
	ID              string   `json:"id"`
	Username        string   `json:"username"`
	Gender          string   `json:"gender"`
	Age             int      `json:"age"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	Specializations []string `json:"specializations"`
	Education       string   `json:"education"`
	DocumentURLs    []string `json:"document_urls"`
	Status          string   `json:"status"`
	UpdatedAt       string   `json:"updated_at"`
	CreatedAt       string   `json:"created_at"`
}

type ExpertList struct {
	Items []gin.H `json:"items"`
	Total int     `json:"total"`
}

func ExpertListPersistenceToAPI(r *model.ExpertList) gin.H {
	items := make([]gin.H, len(r.Items))
	for i, v := range r.Items {

		specs := make([]string, len(v.Specializations))
		for j, k := range v.Specializations {
			specs[j] = Diagnoses[k]
		}

		items[i] = gin.H{
			"id":               v.ID,
			"username":         v.Username,
			"gender":           v.Gender,
			"phone":            v.Phone,
			"email":            v.Email,
			"specializations":  specs,
			"education":        v.Education,
			"document_urls":    v.DocumentURLs,
			"processing_count": v.ProcessingCount,
			"completed_count":  v.CompletedCount,
			"review_count":     v.ReviewCount,
			"session_count":    v.SessionCount,
			"status":           v.Status,
			"updated_at":       v.UpdatedAt.Format(timestampLayout),
			"created_at":       v.CreatedAt.Format(timestampLayout),
		}
	}
	return gin.H{"items": items, "total": r.Total}
}

type RequisitionList struct {
	Items []gin.H `json:"items"`
	Total int     `json:"total"`
}

func RequisitionListPersistenceToAPIForExpert(r *model.RequisitionList, userID string) gin.H {
	items := make([]gin.H, len(r.Items))
	for i, v := range r.Items {

		phone := ""
		username := ""
		feedbackContact := ""

		if v.ExpertID == userID {
			phone = v.Phone
			username = v.Username
			feedbackContact = v.FeedbackContact
		}

		items[i] = gin.H{
			"id":                    v.ID,
			"username":              username,
			"gender":                v.Gender,
			"phone":                 phone,
			"diagnosis":             DiagnosesStripped[v.Diagnosis],
			"diagnosis_description": v.DiagnosisDescription,
			"expert_gender":         v.ExpertGender,
			"feedback_type":         v.FeedbackType,
			"feedback_contact":      feedbackContact,
			"session_count":         v.SessionCount,
			"status":                v.Status,
			"feedback_week_day":     week[v.FeedbackWeekDay],
			"feedback_time":         feedbackTime[v.FeedbackTime],
			"updated_at":            v.UpdatedAt.Format(timestampLayout),
			"created_at":            v.CreatedAt.Format(timestampLayout),
		}
	}
	return gin.H{"items": items, "total": r.Total}
}

func RequisitionListPersistenceToAPIForAdmin(r *model.RequisitionList) gin.H {
	items := make([]gin.H, len(r.Items))
	for i, v := range r.Items {
		items[i] = gin.H{
			"id":                    v.ID,
			"username":              v.Username,
			"gender":                v.Gender,
			"phone":                 v.Phone,
			"diagnosis":             DiagnosesStripped[v.Diagnosis],
			"diagnosis_description": v.DiagnosisDescription,
			"expert_gender":         v.ExpertGender,
			"feedback_type":         v.FeedbackType,
			"feedback_contact":      v.FeedbackContact,
			"session_count":         v.SessionCount,
			"status":                v.Status,
			"feedback_week_day":     week[v.FeedbackWeekDay],
			"feedback_time":         feedbackTime[v.FeedbackTime],
			"updated_at":            v.UpdatedAt.Format(timestampLayout),
			"created_at":            v.CreatedAt.Format(timestampLayout),
		}
	}
	return gin.H{"items": items, "total": r.Total}
}

func RequisitionItemPersistenceToAPIForAdmin(r *model.Requisition) gin.H {
	return gin.H{
		"id":                    r.ID,
		"expert_id":             r.ExpertID,
		"username":              r.Username,
		"gender":                r.Gender,
		"phone":                 r.Phone,
		"diagnosis":             DiagnosesStripped[r.Diagnosis],
		"diagnosis_description": r.DiagnosisDescription,
		"expert_gender":         r.ExpertGender,
		"feedback_type":         r.FeedbackType,
		"feedback_contact":      r.FeedbackContact,
		"session_count":         r.SessionCount,
		"status":                r.Status,
		"feedback_week_day":     week[r.FeedbackWeekDay],
		"feedback_time":         feedbackTime[r.FeedbackTime],
		"updated_at":            r.UpdatedAt.Format(timestampLayout),
		"created_at":            r.CreatedAt.Format(timestampLayout),
	}
}

func RequisitionItemPersistenceToAPIForExpert(r *model.Requisition, userID string) gin.H {
	phone := ""
	username := ""
	feedbackContact := ""

	if r.ExpertID == userID {
		phone = r.Phone
		username = r.Username
		feedbackContact = r.FeedbackContact
	}

	return gin.H{
		"id":                    r.ID,
		"expert_id":             r.ExpertID,
		"username":              username,
		"gender":                r.Gender,
		"phone":                 phone,
		"diagnosis":             DiagnosesStripped[r.Diagnosis],
		"diagnosis_description": r.DiagnosisDescription,
		"expert_gender":         r.ExpertGender,
		"feedback_type":         r.FeedbackType,
		"feedback_contact":      feedbackContact,
		"session_count":         r.SessionCount,
		"status":                r.Status,
		"feedback_week_day":     week[r.FeedbackWeekDay],
		"feedback_time":         feedbackTime[r.FeedbackTime],
		"updated_at":            r.UpdatedAt.Format(timestampLayout),
		"created_at":            r.CreatedAt.Format(timestampLayout),
	}
}

func ExpertViewItemPersistenceToAPI(v *model.Expert) gin.H {
	return gin.H{
		"id":               v.ID,
		"username":         v.Username,
		"gender":           v.Gender,
		"phone":            v.Phone,
		"email":            v.Email,
		"specializations":  GenerateDiagnosesOptions(v.Specializations),
		"education":        v.Education,
		"document_urls":    v.DocumentURLs,
		"processing_count": v.ProcessingCount,
		"completed_count":  v.CompletedCount,
		"review_count":     v.ReviewCount,
		"session_count":    v.SessionCount,
		"status":           v.Status,
		"updated_at":       v.UpdatedAt.Format(timestampLayout),
		"created_at":       v.CreatedAt.Format(timestampLayout),
	}
}

func ReviewListPersistenceToAPI(r *model.ReviewList) gin.H {
	items := make([]gin.H, len(r.Items))
	for i, v := range r.Items {
		items[i] = gin.H{
			"id":                  v.ID,
			"expert_id":           v.ExpertID,
			"expert_username":     v.ExpertUsername,
			"requisition_id":      v.RequisitionID,
			"platform_review":     v.PlatformReview,
			"consultation_count":  v.ConsultationCount,
			"consultation_review": v.ConsultationReview,
			"expert_point":        v.ExpertPoint,
			"expert_review":       v.ExpertReview,
			"status":              v.Status,
			"updated_at":          v.UpdatedAt.Format(timestampLayout),
			"created_at":          v.CreatedAt.Format(timestampLayout),
		}
	}
	return gin.H{"items": items, "total": r.Total}
}

func ReviewItemPersistenceToAPI(v *model.Review) gin.H {
	return gin.H{
		"id":                  v.ID,
		"expert_id":           v.ExpertID,
		"requisition_id":      v.RequisitionID,
		"platform_review":     v.PlatformReview,
		"consultation_count":  v.ConsultationCount,
		"consultation_review": v.ConsultationReview,
		"expert_point":        v.ExpertPoint,
		"expert_review":       v.ExpertReview,
		"status":              v.Status,
		"updated_at":          v.UpdatedAt.Format(timestampLayout),
		"created_at":          v.CreatedAt.Format(timestampLayout),
	}
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

type UpdateExpertRequest struct {
	ID              string   `form:"-"`
	Username        string   `form:"username"`
	Phone           string   `form:"phone"`
	Gender          string   `form:"gender"`
	Education       string   `form:"education"`
	Specializations []string `form:"specializations"`
}
