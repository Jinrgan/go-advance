package id

type UserID string

func (u UserID) String() string {
	return string(u)
}
