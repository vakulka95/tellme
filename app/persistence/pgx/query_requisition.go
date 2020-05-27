package pgx

import (
	"context"
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

func (r *Repository) GetRequisition(id string) (*model.Requisition, error) {
	const query = `
	SELECT 	id,
			expert_id,
			username,
			gender,
			phone,
			diagnosis,
			diagnosis_description,
			expert_gender,
			feedback_type,
			feedback_contact,
			feedback_time,
			feedback_week_day,
			is_adult,
			status,
			created_at,
			updated_at
	  FROM requisitions
	 WHERE id=$1
`

	var (
		ctx         = context.TODO()
		requisition = &model.Requisition{}
	)

	err := r.cli.QueryRow(ctx, query, id).
		Scan(
			&requisition.ID,
			&requisition.ExpertID,
			&requisition.Username,
			&requisition.Gender,
			&requisition.Phone,
			&requisition.Diagnosis,
			&requisition.DiagnosisDescription,
			&requisition.ExpertGender,
			&requisition.FeedbackType,
			&requisition.FeedbackContact,
			&requisition.FeedbackTime,
			&requisition.FeedbackWeekDay,
			&requisition.IsAdult,
			&requisition.Status,
			&requisition.CreatedAt,
			&requisition.UpdatedAt,
		)
	if err != nil {
		return nil, err
	}

	return requisition, nil
}

func (r *Repository) CreateRequisition(req *model.Requisition) (*model.Requisition, error) {
	const query = `
	INSERT INTO requisitions (
		id,
		expert_id,
		username,
		gender,
		phone,
		diagnosis,
		diagnosis_description,
		expert_gender,
		feedback_type,
		feedback_contact,
		feedback_time,
		feedback_week_day,
		is_adult,
		status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
		req.ID,
		req.ExpertID,
		req.Username,
		req.Gender,
		req.Phone,
		req.Diagnosis,
		req.DiagnosisDescription,
		req.ExpertGender,
		req.FeedbackType,
		req.FeedbackContact,
		req.FeedbackTime,
		req.FeedbackWeekDay,
		req.IsAdult,
		req.Status,
	)
	if err != nil {
		return nil, err
	}

	return &model.Requisition{ID: req.ID}, nil
}

func (r *Repository) GetRequisitionList(q *model.QueryRequisitionList) (*model.RequisitionList, error) {
	var (
		ctx  = context.TODO()
		list = &model.RequisitionList{Items: make([]*model.Requisition, 0, 20)} // pseudo default capacity
	)

	rawListQuery := `
			SELECT 	id,
					expert_id,
					username,
					gender,
					phone,
					diagnosis,
					diagnosis_description,
					expert_gender,
					feedback_type,
					feedback_contact,
					feedback_time,
					feedback_week_day,
					is_adult,
					status,
					created_at,
					updated_at
			FROM requisitions
 `

	rawCountQuery := `
			SELECT count(id) FROM requisitions
 `

	listQuery, listArgs := q.BuildWhereOrder(rawListQuery)
	countQuery, countArgs := q.BuildWhere(rawCountQuery)

	rows, err := r.cli.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Requisition{}
		if err = rows.Scan(
			&item.ID,
			&item.ExpertID,
			&item.Username,
			&item.Gender,
			&item.Phone,
			&item.Diagnosis,
			&item.DiagnosisDescription,
			&item.ExpertGender,
			&item.FeedbackType,
			&item.FeedbackContact,
			&item.FeedbackTime,
			&item.FeedbackWeekDay,
			&item.IsAdult,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("failed to scan requisition: %v", err)
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

func (r *Repository) UpdateRequisitionStatus(q *model.Requisition) (*model.Requisition, error) {
	const query = `UPDATE requisitions SET status=$2, expert_id=$3 WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, q.ID, q.Status, q.ExpertID)
	if err != nil {
		return nil, err
	}

	return &model.Requisition{ID: q.ID}, nil
}
