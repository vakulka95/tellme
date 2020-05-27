package model

import (
	"time"
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
}

type ExpertList struct {
	Items []*Expert `db:"items"`
	Total int       `db:"total"`
}
