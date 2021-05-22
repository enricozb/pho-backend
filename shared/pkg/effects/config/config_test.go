package config_test

import (
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/config"
)

func TestConfig_Smoke(t *testing.T) {
	assert := require.New(t)

	dataPath, err := homedir.Expand("~/.pho")
	assert.NoError(err)
	assert.Equal(dataPath, config.Config.DataPath)
}
