package conf

type JWTConfig struct {
	SingingKey string `mapStructure:"signing_key"`
}
