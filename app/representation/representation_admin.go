package representation

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
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
			"status":           v.Status,
			"updated_at":       v.UpdatedAt.String(),
			"created_at":       v.CreatedAt.String(),
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
			"updated_at":            v.UpdatedAt.String(),
			"created_at":            v.CreatedAt.String(),
		}
	}
	return gin.H{"items": items, "total": r.Total}
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}
