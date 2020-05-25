package pgx

import (
	"context"
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

//
// Experts
//
func (r *Repository) GetExpert(id string) (*model.Expert, error) {
	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, r.queries[model.StatementGetExpert], id).
		// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at
		Scan(&expert.ID, &expert.Username, &expert.Gender, &expert.Phone, &expert.Email, &expert.Password, &expert.Specializations,
			&expert.Education, &expert.DocumentURLs, &expert.Status, &expert.UpdatedAt, &expert.CreatedAt)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertByEmail(email string) (*model.Expert, error) {
	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, r.queries[model.StatementGetExpertByEmail], email).
		// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at
		Scan(&expert.ID, &expert.Username, &expert.Gender, &expert.Phone, &expert.Email, &expert.Password, &expert.Specializations,
			&expert.Education, &expert.DocumentURLs, &expert.Status, &expert.UpdatedAt, &expert.CreatedAt)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertByPhone(phone string) (*model.Expert, error) {
	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, r.queries[model.StatementGetExpertByPhone], phone).
		// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at
		Scan(&expert.ID, &expert.Username, &expert.Gender, &expert.Phone, &expert.Email, &expert.Password, &expert.Specializations,
			&expert.Education, &expert.DocumentURLs, &expert.Status, &expert.UpdatedAt, &expert.CreatedAt)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertList(q *model.QueryExpertList) (*model.ExpertList, error) {
	var (
		ctx  = context.TODO()
		list = &model.ExpertList{Items: make([]*model.Expert, 0, 20)} // pseudo default capacity
	)

	query := `
			SELECT id, username, gender, phone, email, password, specializations,
					education, document_urls, status, updated_at, created_at, processing_count, completed_count
			FROM v$experts
`

	count := `
			SELECT count(id) FROM experts
`

	switch {
	case q.Status != "" && len(q.Specializations) != 0:

		query = query + ` WHERE status=$3 AND specializations @> $4  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE status=$1 AND specializations @> $2`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset, q.Status, q.Specializations)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			item := &model.Expert{}
			// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at, processing_count, completed_count
			if err := rows.Scan(&item.ID, &item.Username, &item.Gender, &item.Phone, &item.Email, &item.Password, &item.Specializations,
				&item.Education, &item.DocumentURLs, &item.Status, &item.UpdatedAt, &item.CreatedAt, &item.ProcessingCount, &item.CompletedCount); err != nil {
				log.Printf("failed to scan expert: %v", err)
				continue
			}
			list.Items = append(list.Items, item)
		}

		if err := r.cli.QueryRow(ctx, count, q.Status, q.Specializations).Scan(&list.Total); err != nil {
			return nil, err
		}

		return list, rows.Err()

	case q.Status != "" && len(q.Specializations) == 0:

		query = query + ` WHERE status=$3  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE status=$1`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset, q.Status)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			item := &model.Expert{}
			// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at, processing_count, completed_count
			if err := rows.Scan(&item.ID, &item.Username, &item.Gender, &item.Phone, &item.Email, &item.Password, &item.Specializations,
				&item.Education, &item.DocumentURLs, &item.Status, &item.UpdatedAt, &item.CreatedAt, &item.ProcessingCount, &item.CompletedCount); err != nil {
				log.Printf("failed to scan expert: %v", err)
				continue
			}
			list.Items = append(list.Items, item)
		}

		if err := r.cli.QueryRow(ctx, count, q.Status).Scan(&list.Total); err != nil {
			return nil, err
		}

		return list, rows.Err()
	case q.Status == "" && len(q.Specializations) != 0:

		query = query + ` WHERE specializations @> $3  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE specializations @> $1`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset, q.Specializations)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			item := &model.Expert{}
			// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at, processing_count, completed_count
			if err := rows.Scan(&item.ID, &item.Username, &item.Gender, &item.Phone, &item.Email, &item.Password, &item.Specializations,
				&item.Education, &item.DocumentURLs, &item.Status, &item.UpdatedAt, &item.CreatedAt, &item.ProcessingCount, &item.CompletedCount); err != nil {
				log.Printf("failed to scan expert: %v", err)
				continue
			}
			list.Items = append(list.Items, item)
		}

		if err := r.cli.QueryRow(ctx, count, q.Specializations).Scan(&list.Total); err != nil {
			return nil, err
		}

		return list, rows.Err()

	default:

		rows, err := r.cli.Query(ctx, r.queries[model.StatementGetExpertList], q.Limit, q.Offset)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			item := &model.Expert{}
			// id, username, gender, phone, email, password, specializations, education, document_urls, status, updated_at, created_at, processing_count, completed_count
			if err := rows.Scan(&item.ID, &item.Username, &item.Gender, &item.Phone, &item.Email, &item.Password, &item.Specializations,
				&item.Education, &item.DocumentURLs, &item.Status, &item.UpdatedAt, &item.CreatedAt, &item.ProcessingCount, &item.CompletedCount); err != nil {
				log.Printf("failed to scan expert: %v", err)
				continue
			}
			list.Items = append(list.Items, item)
		}

		if err := r.cli.QueryRow(ctx, r.queries[model.StatementGetExpertCount]).Scan(&list.Total); err != nil {
			return nil, err
		}

		return list, rows.Err()
	}
}

