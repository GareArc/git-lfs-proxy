package config

type DBConfig struct {
	Filename            string `mapstructure:"filename"`
	Username            string `mapstructure:"username"`
	Password            string `mapstructure:"password"`
	EncryptionAlgorithm string `mapstructure:"crypt"`
	Timeout             int    `mapstructure:"timeout"`
}
