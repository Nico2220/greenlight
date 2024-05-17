package data

import (
	"context"
	"database/sql"
	"time"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for _, v := range p {
		if code == v {

			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) getAllForUser(userID int64) (Permissions, error) {
	query := `
			SELECT permissions.code
			FROM permissions
			INNER JOIN users_permissions ON users_permissions.permission_id  = permissions.id
			INNER JOIN users ON users.id = users_permissions.user_id
			WHERE user.id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
