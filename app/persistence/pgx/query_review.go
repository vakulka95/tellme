package pgx

import (
	"context"
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

func (r *Repository) GetReview(id string) (*model.Review, error) {
	const query = `
	SELECT 	id,
			expert_id,
			requisition_id,
			platform_review,
			consultation_count,
			consultation_review,
			expert_point,
			expert_review,
			token,
			status,
			updated_at,
			created_at
	  FROM reviews
	 WHERE id=$1
`

	var (
		ctx    = context.TODO()
		review = &model.Review{}
	)

	err := r.cli.QueryRow(ctx, query, id).
		Scan(
			&review.ID,
			&review.ExpertID,
			&review.RequisitionID,
			&review.PlatformReview,
			&review.ConsultationCount,
			&review.ConsultationReview,
			&review.ExpertPoint,
			&review.ExpertReview,
			&review.Token,
			&review.Status,
			&review.UpdatedAt,
			&review.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Repository) GetReviewByToken(token string) (*model.Review, error) {
	const query = `
	SELECT 	id,
			expert_id,
			requisition_id,
			platform_review,
			consultation_count,
			consultation_review,
			expert_point,
			expert_review,
			token,
			status,
			updated_at,
			created_at
	  FROM reviews
	 WHERE token=$1
`

	var (
		ctx    = context.TODO()
		review = &model.Review{}
	)

	err := r.cli.QueryRow(ctx, query, token).
		Scan(
			&review.ID,
			&review.ExpertID,
			&review.RequisitionID,
			&review.PlatformReview,
			&review.ConsultationCount,
			&review.ConsultationReview,
			&review.ExpertPoint,
			&review.ExpertReview,
			&review.Token,
			&review.Status,
			&review.UpdatedAt,
			&review.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Repository) GetReviewList(q *model.QueryReviewList) (*model.ReviewList, error) {
	var (
		ctx  = context.TODO()
		list = &model.ReviewList{Items: make([]*model.Review, 0, 20)} // pseudo default capacity
	)

	rawListQuery := `
	SELECT 	id,
			expert_id,
			expert_username,
			requisition_id,
			platform_review,
			consultation_count,
			consultation_review,
			expert_point,
			expert_review,
			token,
			status,
			updated_at,
			created_at
	  FROM v$reviews
`

	rawCountQuery := `
			SELECT count(id) FROM v$reviews
`

	listQuery, listArgs := q.BuildWhereOrder(rawListQuery)
	countQuery, countArgs := q.BuildWhere(rawCountQuery)

	rows, err := r.cli.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Review{}
		if err = rows.Scan(
			&item.ID,
			&item.ExpertID,
			&item.ExpertUsername,
			&item.RequisitionID,
			&item.PlatformReview,
			&item.ConsultationCount,
			&item.ConsultationReview,
			&item.ExpertPoint,
			&item.ExpertReview,
			&item.Token,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("failed to scan review: %v", err)
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

func (r *Repository) CreateReview(v *model.Review) (*model.Review, error) {
	const query = `
	INSERT INTO reviews (
			expert_id,
			requisition_id,
			platform_review,
			consultation_count,
			consultation_review,
			expert_point,
			expert_review,
			token,
			status
	) VALUES ($1, $2, $3, $4, $5, $6)`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
		v.ID,
		v.ExpertID,
		v.RequisitionID,
		v.PlatformReview,
		v.ConsultationCount,
		v.ConsultationReview,
		v.ExpertPoint,
		v.ExpertReview,
		v.Token,
		v.Status,
	)
	if err != nil {
		return nil, err
	}

	return &model.Review{ID: v.ID}, nil
}

func (r *Repository) UpdateReviewBodyStatus(v *model.Review) (*model.Review, error) {
	const query = `
	UPDATE reviews 
	   SET platform_review=$2,
           consultation_count=$3,
           consultation_review=$4,
           expert_point=$5,
           expert_review=$6,
		   status=$7
     WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
		v.ID,
		v.PlatformReview,
		v.ConsultationCount,
		v.ConsultationReview,
		v.ExpertPoint,
		v.ExpertReview,
		v.Status,
	)
	if err != nil {
		return nil, err
	}

	return &model.Review{ID: v.ID}, nil
}
