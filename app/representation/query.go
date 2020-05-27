package representation

import (
	"bytes"
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/paginater"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

const (
	defaultQueryLimit  = 10
	defaultQueryOffset = 0
)

type QueryListParams struct {
	Limit           int      `form:"limit"`
	Offset          int      `form:"offset"`
	Status          string   `form:"status"`
	FeedbackTime    string   `form:"feedback_time"`
	FeedbackWeekDay string   `form:"feedback_week_day"`
	Specializations []string `form:"specializations"`
	ExpertID        string   `form:"expert_id"`
}

func QueryExpertAPItoPersistence(q *QueryListParams) *model.QueryExpertList {
	if q.Limit == 0 || q.Limit > 50 {
		q.Limit = defaultQueryLimit
	}

	return &model.QueryExpertList{
		Limit:           q.Limit,
		Offset:          q.Offset,
		Status:          q.Status,
		Specializations: q.Specializations,
	}
}

func QueryRequisitionAPItoPersistence(q *QueryListParams) *model.QueryRequisitionList {
	if q.Limit == 0 || q.Limit > 50 {
		q.Limit = defaultQueryLimit
	}

	return &model.QueryRequisitionList{
		Limit:           q.Limit,
		Offset:          q.Offset,
		Status:          q.Status,
		FeedbackTime:    q.FeedbackTime,
		FeedbackWeekDay: q.FeedbackWeekDay,
		Specializations: q.Specializations,
		ExpertID:        q.ExpertID,
	}
}

func QueryReviewAPItoPersistence(q *QueryListParams) *model.QueryReviewList {
	if q.Limit == 0 || q.Limit > 50 {
		q.Limit = defaultQueryLimit
	}

	return &model.QueryReviewList{
		Limit:    q.Limit,
		Offset:   q.Offset,
		Status:   q.Status,
		ExpertID: q.ExpertID,
	}
}

func (q *QueryListParams) generateQuery() string {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("?status=%s", q.Status))

	for _, v := range q.Specializations {
		buf.WriteString(fmt.Sprintf("&specializations=%s", v))
	}

	buf.WriteString(fmt.Sprintf("&feedback_week_day=%s", q.FeedbackWeekDay))
	buf.WriteString(fmt.Sprintf("&feedback_time=%s", q.FeedbackTime))
	buf.WriteString(fmt.Sprintf("&expert_id=%s", q.ExpertID))

	return buf.String()
}

func (q *QueryListParams) GeneratePagination(total int) gin.H {
	var (
		currentPage = math.Round(float64(q.Offset/q.Limit)) + 1
		paginator   = paginater.New(total, q.Limit, int(currentPage), 1)
	)

	query := q.generateQuery()
	return gin.H{
		"firstQuery":   fmt.Sprintf("%s&limit=%d&offset=%d", query, q.Limit, 0),
		"nextQuery":    fmt.Sprintf("%s&limit=%d&offset=%d", query, q.Limit, q.Limit*paginator.Current()),
		"currentQuery": fmt.Sprintf("%s&limit=%d", query, q.Limit),
		"prevQuery":    fmt.Sprintf("%s&limit=%d&offset=%d", query, q.Limit, q.Limit*(paginator.Previous()-1)),
		"lastQuery":    fmt.Sprintf("%s&limit=%d&offset=%d", query, q.Limit, q.Limit*(paginator.TotalPages()-1)),
		"page":         paginator,
	}
}
