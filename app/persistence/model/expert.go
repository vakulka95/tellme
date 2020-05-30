package model

import (
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"
)

const (
	ExpertStatusActive   = "active"
	ExpertStatusBlocked  = "blocked"
	ExpertStatusOnReview = "on_review"
)

type Expert struct {
	ID              string    `db:"id"`
	Username        string    `db:"username"`
	Gender          string    `db:"gender"`
	Phone           string    `db:"phone"`
	Email           string    `db:"email"`
	Password        string    `db:"password"`
	Specializations []string  `db:"specializations"`
	Education       string    `db:"education"`
	DocumentURLs    []string  `db:"document_urls"`
	ProcessingCount int       `db:"processing_count"`
	CompletedCount  int       `db:"completed_count"`
	ReviewCount     int       `db:"review_count"`
	Status          string    `db:"status"`
	UpdatedAt       time.Time `db:"updated_at"`
	CreatedAt       time.Time `db:"created_at"`
}

type QueryExpertList struct {
	Limit           int
	Offset          int
	Status          string
	Specializations []string
	Search          string
}

func DiscoverExpertExpression(search string) postgres.Expression {

	if strings.Contains(search, "@") {
		return postgres.NewExpression("email", postgres.NewString(strings.ToLower(search)), postgres.OperatorEqual)
	}

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

type ExpertList struct {
	Items []*Expert `db:"items"`
	Total int       `db:"total"`
}
