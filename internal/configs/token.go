package configs

import "time"

type Token struct {
	PrivateKey string
	PublicKey  string
	ExpiresIn  string
}

func (t Token) GetExpiresInDuration() (time.Duration, error) {
	return time.ParseDuration(t.ExpiresIn)
}
