package config_test

import (
	"os"
	"testing"

	"github.com/maitesin/hermes/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// unset environment variables
	variables := []string{
		"TELEGRAM_TOKEN",
		"DB_URL",
		"DB_SSL_MODE",
		"DB_BINARY_PARAMETERS",
	}
	for _, variable := range variables {
		err := os.Unsetenv(variable)
		require.NoError(t, err)
	}

	_, err := config.New()
	require.ErrorAs(t, err, &config.Error{})

	err = os.Setenv("TELEGRAM_TOKEN", "something good")
	require.Nil(t, err)

	cfg, err := config.New()
	require.Nil(t, err)

	require.Equal(t, "postgres://postgres:postgres@localhost:54321/hermes", cfg.SQL.URL)
	require.Equal(t, "disable", cfg.SQL.SSLMode)
	require.Equal(t, "yes", cfg.SQL.BinaryParams)
	require.Equal(t, "something good", cfg.Telegram.Token)

	err = os.Setenv("DB_URL", "https://oscarforner.com")
	require.Nil(t, err)

	cfg, err = config.New()
	require.Nil(t, err)

	require.Equal(t, "https://oscarforner.com", cfg.SQL.URL)
}
