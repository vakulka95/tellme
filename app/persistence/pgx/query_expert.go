package pgx

import (
	"context"
	"log"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"
)

func (r *Repository) GetExpert(id string) (*model.Expert, error) {
	const query = `
	SELECT 	id,
			username,
			gender,
			phone,
			email,
			password,
			specializations,
			education,
			document_urls,
			status,
			updated_at,
			created_at
	  FROM experts
	 WHERE id=$1
`

	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, query, id).
		Scan(
			&expert.ID,
			&expert.Username,
			&expert.Gender,
			&expert.Phone,
			&expert.Email,
			&expert.Password,
			&expert.Specializations,
			&expert.Education,
			&expert.DocumentURLs,
			&expert.Status,
			&expert.UpdatedAt,
			&expert.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertView(id string) (*model.Expert, error) {
	const query = `
	SELECT 	id,
			username,
			gender,
			phone,
			email,
			password,
			specializations,
			education,
			document_urls,
			status,
			updated_at,
			created_at,
			processing_count,
			completed_count,
			review_count,
			session_count
	  FROM v$experts
	 WHERE id=$1
`

	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, query, id).
		Scan(
			&expert.ID,
			&expert.Username,
			&expert.Gender,
			&expert.Phone,
			&expert.Email,
			&expert.Password,
			&expert.Specializations,
			&expert.Education,
			&expert.DocumentURLs,
			&expert.Status,
			&expert.UpdatedAt,
			&expert.CreatedAt,
			&expert.ProcessingCount,
			&expert.CompletedCount,
			&expert.ReviewCount,
			&expert.SessionCount,
		)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertByEmail(email string) (*model.Expert, error) {
	const query = `
	SELECT 	id,
			username,
			gender,
			phone,
			email,
			password,
			specializations,
			education,
			document_urls,
			status,
			updated_at,
			created_at
	  FROM experts
	 WHERE email=$1
`

	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, query, email).
		Scan(
			&expert.ID,
			&expert.Username,
			&expert.Gender,
			&expert.Phone,
			&expert.Email,
			&expert.Password,
			&expert.Specializations,
			&expert.Education,
			&expert.DocumentURLs,
			&expert.Status,
			&expert.UpdatedAt,
			&expert.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertByPhone(phone string) (*model.Expert, error) {
	const query = `
	SELECT 	id,
			username,
			gender,
			phone,
			email,
			password,
			specializations,
			education,
			document_urls,
			status,
			updated_at,
			created_at
	  FROM experts
	 WHERE phone=$1
`

	var (
		ctx    = context.TODO()
		expert = &model.Expert{}
	)

	err := r.cli.QueryRow(ctx, query, phone).
		Scan(
			&expert.ID,
			&expert.Username,
			&expert.Gender,
			&expert.Phone,
			&expert.Email,
			&expert.Password,
			&expert.Specializations,
			&expert.Education,
			&expert.DocumentURLs,
			&expert.Status,
			&expert.UpdatedAt,
			&expert.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return expert, nil
}

func (r *Repository) GetExpertList(q *model.QueryExpertList) (*model.ExpertList, error) {
	var (
		ctx  = context.TODO()
		list = &model.ExpertList{Items: make([]*model.Expert, 0, q.Limit)} // pseudo default capacity
	)

	query := postgres.NewQueryBuilder().
		Select(
			"id",
			"username",
			"gender",
			"phone",
			"email",
			"password",
			"specializations",
			"education",
			"document_urls",
			"status",
			"updated_at",
			"created_at",
			"processing_count",
			"completed_count",
			"review_count",
			"session_count",
		).
		From("v$experts").
		Where(
			postgres.NewExpression("status", postgres.NewString(q.Status), postgres.OperatorEqual),
			postgres.NewExpression("specializations", postgres.NewSliceString(q.Specializations), postgres.OperatorContains),
			model.DiscoverExpertExpression(q.Search),
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
		item := &model.Expert{}
		if err = rows.Scan(
			&item.ID,
			&item.Username,
			&item.Gender,
			&item.Phone,
			&item.Email,
			&item.Password,
			&item.Specializations,
			&item.Education,
			&item.DocumentURLs,
			&item.Status,
			&item.UpdatedAt,
			&item.CreatedAt,
			&item.ProcessingCount,
			&item.CompletedCount,
			&item.ReviewCount,
			&item.SessionCount,
		); err != nil {
			log.Printf("failed to scan expert: %v", err)
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

func (r *Repository) CreateExpert(e *model.Expert) (*model.Expert, error) {
	const query = `
	INSERT INTO experts (
		id,
		username,
		gender,
		phone,
		email,
		password,
		specializations,
		education,
		document_urls,
		status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
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
	const query = `DELETE FROM experts WHERE id=$1`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.ID)
	return err
}

func (r *Repository) UpdateExpert(e *model.Expert) (*model.Expert, error) {
	const query = `
	UPDATE experts
	   SET username=$1,
		   gender=$2,
		   phone=$3,
		   education=$4,
		   specializations=$5,
           updated_at=now()
	 WHERE id=$6
`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query,
		e.Username,
		e.Gender,
		e.Phone,
		e.Education,
		e.Specializations,
		e.ID,
	)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertStatus(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET status=$1, updated_at=now() WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.Status, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertPassword(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET password=$1, updated_at=now() WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.Password, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertDocumentURLs(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET document_urls=$1, updated_at=now() WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.DocumentURLs, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertSpecializations(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET specializations=$1, updated_at=now() WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.Specializations, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}
