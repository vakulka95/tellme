package pgx

import (
	"context"
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

func (r *Repository) GetSession(id string) (*model.Session, error) {
	const query = `
	SELECT 	id,
			expert_id,
			expert_username,
			requisition_id,
			requisition_username,
			updated_at,
			created_at
	  FROM v$sessions
	 WHERE id=$1
`

	var (
		ctx    = context.TODO()
		review = &model.Session{}
	)

	err := r.cli.QueryRow(ctx, query, id).
		Scan(
			&review.ID,
			&review.ExpertID,
			&review.ExpertUsername,
			&review.RequisitionID,
			&review.RequisitionUsername,
			&review.UpdatedAt,
			&review.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Repository) GetRequisitionSessionList(requisitionID string) ([]*model.Session, error) {
	const query = `
	SELECT 	id,
			expert_id,
			requisition_id,
			updated_at,
			created_at
	  FROM sessions
	 WHERE requisition_id=$1
`

	var (
		ctx  = context.TODO()
		list = make([]*model.Session, 0, 3) // pseudo default capacity
	)

	rows, err := r.cli.Query(ctx, query, requisitionID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Session{}
		if err = rows.Scan(
			&item.ID,
			&item.ExpertID,
			&item.RequisitionID,
			&item.UpdatedAt,
			&item.CreatedAt,
		); err != nil {
			log.Printf("failed to scan session: %v", err)
			continue
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *Repository) CreateSession(v *model.Session) (*model.Session, error) {
	const query = `
	INSERT INTO sessions (
			id,
			expert_id,
			requisition_id
	) VALUES ($1, $2, $3)`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
		v.ID,
		v.ExpertID,
		v.RequisitionID,
	)
	if err != nil {
		return nil, err
	}

	return &model.Session{ID: v.ID}, nil
}

func (r *Repository) DeleteRequisitionSessions(requisitionID string) error {
	const query = `DELETE FROM sessions WHERE requisition_id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, requisitionID)
	return err
}

func (r *Repository) DeleteLastRequisitionSession(requisitionID string) error {
	const query = `DELETE FROM sessions WHERE id IN ( SELECT id FROM sessions WHERE requisition_id=$1 ORDER BY id desc LIMIT 1 )`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, requisitionID)
	return err
}
