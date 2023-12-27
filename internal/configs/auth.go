package configs

import "time"

type Hash struct {
	Cost int `yaml:"cost"`
}

type Token struct {
	ExpiresIn                   string `yaml:"expires_in"`
	RegenerateTokenBeforeExpiry string `yaml:"regenerate_token_before_expiry"`
}

func (t Token) GetExpiresInDuration() (time.Duration, error) {
	return time.ParseDuration(t.ExpiresIn)
}

func (t Token) GetRegenerateTokenBeforeExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(t.RegenerateTokenBeforeExpiry)
}

type Auth struct {
	Hash  Hash
	Token Token
}
