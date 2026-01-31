package config

import "testing"

func TestDatabaseConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  *DatabaseConfig
	}{
		{
			name: "default database config",
			cfg: &DatabaseConfig{
				URI:             "postgres://localhost:5432",
				MaxOpenConns:    25,
				MaxIdleConns:    10,
				MinIdleConns:    5,
				ConnMaxLifetime: 3600,
			},
		},
		{
			name: "zero values",
			cfg: &DatabaseConfig{
				URI:             "",
				MaxOpenConns:    0,
				MaxIdleConns:    0,
				MinIdleConns:    0,
				ConnMaxLifetime: 0,
			},
		},
		{
			name: "custom values",
			cfg: &DatabaseConfig{
				URI:             "mysql://user:pass@localhost:3306/db",
				MaxOpenConns:    100,
				MaxIdleConns:    50,
				MinIdleConns:    10,
				ConnMaxLifetime: 7200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg

			if cfg.URI != tt.cfg.URI {
				t.Errorf("URI = %v, want %v", cfg.URI, tt.cfg.URI)
			}
			if cfg.MaxOpenConns != tt.cfg.MaxOpenConns {
				t.Errorf("MaxOpenConns = %v, want %v", cfg.MaxOpenConns, tt.cfg.MaxOpenConns)
			}
			if cfg.MaxIdleConns != tt.cfg.MaxIdleConns {
				t.Errorf("MaxIdleConns = %v, want %v", cfg.MaxIdleConns, tt.cfg.MaxIdleConns)
			}
			if cfg.MinIdleConns != tt.cfg.MinIdleConns {
				t.Errorf("MinIdleConns = %v, want %v", cfg.MinIdleConns, tt.cfg.MinIdleConns)
			}
			if cfg.ConnMaxLifetime != tt.cfg.ConnMaxLifetime {
				t.Errorf("ConnMaxLifetime = %v, want %v", cfg.ConnMaxLifetime, tt.cfg.ConnMaxLifetime)
			}
		})
	}
}

func TestDatabaseConfigIntegration(t *testing.T) {
	dbCfg := &DatabaseConfig{
		URI:             "postgres://localhost:5432/testdb",
		MaxOpenConns:    50,
		MaxIdleConns:    20,
		MinIdleConns:    5,
		ConnMaxLifetime: 1800,
	}

	cfg := New(WithDatabaseConfig(dbCfg))

	if cfg.Database == nil {
		t.Fatal("Expected Database config to be set, got nil")
	}

	if cfg.Database.URI != dbCfg.URI {
		t.Errorf("Database.URI = %v, want %v", cfg.Database.URI, dbCfg.URI)
	}
	if cfg.Database.MaxOpenConns != dbCfg.MaxOpenConns {
		t.Errorf("Database.MaxOpenConns = %v, want %v", cfg.Database.MaxOpenConns, dbCfg.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns != dbCfg.MaxIdleConns {
		t.Errorf("Database.MaxIdleConns = %v, want %v", cfg.Database.MaxIdleConns, dbCfg.MaxIdleConns)
	}
	if cfg.Database.MinIdleConns != dbCfg.MinIdleConns {
		t.Errorf("Database.MinIdleConns = %v, want %v", cfg.Database.MinIdleConns, dbCfg.MinIdleConns)
	}
	if cfg.Database.ConnMaxLifetime != dbCfg.ConnMaxLifetime {
		t.Errorf("Database.ConnMaxLifetime = %v, want %v", cfg.Database.ConnMaxLifetime, dbCfg.ConnMaxLifetime)
	}
}
