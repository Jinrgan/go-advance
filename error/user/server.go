package user

import (
	"database/sql"
	"errors"
	xerrors "github.com/pkg/errors"
	"go-advance/error/dao"
	"log"
)

type User struct {
	Name string
}

type Entity struct {
	Id   string
	User *User
}

type Server struct {
	mySQL *dao.MySQL
}

func (s *Server) GetUser(id string) (*Entity, error) {
	user, err := s.mySQL.GetUser(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		log.Printf("original error: %T %v\n stack trace:\n%+v\n", xerrors.Cause(err), xerrors.Cause(err), err)
	}

	return &Entity{
		Id: id,
		User: &User{
			Name: user.Name,
		},
	}, nil
}
