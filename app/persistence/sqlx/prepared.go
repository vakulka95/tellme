package sqlx

import "github.com/jmoiron/sqlx"

func prepareStatements(db *sqlx.DB, raw map[string]string) (map[string]*sqlx.Stmt, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for name, query := range raw {
		stmt, err := db.Preparex(query)
		if err != nil {
			return nil, err
		}
		stmts[name] = stmt
	}
	return stmts, nil
}
