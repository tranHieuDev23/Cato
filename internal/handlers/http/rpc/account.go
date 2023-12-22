package rpc

type AccountRole string

const (
	AccountRoleAdmin         AccountRole = "admin"
	AccountRoleProblemSetter AccountRole = "problem_setter"
	AccountRoleContestant    AccountRole = "contestant"
	AccountRoleWorker        AccountRole = "worker"
)

type Account struct {
	ID          uint64
	AccountName string      `validate:"alphanum,min=6,max=32"`
	DisplayName string      `validate:"alphanumunicode,min=1,max=32"`
	Role        AccountRole `validate:"enum=admin,problem_setter,contestant,worker"`
}
