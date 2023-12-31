package rpc

type TestCase struct {
	UUID     string
	Input    string `validate:"max=5242880"`
	Output   string `validate:"max=5242880"`
	IsHidden bool
}

type TestCaseSnippet struct {
	UUID     string
	Input    string `validate:"max=100"`
	Output   string `validate:"max=100"`
	IsHidden bool
}
