package pgx

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
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
			review_count
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
		list = &model.ExpertList{Items: make([]*model.Expert, 0, 20)} // pseudo default capacity
	)

	query := `
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
					review_count
			  FROM v$experts
`

	count := `
			SELECT count(id) FROM experts
`
	listRows, countRow, err := r.queryExpertList(ctx, q, query, count)
	if err != nil {
		return nil, err
	}

	for listRows.Next() {
		item := &model.Expert{}
		if err := listRows.Scan(
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
		); err != nil {
			log.Printf("failed to scan expert: %v", err)
			continue
		}
		list.Items = append(list.Items, item)
	}

	if err := countRow.Scan(&list.Total); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *Repository) queryExpertList(ctx context.Context, q *model.QueryExpertList, query, count string) (pgx.Rows, pgx.Row, error) {
	switch {
	case q.Status != "" && len(q.Specializations) != 0:

		query = query + ` WHERE status=$3 AND specializations @> $4  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE status=$1 AND specializations @> $2`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset, q.Status, q.Specializations)
		if err != nil {
			return nil, nil, err
		}

		return rows, r.cli.QueryRow(ctx, count, q.Status, q.Specializations), nil

	case q.Status != "" && len(q.Specializations) == 0:

		query = query + ` WHERE status=$3  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE status=$1`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset, q.Status)
		if err != nil {
			return nil, nil, err
		}

		return rows, r.cli.QueryRow(ctx, count, q.Status), nil

	case q.Status == "" && len(q.Specializations) != 0:

		query = query + ` WHERE specializations @> $3  ORDER BY created_at desc LIMIT $1 OFFSET $2`
		count = count + ` WHERE specializations @> $1`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset, q.Specializations)
		if err != nil {
			return nil, nil, err
		}

		return rows, r.cli.QueryRow(ctx, count, q.Specializations), nil

	default:
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
			review_count
	  FROM v$experts
  ORDER BY created_at desc
	 LIMIT $1
	OFFSET $2
`

		rows, err := r.cli.Query(ctx, query, q.Limit, q.Offset)
		if err != nil {
			return nil, nil, err
		}

		const countQuery = `SELECT COUNT(id) FROM experts`
		return rows, r.cli.QueryRow(ctx, countQuery), nil
	}
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

func (r *Repository) UpdateExpertStatus(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET status=$1 WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.Status, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertPassword(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET password=$1 WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.Password, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}

func (r *Repository) UpdateExpertDocumentURLs(e *model.Expert) (*model.Expert, error) {
	const query = `UPDATE experts SET document_urls=$1 WHERE id=$2`

	var ctx = context.TODO()

	_, err := r.cli.Exec(ctx, query, e.DocumentURLs, e.ID)
	if err != nil {
		return nil, err
	}

	return &model.Expert{ID: e.ID}, nil
}
