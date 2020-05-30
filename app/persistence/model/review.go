package model

import (
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

type ReviewList struct {
	Items []*Review `db:"items"`
	Total int       `db:"total"`
}
