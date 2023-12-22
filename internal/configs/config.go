package configs

type ConfigFilePath string

type Config struct {
	Hash  Hash
	Token Token
}

func NewConfig(filePath ConfigFilePath) (Config, error) {
	return Config{}, nil
}