func (r *Repository) CreateExpert(e *model.Expert) (*model.Expert, error) {
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementCreateExpert],
		e.ID,
		e.Username,
		e.Gender,
		e.Phone,
		e.Email,
		e.Password,
		e.Specializations,
		e.Education,
		e.DocumentURLs,
		e.Status,
	)
	if err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) DeleteExpert(e *model.Expert) error {
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementDeleteExpert], e.ID)
	return err
}

func (r *Repository) UpdateExpertStatus(e *model.Expert) (*model.Expert, error) {
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementUpdateExpertStatus], e.ID, e.Status)
	if err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertPassword(e *model.Expert) (*model.Expert, error) {
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementUpdateExpertPassword], e.ID, e.Password)
	if err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertDocumentURLs(e *model.Expert) (*model.Expert, error) {
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementUpdateExpertDocumentURLs], e.ID, e.DocumentURLs)
	if err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

//
// Requisitions
//
func (r *Repository) CreateRequisition(req *model.Requisition) (*model.Requisition, error) {
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementCreateRequisition],
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
			SELECT id, expert_id, username, gender, phone, diagnosis, diagnosis_description,
					expert_gender, feedback_type, feedback_contact, feedback_time, feedback_week_day, is_adult, status, created_at, updated_at
			FROM requisitions `

	rawCountQuery := `
			SELECT count(id) FROM requisitions
`

	listQuery, listArgs := q.PgxBuildWhereOrder(rawListQuery)
	countQuery, countArgs := q.PgxBuildWhere(rawCountQuery)

	log.Printf("(DEBUG) GetRequisitionList: list query: %s", listQuery)

	rows, err := r.cli.Query(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Requisition{}
		// id, expert_id, username, gender, phone, diagnosis, diagnosis_description, expert_gender, feedback_type, feedback_contact, feedback_time, feedback_week_day, is_adult, status, created_at, updated_at
		if err = rows.Scan(&item.ID, &item.ExpertID, &item.Username, &item.Gender, &item.Phone, &item.Diagnosis, &item.DiagnosisDescription,
			&item.ExpertGender, &item.FeedbackType, &item.FeedbackContact, &item.FeedbackTime, &item.FeedbackWeekDay, &item.IsAdult, &item.Status,
			&item.CreatedAt, &item.UpdatedAt); err != nil {
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
	_, err := r.cli.Exec(context.TODO(), r.queries[model.StatementUpdateRequisitionStatus], q.ID, q.Status, q.ExpertID)
	if err != nil {
		return nil, err
	}
	return &model.Requisition{ID: q.ID}, nil
}

//
// Admins
//
func (r *Repository) GetAdminByLogin(login string) (*model.Admin, error) {
	admin := &model.Admin{}
	err := r.cli.QueryRow(context.TODO(), r.queries[model.StatementGetAdminByLogin], login).
		// id, username, password, status, updated_at, created_at
		Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Status, &admin.UpdatedAt, &admin.CreatedAt)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
