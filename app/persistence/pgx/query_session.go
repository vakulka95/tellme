package pgx

import (
	"context"

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
