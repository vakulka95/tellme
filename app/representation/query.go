package representation

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unknwon/paginater"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

const (
	defaultQueryLimit  = 10
	defaultQueryOffset = 0
)

type QueryListParams struct {
	Limit           int        `form:"limit"`
	Offset          int        `form:"offset"`
	OrderBy         string     `form:"order_by"`
	OrderDir        string     `form:"order_dir"`
	Status          string     `form:"status"`
	FeedbackTime    string     `form:"feedback_time"`
	FeedbackWeekDay string     `form:"feedback_week_day"`
	Specializations []string   `form:"specializations"`
	ExpertID        string     `form:"expert_id"`
	Search          string     `form:"search"`
	CreatedAtFrom   *time.Time `form:"datetime_from" time_format:"2006-01-02T15:04"`
	CreatedAtTo     *time.Time `form:"datetime_to" time_format:"2006-01-02T15:04"`
}

func (q *QueryListParams) GetDatetimeRange() (from string, to string) {
	if q.CreatedAtFrom != nil && !q.CreatedAtFrom.IsZero() {
		from = q.CreatedAtFrom.Format("2006-01-02T15:04")
	}
	if q.CreatedAtTo != nil && !q.CreatedAtTo.IsZero() {
		to = q.CreatedAtTo.Format("2006-01-02T15:04")
	}
	return
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
		Search:          q.Search,
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
		Search:          q.Search,
		CreatedAtFrom:   q.CreatedAtFrom,
		CreatedAtTo:     q.CreatedAtTo,
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

func QueryExpertRatingAPItoPersistence(q *QueryListParams) *model.QueryExpertRatingList {
	if q.Limit == 0 || q.Limit > 50 {
		q.Limit = defaultQueryLimit
	}

	return &model.QueryExpertRatingList{
		Limit:    q.Limit,
		Offset:   q.Offset,
		OrderBy:  q.OrderBy,
		OrderDir: q.OrderDir,
		Status:   q.Status,
	}
}

func (q *QueryListParams) generateQuery() string {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("?status=%s", q.Status))

	for _, v := range q.Specializations {
		buf.WriteString(fmt.Sprintf("&specializations=%s", v))
	}

	from, to := q.GetDatetimeRange()
	buf.WriteString(fmt.Sprintf("&datetime_from=%s&datetime_to=%s", from, to))
	buf.WriteString(fmt.Sprintf("&feedback_week_day=%s", q.FeedbackWeekDay))
	buf.WriteString(fmt.Sprintf("&feedback_time=%s", q.FeedbackTime))
	buf.WriteString(fmt.Sprintf("&expert_id=%s", q.ExpertID))
	buf.WriteString(fmt.Sprintf("&search=%s", q.Search))
	buf.WriteString(fmt.Sprintf("&order_by=%s", q.OrderBy))
	buf.WriteString(fmt.Sprintf("&order_dir=%s", strings.ToLower(q.OrderDir)))

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
