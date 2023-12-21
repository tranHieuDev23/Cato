package rpc

type UserRole string

const (
	UserRoleAdmin         UserRole = "admin"
	UserRoleProblemSetter UserRole = "problem_setter"
	UserRoleContestant    UserRole = "contestant"
)

type User struct {
	ID          uint64
	Username    string   `validate:"alphanum,min=6,max=32"`
	DisplayName string   `validate:"alphanumunicode,min=1,max=32"`
	Role        UserRole `validate:"enum=admin,problem_setter,contestant"`
}
