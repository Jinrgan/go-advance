package dao

import (
	"database/sql"
	"github.com/pkg/errors"
)

type User struct {
	Name string
}

type UserRecord struct {
	Id string
	*User
}

type MySQL struct {
	DB *sql.DB
}

func (m *MySQL) GetUser(id string) (*UserRecord, error) {
	row := m.DB.QueryRow("select * from user where id = %s", id)
	if err := row.Err(); err != nil {
		return nil, errors.Wrap(err, "cannot query row")
	}

	var user UserRecord
	err := row.Scan(&user)
	if err != nil {
		return nil, errors.Wrap(err, "cannot scan row")
	}

	return &user, nil
}
