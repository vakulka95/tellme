package pgx

import (
	"context"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

func (r *Repository) GetAdminByLogin(login string) (*model.Admin, error) {
	const query = `
	SELECT 	id,
			username,
			password,
			status,
			updated_at,
			created_at
	  FROM admins 
	 WHERE username=$1
`

	var (
		ctx   = context.TODO()
		admin = &model.Admin{}
	)

	err := r.cli.QueryRow(ctx, query, login).
		Scan(
			&admin.ID,
			&admin.Username,
			&admin.Password,
			&admin.Status,
			&admin.UpdatedAt,
			&admin.CreatedAt,
		)
	if err != nil {
		return nil, err
	}

	return admin, nil
}
