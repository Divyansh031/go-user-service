package unit

import (
	"os"
	"testing"

	"github.com/Divyansh031/user-service/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("ENV", "test")
	os.Setenv("GRPC_PORT", "50051")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("SCYLLA_HOSTS", "localhost")
	os.Setenv("SCYLLA_PORT", "9042")
	os.Setenv("SCYLLA_KEYSPACE", "userservice")
	os.Setenv("SCYLLA_CONSISTENCY", "QUORUM")
	os.Setenv("LOG_LEVEL", "info")

	// Remove CONFIG_PATH to force env loading
	os.Unsetenv("CONFIG_PATH")

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "test", cfg.Env)
	assert.Equal(t, 50051, cfg.GRPC.Port)
	assert.Equal(t, 8080, cfg.HTTP.Port)
	assert.Equal(t, "userservice", cfg.ScyllaDB.Keyspace)
	assert.Equal(t, "QUORUM", cfg.ScyllaDB.Consistency)
	assert.Equal(t, "info", cfg.Log.Level)
}

func TestLoadConfigDefaults(t *testing.T) {
	// Clear all env vars
	os.Unsetenv("ENV")
	os.Unsetenv("GRPC_PORT")
	os.Unsetenv("HTTP_PORT")
	os.Unsetenv("CONFIG_PATH")

	cfg, err := config.Load()

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "development", cfg.Env) // default
	assert.Equal(t, 50051, cfg.GRPC.Port)   // default
	assert.Equal(t, 8080, cfg.HTTP.Port)    // default
}
