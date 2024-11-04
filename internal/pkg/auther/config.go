package auther

import "time"

type Config struct {
	JwtAccess             []byte        `env:"AUTH_JWT_SECRET_KEY" env-required:"true"`
	AccessExpirationTime  time.Duration `yaml:"accessExpirationTime" env-default:"5m"`
	RefreshExpirationTime time.Duration `yaml:"refreshExpirationTime" env-default:"24h"`
	Issuer                string
}
