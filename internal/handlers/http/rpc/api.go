package rpc

type ErrorCode int

const (
	ErrorCodeOK                 ErrorCode = 1
	ErrorCodeCanceled           ErrorCode = 2
	ErrorCodeUnknown            ErrorCode = 3
	ErrorCodeInvalidArgument    ErrorCode = 4
	ErrorCodeDeadlineExceeded   ErrorCode = 5
	ErrorCodeNotFound           ErrorCode = 6
	ErrorCodeAlreadyExists      ErrorCode = 7
	ErrorCodePermissionDenied   ErrorCode = 8
	ErrorCodeResourceExhausted  ErrorCode = 9
	ErrorCodeFailedPrecondition ErrorCode = 10
	ErrorCodeAborted            ErrorCode = 11
	ErrorCodeOutOfRange         ErrorCode = 12
	ErrorCodeUnimplemented      ErrorCode = 13
	ErrorCodeInternal           ErrorCode = 14
	ErrorCodeUnavailable        ErrorCode = 15
	ErrorCodeDataLoss           ErrorCode = 16
	ErrorCodeUnauthenticated    ErrorCode = 17
)

//go:generate genpjrpc -search.name=API -print.place.path_swagger_file=../../../../api/swagger.json
type API interface {
	GetServerInfo(GetServerInfoRequest) GetServerInfoResponse

	CreateAccount(CreateAccountRequest) CreateAccountResponse
	GetAccountList(GetAccountListRequest) GetAccountListResponse
	GetAccount(GetAccountRequest) GetAccountResponse
	UpdateAccount(UpdateAccountRequest) UpdateAccountResponse

	CreateSession(CreateSessionRequest) CreateSessionResponse
	GetSession(GetSessionRequest) GetSessionResponse
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
	UpdateSubmission(UpdateSubmissionRequest) UpdateSubmissionResponse
	DeleteSubmission(DeleteSubmissionRequest) DeleteSubmissionResponse

	GetAccountSubmissionSnippetList(GetAccountSubmissionSnippetListRequest) GetAccountSubmissionSnippetListResponse
	GetProblemSubmissionSnippetList(GetProblemSubmissionSnippetListRequest) GetProblemSubmissionSnippetListResponse
	GetAccountProblemSubmissionSnippetList(
		GetAccountProblemSubmissionSnippetListRequest,
	) GetAccountProblemSubmissionSnippetListResponse

	GetAndUpdateFirstSubmittedSubmissionToExecuting(
		GetAndUpdateFirstSubmittedSubmissionToExecutingRequest,
	) GetAndUpdateFirstSubmittedSubmissionToExecutingResponse

	UpdateSetting(UpdateSettingRequest) UpdateSettingResponse
}

type GetServerInfoRequest struct{}

type GetServerInfoResponse struct {
	IsDistributed         bool
	SupportedLanguageList []Language
	Setting               Setting
}

type CreateAccountRequest struct {
	AccountName string `validate:"alphanum,min=6,max=32"`
	DisplayName string `validate:"min=1,max=32"`
	Role        string `validate:"oneof=admin problem_setter contestant worker"`
	Password    string `validate:"min=8"`
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
	Account Account
}

type UpdateAccountRequest struct {
	ID          uint64
	DisplayName *string `validate:"omitnil,min=1,max=32"`
	Role        *string `validate:"omitnil,oneof=admin problem_setter contestant worker"`
	Password    *string `validate:"omitnil,min=8"`
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
	Token   string
}

type GetSessionRequest struct{}

type GetSessionResponse struct {
	Account Account
}

type DeleteSessionRequest struct{}

type DeleteSessionResponse struct{}

type CreateProblemRequest struct {
	DisplayName            string           `validate:"min=1,max=256"`
	Description            string           `validate:"max=5000"`
	TimeLimitInMillisecond uint64           `validate:"max=10000"`
	MemoryLimitInByte      uint64           `validate:"max=8589934592"`
	ExampleList            []ProblemExample `validate:"max=5"`
}

type CreateProblemResponse struct {
	Problem Problem
}

