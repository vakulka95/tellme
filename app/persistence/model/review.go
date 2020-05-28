package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	ReviewStatusCompleted = "completed"
	ReviewStatusRequested = "requested"
)

type Review struct {
	ID                 string    `db:"id"`
	ExpertID           string    `db:"expert_id"`
	ExpertUsername     string    `db:"expert_username"`
	RequisitionID      string    `db:"requisition_id"`
	PlatformReview     string    `db:"platform_review"`
	ConsultationCount  int       `db:"consultation_count"`
	ConsultationReview string    `db:"consultation_review"`
	ExpertPoint        int       `db:"expert_point"`
	ExpertReview       string    `db:"expert_review"`
	Token              string    `db:"token"`
	Status             string    `db:"status"`
	UpdatedAt          time.Time `db:"updated_at"`
	CreatedAt          time.Time `db:"created_at"`
}
type QueryReviewList struct {
	Limit    int
	Offset   int
	Status   string
	ExpertID string
}

func (q *QueryReviewList) BuildWhereOrder(rawQuery string) (string, []interface{}) {
	rawArgs := map[string]interface{}{}

	if q.Status != "" {
		rawArgs["status"] = q.Status
	}

	if q.ExpertID != "" {
		rawArgs["expert_id"] = q.ExpertID
	}

	var (
		index = 1
		query = rawQuery
		args  = make([]interface{}, 0)
	)

	if len(rawArgs) != 0 {
		query = rawQuery + " WHERE "

		for key, arg := range rawArgs {
			query = query + fmt.Sprintf(" %s=$%d AND", key, index)
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

func (q *QueryReviewList) BuildWhere(rawQuery string) (string, []interface{}) {
	rawArgs := map[string]interface{}{}

	if q.Status != "" {
		rawArgs["status"] = q.Status
	}

	if q.ExpertID != "" {
		rawArgs["expert_id"] = q.ExpertID
	}

	if len(rawArgs) == 0 {
		return rawQuery, []interface{}{}
	}

	index := 1
	args := make([]interface{}, 0)
	query := rawQuery + " WHERE "

	for key, arg := range rawArgs {
		query = query + fmt.Sprintf(" %s=$%d AND", key, index)
		args = append(args, arg)

		index++
	}

	return strings.TrimSuffix(query, "AND"), args
}

type ReviewList struct {
	Items []*Review `db:"items"`
	Total int       `db:"total"`
}
