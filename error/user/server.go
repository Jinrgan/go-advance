package user

import (
	"database/sql"
	"go-advance/error/dao"
	"log"
)
import "github.com/pkg/errors"

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
		log.Fatalf("original error: %T %v\n stack trace:\n%+v\n", errors.Cause(err), errors.Cause(err), err)
	}

	return &Entity{
		Id: id,
		User: &User{
			Name: user.Name,
		},
	}, nil
}
