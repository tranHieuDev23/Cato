package rpc

type SubmissionStatus uint8

const (
	SubmissionStatusSubmitted SubmissionStatus = 1
	SubmissionStatusExecuting SubmissionStatus = 2
	SubmissionStatusFinished  SubmissionStatus = 3
)

type SubmissionProblemSnippet struct {
	ID          uint64
	DisplayName string
}

type Submission struct {
	ID          uint64
	Problem     SubmissionProblemSnippet
	Author      User
	Content     string           `validate:"max=5120"`
	Language    string           `validate:"max=32"`
	Status      SubmissionStatus `validate:"enum=1,2,3"`
	CreatedTime uint64
}

type SubmissionSnippet struct {
	ID          uint64
	Problem     SubmissionProblemSnippet
	Author      User
	Language    string           `validate:"max=32"`
	Status      SubmissionStatus `validate:"enum=1,2,3"`
	CreatedTime uint64
}
