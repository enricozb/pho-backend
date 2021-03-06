package config_test

import (
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/config"
)

func TestConfig_Smoke(t *testing.T) {
	assert := require.New(t)

	dbDir, err := homedir.Expand("~/.pho/db")
	assert.NoError(err)
	assert.Equal(dbDir, config.Config.DBDir)

	mediaDir, err := homedir.Expand("~/.pho/imported/media")
	assert.NoError(err)
	assert.Equal(mediaDir, config.Config.MediaDir)
}
