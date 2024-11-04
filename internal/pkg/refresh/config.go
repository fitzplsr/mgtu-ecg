package refresh

import "time"

type Config struct {
	RefreshExpirationTime time.Duration `yaml:"refreshExpirationTime" env-default:"24h"`
}
