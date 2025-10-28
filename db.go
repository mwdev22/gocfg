package config

type DatabaseConfig struct {
	URI             string
	MaxOpenConns    int
	MaxIdleConns    int
	MinIdleConns    int
	ConnMaxLifetime int
}
