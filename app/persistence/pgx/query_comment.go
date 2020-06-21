package pgx

import (
	"context"
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"
)

func (r *Repository) GetComment(id string) (*model.Comment, error) {
	const query = `
	SELECT 	id,
			admin_id,
			expert_id,
			admin_username,
			body,
			updated_at,
			created_at
	  FROM v$comments
	 WHERE id=$1
`

	var (
		ctx    = context.TODO()
		review = &model.Comment{}
	)

	err := r.cli.QueryRow(ctx, query, id).
		Scan(
			&review.ID,
			&review.AdminID,
			&review.ExpertID,
			&review.AdminUsername,
			&review.Body,
			&review.UpdatedAt,
			&review.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Repository) GetExpertCommentList(expertID string) ([]*model.Comment, error) {
	var (
		ctx  = context.TODO()
		list = make([]*model.Comment, 0, 20) // pseudo default capacity
	)

	query := postgres.NewQueryBuilder().
		Select(
			"id",
			"admin_id",
			"expert_id",
			"admin_username",
			"body",
			"updated_at",
			"created_at",
		).
		From("v$comments").
		Where(
			postgres.NewExpression("expert_id", postgres.NewString(expertID), postgres.OperatorEqual),
		).
		OrderBy("created_at").
		OrderDir(postgres.OrderDirDESC)

	listQuery, listArgs := query.Build()

	rows, err := r.cli.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Comment{}
		if err = rows.Scan(
			&item.ID,
			&item.AdminID,
			&item.ExpertID,
			&item.AdminUsername,
			&item.Body,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("failed to scan comment: %v", err)
			continue
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *Repository) GetCommentList(q *model.QueryCommentList) (*model.CommentList, error) {
	var (
		ctx  = context.TODO()
		list = &model.CommentList{Items: make([]*model.Comment, 0, 20)} // pseudo default capacity
	)

	query := postgres.NewQueryBuilder().
		Select(
			"id",
			"admin_id",
			"expert_id",
			"admin_username",
			"body",
			"updated_at",
			"created_at",
		).
		From("v$comments").
		Where(
			postgres.NewExpression("expert_id", postgres.NewString(q.ExpertID), postgres.OperatorEqual),
		).
		OrderBy("created_at").
		OrderDir(postgres.OrderDirDESC).
		Limit(q.Limit).
		Offset(q.Offset)

	listQuery, listArgs := query.Build()
	countQuery, countArgs := query.BuildCount()

	rows, err := r.cli.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Comment{}
		if err = rows.Scan(
			&item.ID,
			&item.AdminID,
			&item.ExpertID,
			&item.AdminUsername,
			&item.Body,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("failed to scan comment: %v", err)
			continue
		}
		list.Items = append(list.Items, item)
	}

	err = r.cli.QueryRow(ctx, countQuery, countArgs...).Scan(&list.Total)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) CreateComment(c *model.Comment) (*model.Comment, error) {
	const query = `
	INSERT INTO comments (
			id,
			admin_id,
			expert_id,
			body
	) VALUES ($1, $2, $3, $4)`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
		c.ID,
		c.AdminID,
		c.ExpertID,
		c.Body,
	)
	if err != nil {
		return nil, err
	}

	return &model.Comment{ID: c.ID}, nil
}

func (r *Repository) DeleteComment(c *model.Comment) error {
	const query = `DELETE FROM comments WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, c.ID)
	return err
}
