package rpc

//go:generate genpjrpc -search.name=API -print.place.path_swagger_file=../../../../api/swagger.json
type API interface {
	CreateUser(CreateUserRequest) CreateUserResponse
	GetUserList(GetUserListRequest) GetUserListResponse
	GetUser(GetUserRequest) GetUserResponse
	UpdateUser(UpdateUserRequest) UpdateUserResponse

	CreateSession(CreateSessionRequest) CreateSessionResponse
	DeleteSession(DeleteSessionRequest) DeleteSessionResponse

	CreateProblem(CreateProblemRequest) CreateProblemResponse
	GetProblemSnippetList(GetProblemSnippetListRequest) GetProblemSnippetListResponse
	GetProblem(GetProblemRequest) GetProblemResponse
	UpdateProblem(UpdateProblemRequest) UpdateProblemResponse
	DeleteProblem(DeleteProblemRequest) DeleteProblemResponse

	GetUserProblemSnippetList(GetUserProblemSnippetListRequest) GetUserProblemSnippetListResponse

	CreateSubmission(CreateSubmissionRequest) CreateSubmissionResponse
	GetSubmissionSnippetList(GetSubmissionSnippetListRequest) GetSubmissionSnippetListResponse
	GetSubmission(GetSubmissionRequest) GetSubmissionResponse
	DeleteSubmission(DeleteSubmissionRequest) DeleteSubmissionResponse

	GetUserSubmissionSnippetList(GetUserSubmissionSnippetListRequest) GetUserSubmissionSnippetListResponse
	GetProblemSubmissionSnippetList(GetProblemSubmissionSnippetListRequest) GetProblemSubmissionSnippetListResponse
}

type CreateUserRequest struct {
	Username    string   `validate:"alphanum,min=6,max=32"`
	DisplayName string   `validate:"alphanumunicode,min=1,max=32"`
	Role        UserRole `validate:"enum=admin,problem_setter,contestant"`
	Password    string   `validate:"min=8"`
}

type CreateUserResponse struct {
	User User
}

type GetUserListRequest struct {
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetUserListResponse struct {
	TotalUserCount uint64
	UserList       []User
}

type GetUserRequest struct {
	ID uint64
}

type GetUserResponse struct {
	User *User
}

type UpdateUserRequest struct {
	ID          uint64
	DisplayName string   `validate:"alphanumunicode,min=1,max=32"`
	Role        UserRole `validate:"enum=admin,problem_setter,contestant"`
}

type UpdateUserResponse struct {
	User User
}

type CreateSessionRequest struct {
	Username string `validate:"alphanum,min=6,max=32"`
	Password string `validate:"min=8"`
}

type CreateSessionResponse struct {
	User User
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

type GetUserProblemSnippetListRequest struct {
	UserID uint64
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetUserProblemSnippetListResponse struct {
	TotalProblemCount  uint64
	ProblemSnippetList []Problem
}

type CreateSubmissionRequest struct {
	ProblemID    uint64
	AuthorUserID uint64
	Content      string `validate:"max=5120"`
	Language     string `validate:"max=32"`
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

type GetUserSubmissionSnippetListRequest struct {
	UserID uint64
	Offset uint64
	Limit  uint64 `validate:"max=100"`
}

type GetUserSubmissionSnippetListResponse struct {
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
