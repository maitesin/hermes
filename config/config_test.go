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

	require.Equal(t, "something good", cfg.Telegram.Token)

	cfg, err = config.New()
	require.Nil(t, err)
}
