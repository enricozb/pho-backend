package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/config"
)

func TestConfig_Smoke(t *testing.T) {
	assert := require.New(t)

	assert.Equal("~/.pho", config.Config.DataPath)
}
