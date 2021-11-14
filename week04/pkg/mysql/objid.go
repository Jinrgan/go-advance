package mysql

import (
	"fmt"
	"go-advance/week04/pkg/id"
	"strconv"
)

type ObjectID int

// ObjectIDFromID converts an id to objected id.
func ObjectIDFromID(id fmt.Stringer) (ObjectID, error) {
	oid, err := strconv.Atoi(id.String())
	if err != nil {
		return 0, err
	}

	return ObjectID(oid), err
}

// ObjectIDMustFromID converts an id to objected id, panics on error.
func ObjectIDMustFromID(id fmt.Stringer) ObjectID {
	oid, err := ObjectIDFromID(id)
	if err != nil {
		panic(err)
	}

	return oid
}

// ObjectIDToUserID converts object id to account id.
func ObjectIDToUserID(oid ObjectID) id.UserID {
	return id.UserID(strconv.Itoa(int(oid)))
}

func (id ObjectID) String() string {
	return strconv.Itoa(int(id))
}
