package model

import (
	"time"
)

type Session struct {
	ID                  string    `db:"id"`
	ExpertID            string    `db:"expert_id"`
	ExpertUsername      string    `db:"expert_username"`
	RequisitionID       string    `db:"requisition_id"`
	RequisitionUsername string    `db:"requisition_username"`
	UpdatedAt           time.Time `db:"updated_at"`
	CreatedAt           time.Time `db:"created_at"`
}
