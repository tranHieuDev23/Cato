package rpc

type ProblemExample struct {
	ID     uint64
	Input  string `validate:"max=5000"`
	Output string `validate:"max=5000"`
}

type Problem struct {
	ID                     uint64
	DisplayName            string `validate:"alphanumunicode,min=1,max=256"`
	Author                 User
	Description            string           `validate:"max=5000"`
	TimeLimitInMillisecond uint64           `validate:"max=10000"`
	MemoryLimitInByte      uint64           `validate:"max=8589934592"`
	ExampleList            []ProblemExample `validate:"max=5"`
	CreatedTime            uint64
	UpdatedTime            uint64
}

type ProblemSnippet struct {
	ID                     uint64
	DisplayName            string `validate:"alphanumunicode,min=1,max=256"`
	Author                 User
	TimeLimitInMillisecond uint64           `validate:"max=10000"`
	MemoryLimitInByte      uint64           `validate:"max=8589934592"`
	ExampleList            []ProblemExample `validate:"max=5"`
	CreatedTime            uint64
	UpdatedTime            uint64
}
