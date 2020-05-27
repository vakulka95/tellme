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

func RequisitionListPersistenceToAPI(r *model.RequisitionList) gin.H {
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
			"status":                v.Status,
			"feedback_week_day":     week[v.FeedbackWeekDay],
			"feedback_time":         v.FeedbackTime,
			"updated_at":            v.UpdatedAt.Format(timestampLayout),
			"created_at":            v.CreatedAt.Format(timestampLayout),
		}
	}
	return gin.H{"items": items, "total": r.Total}
}

func RequisitionItemPersistenceToAPI(r *model.Requisition) gin.H {
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
		"status":                r.Status,
		"feedback_week_day":     week[r.FeedbackWeekDay],
		"feedback_time":         r.FeedbackTime,
		"updated_at":            r.UpdatedAt.Format(timestampLayout),
		"created_at":            r.CreatedAt.Format(timestampLayout),
	}
}

func ExpertViewItemPersistenceToAPI(v *model.Expert) gin.H {
	specs := make([]string, len(v.Specializations))
	for j, k := range v.Specializations {
		specs[j] = Diagnoses[k]
	}
	return gin.H{
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
		"token":               v.Token, // TODO: remove
		"status":              v.Status,
		"updated_at":          v.UpdatedAt.Format(timestampLayout),
		"created_at":          v.CreatedAt.Format(timestampLayout),
	}
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}
