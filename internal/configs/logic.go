package configs

import (
	"time"

	"github.com/dustin/go-humanize"
)

type FirstAccount struct {
	AccountName string `yaml:"account_name"`
	DisplayName string `yaml:"display_name"`
	Password    string `yaml:"password"`
}

type FirstAccounts struct {
	Admin  FirstAccount `yaml:"admin"`
	Worker FirstAccount `yaml:"worker"`
}

type ProblemTestCaseHash struct {
	BatchSize uint64 `yaml:"batch_size"`
}

type Compile struct {
	Image           string   `yaml:"image"`
	CommandTemplate []string `yaml:"command_template"`
	Timeout         string   `yaml:"timeout"`
	CPUQuota        int64    `yaml:"cpu_quota"`
	Memory          string   `yaml:"memory"`
	WorkingDir      string   `yaml:"working_dir"`
	SourceFileName  string   `yaml:"source_file_name"`
	ProgramFileName string   `yaml:"program_file_name"`
}

func (c Compile) GetTimeoutDuration() (time.Duration, error) {
	return time.ParseDuration(c.Timeout)
}

func (c Compile) GetMemoryInBytes() (uint64, error) {
	return humanize.ParseBytes(c.Memory)
}

type TestCaseRun struct {
	Image           string   `yaml:"image"`
	CommandTemplate []string `yaml:"command_template"`
	CPUQuota        int64    `yaml:"cpu_quota"`
	WorkingDir      string   `yaml:"working_dir"`
}

type Language struct {
	Value       string      `yaml:"value"`
	Name        string      `yaml:"name"`
	Compile     *Compile    `yaml:"compile"`
	TestCaseRun TestCaseRun `yaml:"test_case_run"`
}

type Judge struct {
	Languages            []Language `yaml:"languages"`
	SubmissionRetryDelay string     `yaml:"submission_retry_delay"`
}

func (j Judge) GetSubmissionRetryDelayDuration() (time.Duration, error) {
	return time.ParseDuration(j.SubmissionRetryDelay)
}

type SyncProblem struct {
	Schedule                        string `yaml:"schedule"`
	GetProblemSnippetListBatchSize  uint64 `yaml:"get_problem_snippet_list_batch_size"`
	GetTestCaseSnippetListBatchSize uint64 `yaml:"get_test_case_snippet_list_batch_size"`
}

type Logic struct {
	FirstAccounts       FirstAccounts       `yaml:"first_accounts"`
	ProblemTestCaseHash ProblemTestCaseHash `yaml:"problem_test_case_hash"`
	Judge               Judge               `yaml:"judge"`
	SyncProblem         SyncProblem         `yaml:"sync_problem"`
}
