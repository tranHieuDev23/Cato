package rpc

type AccountSetting struct {
	DisableAccountCreation                 bool `json:"DisableAccountCreation"`
	DisableAccountUpdate                   bool `json:"DisableAccountUpdate"`
	DisableSessionCreationForContestant    bool `json:"DisableSessionCreationForContestant"`
	DisableSessionCreationForProblemSetter bool `json:"DisableSessionCreationForProblemSetter"`
}

type ProblemSetting struct {
	DisableProblemCreation bool `json:"DisableProblemCreation"`
	DisableProblemUpdate   bool `json:"DisableProblemUpdate"`
}

type SubmissionSetting struct {
	DisableSubmissionCreation bool `json:"DisableSubmissionCreation"`
}

type Setting struct {
	Account    AccountSetting    `json:"Account"`
	Problem    ProblemSetting    `json:"Problem"`
	Submission SubmissionSetting `json:"Submission"`
}
