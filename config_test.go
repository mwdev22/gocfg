package config

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		envVars   map[string]string
		opts      []func(*Config)
		wantAddr  string
		wantKey   string
		wantDB    *DatabaseConfig
	}{
		{
			name:     "default values when no env vars",
			envVars:  map[string]string{},
			opts:     nil,
			wantAddr: ":8080",
			wantKey:  "",
			wantDB:   nil,
		},
		{
			name: "custom addr from env",
			envVars: map[string]string{
				"ADDR": ":3000",
			},
			opts:     nil,
			wantAddr: ":3000",
			wantKey:  "",
			wantDB:   nil,
		},
		{
			name: "custom secret key from env",
			envVars: map[string]string{
				"SECRET_KEY": "mysecret123",
			},
			opts:     nil,
			wantAddr: ":8080",
			wantKey:  "mysecret123",
			wantDB:   nil,
		},
		{
			name: "all env vars set",
			envVars: map[string]string{
				"ADDR":       ":9000",
				"SECRET_KEY": "topsecret",
			},
			opts:     nil,
			wantAddr: ":9000",
			wantKey:  "topsecret",
			wantDB:   nil,
		},
		{
			name:    "with database config option",
			envVars: map[string]string{},
			opts: []func(*Config){
				WithDatabaseConfig(&DatabaseConfig{
					URI:          "postgres://localhost:5432",
					MaxOpenConns: 10,
				}),
			},
			wantAddr: ":8080",
			wantKey:  "",
			wantDB: &DatabaseConfig{
				URI:          "postgres://localhost:5432",
				MaxOpenConns: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars
			os.Clearenv()

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			cfg := New(tt.opts...)

			if cfg.Addr != tt.wantAddr {
				t.Errorf("New().Addr = %v, want %v", cfg.Addr, tt.wantAddr)
			}

			if string(cfg.SecretKey) != tt.wantKey {
				t.Errorf("New().SecretKey = %v, want %v", string(cfg.SecretKey), tt.wantKey)
			}

			if tt.wantDB == nil && cfg.Database != nil {
				t.Errorf("New().Database = %v, want nil", cfg.Database)
			} else if tt.wantDB != nil {
				if cfg.Database == nil {
					t.Errorf("New().Database = nil, want %v", tt.wantDB)
				} else {
					if cfg.Database.URI != tt.wantDB.URI {
						t.Errorf("New().Database.URI = %v, want %v", cfg.Database.URI, tt.wantDB.URI)
					}
					if cfg.Database.MaxOpenConns != tt.wantDB.MaxOpenConns {
						t.Errorf("New().Database.MaxOpenConns = %v, want %v", cfg.Database.MaxOpenConns, tt.wantDB.MaxOpenConns)
					}
				}
			}
		})
	}
}

func TestWithDatabaseConfig(t *testing.T) {
	dbCfg := &DatabaseConfig{
		URI:             "postgres://test",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		MinIdleConns:    2,
		ConnMaxLifetime: 3600,
	}

	cfg := &Config{
		Addr: ":8080",
	}

	opt := WithDatabaseConfig(dbCfg)
	opt(cfg)

	if cfg.Database != dbCfg {
		t.Errorf("WithDatabaseConfig() did not set database config correctly")
	}
	if cfg.Database.URI != "postgres://test" {
		t.Errorf("Database.URI = %v, want postgres://test", cfg.Database.URI)
	}
	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("Database.MaxOpenConns = %v, want 25", cfg.Database.MaxOpenConns)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		setEnv       bool
		want         string
	}{
		{
			name:         "env var is set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			setEnv:       true,
			want:         "custom",
		},
		{
			name:         "env var is not set",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			setEnv:       false,
			want:         "default",
		},
		{
			name:         "env var is empty string",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "",
			setEnv:       true,
			want:         "default",
		},
		{
			name:         "empty default value",
			key:          "TEST_VAR",
			defaultValue: "",
			envValue:     "value",
			setEnv:       true,
			want:         "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got := GetEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		setEnv       bool
		want         int
	}{
		{
			name:         "valid integer env var",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "42",
			setEnv:       true,
			want:         42,
		},
		{
			name:         "env var not set",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "",
			setEnv:       false,
			want:         10,
		},
		{
			name:         "env var is empty string",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "",
			setEnv:       true,
			want:         10,
		},
		{
			name:         "invalid integer value",
			key:          "TEST_INT",
			defaultValue: 10,
			envValue:     "not-a-number",
			setEnv:       true,
			want:         10,
		},
		{
			name:         "negative integer",
			key:          "TEST_INT",
			defaultValue: 0,
			envValue:     "-5",
			setEnv:       true,
			want:         -5,
		},
		{
			name:         "zero value",
			key:          "TEST_INT",
			defaultValue: 100,
			envValue:     "0",
			setEnv:       true,
			want:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			if tt.setEnv {
				os.Setenv(tt.key, tt.envValue)
			}

			got := GetEnvAsInt(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetEnvAsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
