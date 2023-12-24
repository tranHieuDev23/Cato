package configs

type FirstAdmin struct {
	AccountName string `yaml:"account_name"`
	DisplayName string `yaml:"display_name"`
	Password    string `yaml:"password"`
}

type Logic struct {
	FirstAdmin FirstAdmin `yaml:"first_admin"`
}
