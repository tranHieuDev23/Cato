package rpc

type ProblemExample struct {
	Input  string `validate:"max=5000"`
	Output string `validate:"max=5000"`
}

type Problem struct {
	UUID                   string
	DisplayName            string `validate:"min=1,max=256"`
	Author                 Account
	Description            string           `validate:"max=5000"`
	TimeLimitInMillisecond uint64           `validate:"max=10000"`
	MemoryLimitInByte      uint64           `validate:"max=8589934592"`
	ExampleList            []ProblemExample `validate:"max=5"`
	CreatedTime            uint64
	UpdatedTime            uint64
}

type ProblemSnippet struct {
	UUID                   string
	DisplayName            string `validate:"min=1,max=256"`
	Author                 Account
	TimeLimitInMillisecond uint64 `validate:"max=10000"`
	MemoryLimitInByte      uint64 `validate:"max=8589934592"`
	CreatedTime            uint64
	UpdatedTime            uint64
}
