package rpc

//go:generate genpjrpc -search.name=API -print.place.path_swagger_file=../../../../api/swagger.json
type API interface {
	CreateAccount(CreateAccountRequest) CreateAccountResponse
	GetAccountList(GetAccountListRequest) GetAccountListResponse
	GetAccount(GetAccountRequest) GetAccountResponse
	UpdateAccount(UpdateAccountRequest) UpdateAccountResponse

	CreateSession(CreateSessionRequest) CreateSessionResponse
	DeleteSession(DeleteSessionRequest) DeleteSessionResponse

	CreateProblem(CreateProblemRequest) CreateProblemResponse
	GetProblemSnippetList(GetProblemSnippetListRequest) GetProblemSnippetListResponse
	GetProblem(GetProblemRequest) GetProblemResponse
	UpdateProblem(UpdateProblemRequest) UpdateProblemResponse
	DeleteProblem(DeleteProblemRequest) DeleteProblemResponse

	CreateTestCase(CreateTestCaseRequest) CreateTestCaseResponse
	CreateTestCaseList(CreateTestCaseListRequest) CreateTestCaseListResponse
	GetProblemTestCaseSnippetList(GetProblemTestCaseSnippetListRequest) GetProblemTestCaseSnippetListResponse
	GetTestCase(GetTestCaseRequest) GetTestCaseResponse
	UpdateTestCase(UpdateTestCaseRequest) UpdateTestCaseResponse
	DeleteTestCase(DeleteTestCaseRequest) DeleteTestCaseResponse

	GetAccountProblemSnippetList(GetAccountProblemSnippetListRequest) GetAccountProblemSnippetListResponse

	CreateSubmission(CreateSubmissionRequest) CreateSubmissionResponse
	GetSubmissionSnippetList(GetSubmissionSnippetListRequest) GetSubmissionSnippetListResponse
	GetSubmission(GetSubmissionRequest) GetSubmissionResponse
	DeleteSubmission(DeleteSubmissionRequest) DeleteSubmissionResponse

	GetAccountSubmissionSnippetList(GetAccountSubmissionSnippetListRequest) GetAccountSubmissionSnippetListResponse
	GetProblemSubmissionSnippetList(GetProblemSubmissionSnippetListRequest) GetProblemSubmissionSnippetListResponse
}

type CreateAccountRequest struct {
	AccountName string      `validate:"alphanum,min=6,max=32"`
	DisplayName string      `validate:"alphanumunicode,min=1,max=32"`
	Role        AccountRole `validate:"enum=admin,problem_setter,contestant,worker"`
	Password    string      `validate:"min=8"`
}

type CreateAccountResponse struct {
	Account Account
}

type GetAccountListRequest struct {
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetAccountListResponse struct {
	TotalAccountCount uint64
	AccountList       []Account
}

type GetAccountRequest struct {
	ID uint64
}

type GetAccountResponse struct {
	Account *Account
}

type UpdateAccountRequest struct {
	ID          uint64
	DisplayName string      `validate:"alphanumunicode,min=1,max=32"`
	Role        AccountRole `validate:"enum=admin,problem_setter,contestant,worker"`
}

type UpdateAccountResponse struct {
	Account Account
}

type CreateSessionRequest struct {
	AccountName string `validate:"alphanum,min=6,max=32"`
	Password    string `validate:"min=8"`
}

type CreateSessionResponse struct {
	Account Account
}

type DeleteSessionRequest struct{}

type DeleteSessionResponse struct{}

type CreateProblemRequest struct {
	DisplayName            string           `validate:"alphanumunicode,min=1,max=256"`
	Description            string           `validate:"max=5000"`
	TimeLimitInMillisecond uint64           `validate:"max=10000"`
	MemoryLimitInByte      uint64           `validate:"max=8589934592"`
	ExampleList            []ProblemExample `validate:"max=5"`
}

type CreateProblemResponse struct {
	Problem
}

type GetProblemSnippetListRequest struct {
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetProblemSnippetListResponse struct {
	TotalProblemCount  uint64
	ProblemSnippetList []Problem
}

type GetProblemRequest struct {
	ID uint64
}

type GetProblemResponse struct {
	Problem *Problem
}

type UpdateProblemRequest struct {
	ID                     uint64
	DisplayName            string           `validate:"alphanumunicode,min=1,max=256"`
	Description            string           `validate:"max=5000"`
	TimeLimitInMillisecond uint64           `validate:"max=10000"`
	MemoryLimitInByte      uint64           `validate:"max=8589934592"`
	ExampleList            []ProblemExample `validate:"max=5"`
}

type UpdateProblemResponse struct {
	Problem Problem
}

type DeleteProblemRequest struct {
	ID uint64
}

type DeleteProblemResponse struct{}

type CreateTestCaseRequest struct {
	ProblemID uint64
	Input     string `validate:"max=5000"`
	Output    string `validate:"max=5000"`
	IsHidden  bool
}

type CreateTestCaseResponse struct {
	TestCaseSnippet TestCaseSnippet
}

type CreateTestCaseListRequest struct {
	ProblemID      uint64
	ZippedTestData []byte
}

type CreateTestCaseListResponse struct{}

type GetProblemTestCaseSnippetListRequest struct {
	ProblemID uint64
	Offset    uint64
	Limit     uint64 `validate:"max=100"`
}

type GetProblemTestCaseSnippetListResponse struct {
	TotalTestCaseCount  uint64
	TestCaseSnippetList []TestCaseSnippet
}

type GetTestCaseRequest struct {
	ID uint64
}

type GetTestCaseResponse struct {
	TestCase TestCase
}

type UpdateTestCaseRequest struct {
	ID       uint64
	Input    string `validate:"max=5000"`
	Output   string `validate:"max=5000"`
	IsHidden bool
}

type UpdateTestCaseResponse struct {
	TestCase TestCase
}

type DeleteTestCaseRequest struct {
	ID uint64
}

type DeleteTestCaseResponse struct{}

type GetAccountProblemSnippetListRequest struct {
	AccountID uint64
	Offset    uint64
	Limit     uint64 `validate:"max=100"`
}

type GetAccountProblemSnippetListResponse struct {
	TotalProblemCount  uint64
	ProblemSnippetList []Problem
}

type CreateSubmissionRequest struct {
	ProblemID       uint64
	AuthorAccountID uint64
	Content         string `validate:"max=5120"`
	Language        string `validate:"max=32"`
}

type CreateSubmissionResponse struct {
	Problem Problem
}

type GetSubmissionSnippetListRequest struct {
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetSubmissionSnippetListResponse struct {
	TotalSubmissionCount  uint64
	SubmissionSnippetList []SubmissionSnippet
}

type GetSubmissionRequest struct {
	ID uint64
}

type GetSubmissionResponse struct {
	Submission *Submission
}

type DeleteSubmissionRequest struct {
	ID uint64
}

type DeleteSubmissionResponse struct{}

type GetAccountSubmissionSnippetListRequest struct {
	AccountID uint64
	Offset    uint64
	Limit     uint64 `validate:"max=100"`
}

type GetAccountSubmissionSnippetListResponse struct {
	TotalSubmissionCount  uint64
	SubmissionSnippetList []SubmissionSnippet
}

type GetProblemSubmissionSnippetListRequest struct {
	ProblemID uint64
	Offset    uint64
	Limit     uint64 `validate:"max=100"`
}

type GetProblemSubmissionSnippetListResponse struct {
	TotalSubmissionCount  uint64
	SubmissionSnippetList []SubmissionSnippet
}
