package configs

import (
	"time"

	"github.com/dustin/go-humanize"
)

type FirstAdmin struct {
	AccountName string `yaml:"account_name"`
	DisplayName string `yaml:"display_name"`
	Password    string `yaml:"password"`
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
	ProgramFileName string   `yaml:"program_file_name"`
}

type Language struct {
	Compile     *Compile    `yaml:"compile"`
	TestCaseRun TestCaseRun `yaml:"test_case_run"`
}

type Judge struct {
	Languages map[string]Language `yaml:"languages"`
}

type Logic struct {
	FirstAdmin          FirstAdmin          `yaml:"first_admin"`
	ProblemTestCaseHash ProblemTestCaseHash `yaml:"problem_test_case_hash"`
	Judge               Judge               `yaml:"judge"`
}