type GetProblemSnippetListRequest struct {
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetProblemSnippetListResponse struct {
	TotalProblemCount  uint64
	ProblemSnippetList []ProblemSnippet
}

type GetProblemRequest struct {
	UUID string
}

type GetProblemResponse struct {
	Problem Problem
}

type UpdateProblemRequest struct {
	UUID                   string
	DisplayName            *string           `omitnil,validate:"min=1,max=256"`
	Description            *string           `omitnil,validate:"max=5000"`
	TimeLimitInMillisecond *uint64           `omitnil,validate:"max=10000"`
	MemoryLimitInByte      *uint64           `omitnil,validate:"max=8589934592"`
	ExampleList            *[]ProblemExample `omitnil,validate:"max=5"`
}

type UpdateProblemResponse struct {
	Problem Problem
}

type DeleteProblemRequest struct {
	UUID string
}

type DeleteProblemResponse struct{}

type CreateTestCaseRequest struct {
	ProblemUUID string
	Input       string `validate:"max=5242880"`
	Output      string `validate:"max=5242880"`
	IsHidden    bool
}

type CreateTestCaseResponse struct {
	TestCaseSnippet TestCaseSnippet
}

type CreateTestCaseListRequest struct {
	ProblemUUID    string
	ZippedTestData string `validate:"max=5242880"`
}

type CreateTestCaseListResponse struct{}

type GetProblemTestCaseSnippetListRequest struct {
	ProblemUUID string
	Offset      uint64
	Limit       uint64 `validate:"max=100"`
}

type GetProblemTestCaseSnippetListResponse struct {
	TotalTestCaseCount  uint64
	TestCaseSnippetList []TestCaseSnippet
}

type GetTestCaseRequest struct {
	UUID string
}

type GetTestCaseResponse struct {
	TestCase TestCase
}

type UpdateTestCaseRequest struct {
	UUID     string
	Input    *string `validate:"omitnil,max=5242880"`
	Output   *string `validate:"omitnil,max=5242880"`
	IsHidden *bool
}

type UpdateTestCaseResponse struct {
	TestCaseSnippet TestCaseSnippet
}

type DeleteTestCaseRequest struct {
	UUID string
}

type DeleteTestCaseResponse struct{}

type GetAccountProblemSnippetListRequest struct {
	AccountID uint64
	Offset    uint64
	Limit     uint64 `validate:"max=100"`
}

type GetAccountProblemSnippetListResponse struct {
	TotalProblemCount  uint64
	ProblemSnippetList []ProblemSnippet
}

type CreateSubmissionRequest struct {
	ProblemUUID string
	Content     string `validate:"min=1,max=64000"`
	Language    string `validate:"max=32"`
}

type CreateSubmissionResponse struct {
	SubmissionSnippet SubmissionSnippet
}

type UpdateSubmissionRequest struct {
	ID     uint64
	Status uint8 `validate:"oneof=1 2 3"`
	Result uint8 `validate:"oneof=1 2 3 4 5 6 7"`
}

type UpdateSubmissionResponse struct {
	SubmissionSnippet SubmissionSnippet
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
	Submission Submission
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
	ProblemUUID string
	Offset      uint64
	Limit       uint64 `validate:"max=100"`
}

type GetProblemSubmissionSnippetListResponse struct {
	TotalSubmissionCount  uint64
	SubmissionSnippetList []SubmissionSnippet
}

type GetAccountProblemSubmissionSnippetListRequest struct {
	AccountID   uint64
	ProblemUUID string
	Offset      uint64
	Limit       uint64 `validate:"max=100"`
}

type GetAccountProblemSubmissionSnippetListResponse struct {
	TotalSubmissionCount  uint64
	SubmissionSnippetList []SubmissionSnippet
}

type GetAndUpdateFirstSubmittedSubmissionToExecutingRequest struct{}

type GetAndUpdateFirstSubmittedSubmissionToExecutingResponse struct {
	Submission Submission
}

type UpdateSettingRequest struct {
	Setting Setting
}

type UpdateSettingResponse struct {
	Setting Setting
}
