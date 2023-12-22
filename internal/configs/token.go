package configs

import "time"

type Token struct {
	PrivateKey                  string
	PublicKey                   string
	ExpiresIn                   string
	RegenerateTokenBeforeExpiry string
}

func (t Token) GetExpiresInDuration() (time.Duration, error) {
	return time.ParseDuration(t.ExpiresIn)
}

func (t Token) GetRegenerateTokenBeforeExpiryDuration() (time.Duration, error) {
	return time.ParseDuration(t.RegenerateTokenBeforeExpiry)
}
