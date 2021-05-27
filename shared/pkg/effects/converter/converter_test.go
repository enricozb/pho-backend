package converter_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/enricozb/pho/shared/pkg/effects/converter"
	"github.com/enricozb/pho/shared/pkg/lib/file"
	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func TestMediaConverter_Smoke(t *testing.T) {
	assert := require.New(t)

	tmp, err := os.MkdirTemp("", "pho-tests-converter-*")
	assert.NoError(err)

	// convert the files
	converter := converter.NewMediaConverter()

	var supportedCount int64

	err = filepath.WalkDir(testutil.MediaFixturesPath, func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() {
			return nil
		}

		relpath, err := filepath.Rel(testutil.MediaFixturesPath, path)
		assert.NoError(err)

		// create the destination directory
		assert.NoError(os.MkdirAll(filepath.Join(tmp, filepath.Dir(relpath)), os.ModePerm))
		dstpath := filepath.Join(tmp, relpath)

		if isSupported, _, mimetype := file.Kind(path); isSupported {
			supportedCount++
			assert.NoError(converter.Convert(path, dstpath, mimetype))
		}

		return nil
	})

	assert.NoError(err)
	assert.NoError(converter.Finish())
	assert.Equal(supportedCount, testutil.NumFilesInFixture)
}
