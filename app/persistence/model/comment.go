package model

import "time"

type Comment struct {
	ID            string    `db:"id"`
	AdminID       string    `db:"admin_id"`
	AdminUsername string    `db:"admin_usernsame"`
	ExpertID      string    `db:"expert_id"`
	Body          string    `db:"body"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedAt     time.Time `db:"created_at"`
}
type QueryCommentList struct {
	Limit    int
	Offset   int
	Status   string
	ExpertID string
}

type CommentList struct {
	Items []*Comment `db:"items"`
	Total int        `db:"total"`
}
