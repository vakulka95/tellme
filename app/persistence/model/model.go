package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

const (
	StatusActive   = "active"
	StatusBlocked  = "blocked"
	StatusOnReview = "on_review"

	StatusCreated    = "created"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)

type Expert struct {
	ID              string         `db:"id"`
	Username        string         `db:"username"`
	Gender          string         `db:"gender"`
	Phone           string         `db:"phone"`
	Email           string         `db:"email"`
	Password        string         `db:"password"`
	Specializations pq.StringArray `db:"specializations"`
	Education       string         `db:"education"`
	DocumentURLs    pq.StringArray `db:"document_urls"`
	ProcessingCount int            `db:"processing_count"`
	CompletedCount  int            `db:"completed_count"`
	Status          string         `db:"status"`
	UpdatedAt       time.Time      `db:"updated_at"`
	CreatedAt       time.Time      `db:"created_at"`
}

type QueryExpertList struct {
	Limit           int
	Offset          int
	Status          string
	Specializations pq.StringArray
}

type ExpertList struct {
	Items []*Expert `db:"items"`
	Total int       `db:"total"`
}

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
	Status               string    `db:"status"`
	UpdatedAt            time.Time `db:"updated_at"`
	CreatedAt            time.Time `db:"created_at"`
}

type QueryRequisitionList struct {
	Limit           int            `db:"limit"`
	Offset          int            `db:"offset"`
	Status          string         `db:"status"`
	ExpertID        string         `db:"expert_id"`
	FeedbackTime    string         `db:"feedback_time"`
	FeedbackWeekDay string         `db:"feedback_week_day"`
	Specializations pq.StringArray `db:"specializations"`
}

func (q *QueryRequisitionList) SqlxBuildWhere() string {
	query := ""

	if q.Status != "" {
		query = query + ` status=:status AND`
	}

	if len(q.Specializations) != 0 {
		query = query + ` diagnosis=any(:specializations) AND`
	}

	if q.ExpertID != "" {
		query = query + ` expert_id=:expert_id AND`
	}

	if q.FeedbackWeekDay != "" {
		query = query + ` feedback_week_day=:feedback_week_day AND`
	}

	if q.FeedbackTime != "" {
		query = query + ` feedback_time=:feedback_time AND`
	}

	query = strings.TrimSuffix(query, "AND")

	if query == "" {
		return ""
	}

	return " WHERE " + query
}

func (q *QueryRequisitionList) PgxBuildWhereOrder(rawQuery string) (string, []interface{}) {
	rawArgs := map[string]interface{}{}

	if q.Status != "" {
		rawArgs["status"] = q.Status
	}

	if len(q.Specializations) != 0 {
		rawArgs["diagnosis"] = q.Specializations
	}

	if q.ExpertID != "" {
		rawArgs["expert_id"] = q.ExpertID
	}

	if q.FeedbackWeekDay != "" {
		rawArgs["feedback_week_day"] = q.FeedbackWeekDay
	}

	if q.FeedbackTime != "" {
		rawArgs["feedback_time"] = q.FeedbackTime
	}

	var (
		index = 1
		query = rawQuery
		args  = make([]interface{}, 0)
	)

	if len(rawArgs) != 0 {

		query = rawQuery + " WHERE "
		for key, arg := range rawArgs {

			// workaround!!
			if key == "diagnosis" {
				query = query + fmt.Sprintf(" %s=any($%d) AND", key, index)
			} else {
				query = query + fmt.Sprintf(" %s=$%d AND", key, index)
			}
			args = append(args, arg)

			index++
		}

		query = strings.TrimSuffix(query, "AND")
	}
	// ordering
	query = query + fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", index, index+1)
	args = append(args, q.Limit, q.Offset)

	return query, args
}

func (q *QueryRequisitionList) PgxBuildWhere(rawQuery string) (string, []interface{}) {
	rawArgs := map[string]interface{}{}

	if q.Status != "" {
		rawArgs["status"] = q.Status
	}

	if len(q.Specializations) != 0 {
		rawArgs["diagnosis"] = q.Specializations
	}

	if q.ExpertID != "" {
		rawArgs["expert_id"] = q.ExpertID
	}

	if q.FeedbackWeekDay != "" {
		rawArgs["feedback_week_day"] = q.FeedbackWeekDay
	}

	if q.FeedbackTime != "" {
		rawArgs["feedback_time"] = q.FeedbackTime
	}

	if len(rawArgs) == 0 {
		return rawQuery, []interface{}{}
	}

	index := 1
	args := make([]interface{}, 0)
	query := rawQuery + " WHERE "

	for key, arg := range rawArgs {

		// workaround!!
		if key == "diagnosis" {
			query = query + fmt.Sprintf(" %s=any($%d) AND", key, index)
		} else {
			query = query + fmt.Sprintf(" %s=$%d AND", key, index)
		}
		args = append(args, arg)

		index++
	}

	return strings.TrimSuffix(query, "AND"), args
}

type RequisitionList struct {
	Items []*Requisition `db:"items"`
	Total int            `db:"total"`
}

type Admin struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}
