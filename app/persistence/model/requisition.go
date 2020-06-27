package model

import (
	"strings"
	"time"

	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"

	uuid "github.com/satori/go.uuid"
)

const (
	RequisitionStatusCreated    = "created"
	RequisitionStatusProcessing = "processing"
	RequisitionStatusCompleted  = "completed"
)

type Requisition struct {
	ID                   string    `db:"id"`
	ExpertID             string    `db:"expert_id"`
	Username             string    `db:"username"`
	Gender               string    `db:"gender"`
	Phone                string    `db:"phone"`
	Diagnosis            string    `db:"diagnosis"`
	DiagnosisDescription string    `db:"diagnosis_description"`
	ExpertGender         string    `db:"expert_gender"`
	FeedbackType         string    `db:"feedback_type"`
	FeedbackContact      string    `db:"feedback_contact"`
	FeedbackTime         string    `db:"feedback_time"`
	FeedbackWeekDay      string    `db:"feedback_week_day"`
	IsAdult              bool      `db:"is_adult"`
	SessionCount         int       `db:"session_count"`
	SMSReplyCount        int       `db:"sms_reply_count"`
	Status               string    `db:"status"`
	UpdatedAt            time.Time `db:"updated_at"`
	CreatedAt            time.Time `db:"created_at"`
}

type QueryRequisitionList struct {
	Limit           int        `db:"limit"`
	Offset          int        `db:"offset"`
	Status          string     `db:"status"`
	ExpertID        string     `db:"expert_id"`
	FeedbackTime    string     `db:"feedback_time"`
	FeedbackWeekDay string     `db:"feedback_week_day"`
	Specializations []string   `db:"specializations"`
	Search          string     `db:"search"`
	CreatedAtFrom   *time.Time `db:"created_at"`
	CreatedAtTo     *time.Time `db:"created_at"`
}

func DiscoverRequisitionExpression(search string) postgres.Expression {

	_, err := uuid.FromString(search)
	if err == nil {
		return postgres.NewExpression("id", postgres.NewString(search), postgres.OperatorEqual)
	}

	phone := strings.TrimPrefix(search, "+")
	phone = strings.TrimPrefix(phone, "3")
	phone = strings.TrimPrefix(phone, "8")
	if strings.IndexFunc(phone, isNotDigit) == -1 {
		return postgres.NewExpression("phone", postgres.NewString(phone), postgres.OperatorLike)
	}

	return postgres.NewExpression("lower(username)", postgres.NewString(strings.ToLower(search)), postgres.OperatorLike)
}

type RequisitionList struct {
	Items []*Requisition `db:"items"`
	Total int            `db:"total"`
}
