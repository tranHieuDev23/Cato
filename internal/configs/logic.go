package configs

type FirstAdmin struct {
	AccountName string `yaml:"account_name"`
	DisplayName string `yaml:"display_name"`
	Password    string `yaml:"password"`
}

type ProblemTestCaseHash struct {
	BatchSize uint64 `yaml:"batch_size"`
}

type Logic struct {
	FirstAdmin          FirstAdmin          `yaml:"first_admin"`
	ProblemTestCaseHash ProblemTestCaseHash `yaml:"problem_test_case_hash"`
}
