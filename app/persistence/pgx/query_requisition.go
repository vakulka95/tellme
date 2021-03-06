package pgx

import (
	"context"
	"log"
	"time"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"
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
		    session_count,
			status,
			created_at,
			updated_at
	  FROM v$requisitions
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
			&requisition.SessionCount,
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

	query := postgres.NewQueryBuilder().
		Select(
			"id",
			"expert_id",
			"username",
			"gender",
			"phone",
			"diagnosis",
			"diagnosis_description",
			"expert_gender",
			"feedback_type",
			"feedback_contact",
			"feedback_time",
			"feedback_week_day",
			"is_adult",
			"session_count",
			"status",
			"created_at",
			"updated_at",
		).
		From("v$requisitions").
		Where(
			postgres.NewExpression("status", postgres.NewString(q.Status), postgres.OperatorEqual),
			postgres.NewExpression("expert_id", postgres.NewString(q.ExpertID), postgres.OperatorEqual),
			//postgres.NewExpression("feedback_time", postgres.NewString(q.FeedbackTime), postgres.OperatorEqual),
			//postgres.NewExpression("feedback_week_day", postgres.NewString(q.FeedbackWeekDay), postgres.OperatorEqual),
			postgres.NewExpression("diagnosis", postgres.NewSliceString(q.Specializations), postgres.OperatorAny),
			postgres.NewExpression("created_at", postgres.NewTimestamp(q.CreatedAtFrom), postgres.OperatorGreaterEqualOrEqual),
			postgres.NewExpression("created_at", postgres.NewTimestamp(q.CreatedAtTo), postgres.OperatorLessThenOrEqual),
			model.DiscoverRequisitionExpression(q.Search),
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
			&item.SessionCount,
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
	const query = `UPDATE requisitions SET status=$2, expert_id=$3, updated_at=now() WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, q.ID, q.Status, q.ExpertID)
	if err != nil {
		return nil, err
	}

	return &model.Requisition{ID: q.ID}, nil
}

func (r *Repository) UpdateRequisitionSMSReplyCount(q *model.Requisition) (*model.Requisition, error) {
	const query = `UPDATE requisitions SET sms_reply_count=$2, updated_at=now() WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, q.ID, q.SMSReplyCount)
	if err != nil {
		return nil, err
	}

	return &model.Requisition{ID: q.ID}, nil
}

func (r *Repository) GetNotReviewedRequisition() ([]*model.Requisition, error) {
	var (
		ctx  = context.TODO()
		list = make([]*model.Requisition, 0, 20) // pseudo default capacity
	)

	const query = `
	SELECT id,
		   expert_id,
		   phone
	  FROM v$requisition_not_reviewed
     WHERE review_id is null 
	   AND status='completed'
`

	rows, err := r.cli.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Requisition{}
		if err = rows.Scan(
			&item.ID,
			&item.ExpertID,
			&item.Phone,
		); err != nil {
			log.Printf("failed to scan requisition for not reviewed status: %v", err)
			continue
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *Repository) GetNotProcessedRequisition(lifetimeTo *time.Time) ([]*model.Requisition, error) {
	var (
		ctx  = context.TODO()
		list = make([]*model.Requisition, 0, 20) // pseudo default capacity
	)

	query := postgres.NewQueryBuilder().
		Select(
			"id",
			"username",
			"phone",
			"gender",
			"status",
			"sms_reply_count",
			"created_at",
			"updated_at",
		).
		From("v$requisitions").
		Where(
			postgres.NewExpression("status", postgres.NewString(model.RequisitionStatusCreated), postgres.OperatorEqual),
			postgres.NewExpression("sms_reply_count", postgres.NewInt(2), postgres.OperatorLessThen),
			postgres.NewExpression("created_at", postgres.NewTimestamp(lifetimeTo), postgres.OperatorLessThenOrEqual)).
		OrderBy("created_at").
		OrderDir(postgres.OrderDirDESC)

	listQuery, listArgs := query.Build()

	rows, err := r.cli.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Requisition{}
		if err = rows.Scan(
			&item.ID,
			&item.Username,
			&item.Phone,
			&item.Gender,
			&item.Status,
			&item.SMSReplyCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			log.Printf("failed to scan requisition: %v", err)
			continue
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *Repository) DeleteRequisition(v *model.Requisition) error {
	const query = `DELETE FROM requisitions WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, v.ID)
	return err
}
