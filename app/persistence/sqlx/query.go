package sqlx

import (
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

//
// Experts
//
func (r *Repository) GetExpert(id string) (*model.Expert, error) {
	expert := &model.Expert{}
	if err := r.statements[model.StatementGetExpert].Get(expert, id); err != nil {
		return nil, err
	}
	return expert, nil
}

func (r *Repository) GetExpertByEmail(email string) (*model.Expert, error) {
	expert := &model.Expert{}
	if err := r.statements[model.StatementGetExpertByEmail].Get(expert, email); err != nil {
		return nil, err
	}
	return expert, nil
}

func (r *Repository) GetExpertByPhone(phone string) (*model.Expert, error) {
	expert := &model.Expert{}
	if err := r.statements[model.StatementGetExpertByPhone].Get(expert, phone); err != nil {
		return nil, err
	}
	return expert, nil
}

func (r *Repository) GetExpertList(q *model.QueryExpertList) (*model.ExpertList, error) {
	var (
		list = &model.ExpertList{Items: make([]*model.Expert, 0)}
	)

	query := `
			SELECT id, username, gender, phone, email, password, specializations,
					education, document_urls, status, updated_at, created_at
			FROM experts 
`

	count := `
			SELECT count(id) FROM experts
`

	switch {
	case q.Status != "" && len(q.Specializations) != 0:

		query = query + ` WHERE status=$3 AND specializations @> $4  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE status=$1 AND specializations @> $2`

		if err := r.cli.Select(&list.Items, query, q.Limit, q.Offset, q.Status, q.Specializations); err != nil {
			return nil, err
		}

		if err := r.cli.Get(&list.Total, count, q.Status, q.Specializations); err != nil {
			return nil, err
		}

		return list, nil

	case q.Status != "" && len(q.Specializations) == 0:

		query = query + ` WHERE status=$3  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE status=$1`

		if err := r.cli.Select(&list.Items, query, q.Limit, q.Offset, q.Status); err != nil {
			return nil, err
		}

		if err := r.cli.Get(&list.Total, count, q.Status); err != nil {
			return nil, err
		}

		return list, nil
	case q.Status == "" && len(q.Specializations) != 0:

		query = query + ` WHERE specializations @> $3  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE specializations @> $1`

		if err := r.cli.Select(&list.Items, query, q.Limit, q.Offset, q.Specializations); err != nil {
			return nil, err
		}

		if err := r.cli.Get(&list.Total, count, q.Specializations); err != nil {
			return nil, err
		}

		return list, nil

	default:

		if err := r.statements[model.StatementGetExpertList].Select(&list.Items, q.Limit, q.Offset); err != nil {
			return nil, err
		}

		if err := r.statements[model.StatementGetExpertCount].Get(&list.Total); err != nil {
			return nil, err
		}

		return list, nil
	}
}

func (r *Repository) CreateExpert(e *model.Expert) (*model.Expert, error) {
	if _, err := r.statements[model.StatementCreateExpert].Exec(
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
	); err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) DeleteExpert(e *model.Expert) error {
	_, err := r.statements[model.StatementDeleteExpert].Exec(e.ID)
	return err
}

func (r *Repository) UpdateExpertStatus(e *model.Expert) (*model.Expert, error) {
	if _, err := r.statements[model.StatementUpdateExpertStatus].
		Exec(e.ID, e.Status); err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertPassword(e *model.Expert) (*model.Expert, error) {
	if _, err := r.statements[model.StatementUpdateExpertPassword].
		Exec(e.ID, e.Password); err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertDocumentURLs(e *model.Expert) (*model.Expert, error) {
	if _, err := r.statements[model.StatementUpdateExpertDocumentURLs].
		Exec(e.ID, e.DocumentURLs); err != nil {
		return nil, err
	}
	return &model.Expert{ID: e.ID}, nil
}

//
// Requisitions
//
func (r *Repository) CreateRequisition(req *model.Requisition) (*model.Requisition, error) {
	if _, err := r.statements[model.StatementCreateRequisition].Exec(
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
	); err != nil {
		return nil, err
	}

	return &model.Requisition{ID: req.ID}, nil
}

func (r *Repository) GetRequisitionList(q *model.QueryRequisitionList) (*model.RequisitionList, error) {
	list := &model.RequisitionList{Items: make([]*model.Requisition, 0)}

	query := `
			SELECT id, expert_id, username, gender, phone, diagnosis, diagnosis_description, 
					expert_gender, feedback_type, feedback_contact, feedback_time, feedback_week_day, is_adult, status, created_at, updated_at
			FROM requisitions `

	count := `
			SELECT count(id) FROM requisitions
`

	where := q.SqlxBuildWhere()
	order := ` ORDER BY created_at DESC LIMIT :limit OFFSET :offset`

	rows, err := r.cli.NamedQuery(query+where+order, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &model.Requisition{}

		if err = rows.StructScan(item); err != nil {
			log.Printf("failed to scan requisition: %v", err)
			continue
		}

		list.Items = append(list.Items, item)
	}

	countRows, err := r.cli.NamedQuery(count+where, q)
	if err != nil {
		return nil, err
	}

	if countRows.Next() {
		if err = countRows.Scan(&list.Total); err != nil {
			return nil, err
		}
	}

	return list, nil
}

func (r *Repository) UpdateRequisitionStatus(q *model.Requisition) (*model.Requisition, error) {
	if _, err := r.statements[model.StatementUpdateRequisitionStatus].
		Exec(q.ID, q.Status, q.ExpertID); err != nil {
		return nil, err
	}
	return &model.Requisition{ID: q.ID}, nil
}

//
// Admins
//
func (r *Repository) GetAdminByLogin(login string) (*model.Admin, error) {
	admin := &model.Admin{}
	if err := r.statements[model.StatementGetAdminByLogin].Get(admin, login); err != nil {
		return nil, err
	}
	return admin, nil
}
