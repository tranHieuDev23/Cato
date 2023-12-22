package rpc

type SubmissionStatus uint8
type SubmissionResult uint8

const (
	SubmissionStatusSubmitted SubmissionStatus = 1
	SubmissionStatusExecuting SubmissionStatus = 2
	SubmissionStatusFinished  SubmissionStatus = 3

	SubmissionResultOK                SubmissionResult = 1
	SubmissionResultCompileError      SubmissionResult = 2
	SubmissionResultRuntimeError      SubmissionResult = 3
	SubmissionResultTimeLimitExceeded SubmissionResult = 4
	SubmissionResultMemoryLimitExceed SubmissionResult = 5
	SubmissionResultWrongAnswer       SubmissionResult = 6
)

type SubmissionProblemSnippet struct {
	ID          uint64
	DisplayName string
}

type Submission struct {
	ID          uint64
	Problem     SubmissionProblemSnippet
	Author      Account
	Content     string           `validate:"max=5120"`
	Language    string           `validate:"max=32"`
	Status      SubmissionStatus `validate:"enum=1,2,3"`
	Result      SubmissionResult
	CreatedTime uint64
}

type SubmissionSnippet struct {
	ID          uint64
	Problem     SubmissionProblemSnippet
	Author      Account
	Language    string           `validate:"max=32"`
	Status      SubmissionStatus `validate:"enum=1,2,3"`
	Result      SubmissionResult `validate:"enum=1,2,3,4,5,6"`
	CreatedTime uint64
}
