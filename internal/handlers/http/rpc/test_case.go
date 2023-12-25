package rpc

type TestCase struct {
	ID       uint64
	Input    string `validate:"max=5000"`
	Output   string `validate:"max=5000"`
	IsHidden bool
}

type TestCaseSnippet struct {
	ID       uint64
	Input    string `validate:"max=100"`
	Output   string `validate:"max=100"`
	IsHidden bool
}
