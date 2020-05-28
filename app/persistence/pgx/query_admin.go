package pgx

import (
	"context"

	"gitlab.com/tellmecomua/tellme.api/pkg/postgres"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

func (r *Repository) GetAdminByLogin(login string) (*model.Admin, error) {

	sql, args := postgres.NewQueryBuilder().
		Select(
			"id",
			"username",
			"password",
			"status",
			"updated_at",
			"created_at",
		).
		From("admins").
		Where(postgres.NewEqual("username", login)).
		Build()

	var (
		ctx   = context.TODO()
		admin = &model.Admin{}
	)

	err := r.cli.QueryRow(ctx, sql, args...).
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
