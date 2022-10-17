package config_test

import (
	"fmt"
	"testing"

	"github.com/maitesin/hermes/config"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	reason := "something bad happened"
	err := config.NewError(reason)

	require.ErrorAs(t, err, &config.Error{})
	require.Equal(t, err.Error(), fmt.Sprintf("unable to build configuration object: %s", reason))
}
