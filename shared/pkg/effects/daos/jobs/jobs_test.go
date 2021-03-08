package jobs_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/enricozb/pho/shared/pkg/lib/testutil"
)

func setup(t *testing.T) (*require.Assertions, *gorm.DB, func()) {
	assert := require.New(t)
	db, cleanup := testutil.MockDB(t)

	return assert, db, cleanup
}
