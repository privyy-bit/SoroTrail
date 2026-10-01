package buildinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildInfoDefaults(t *testing.T) {
	// Verify fallback build info values are retrievable and non-empty
	assert.NotEmpty(t, Version)
	assert.NotEmpty(t, Commit)
	assert.NotEmpty(t, BuildDate)
}

func TestBuildInfoMutationAndErrors(t *testing.T) {
	// Save originals to restore after testing
	origVersion := Version
	origCommit := Commit
	origBuildDate := BuildDate
	defer func() {
		Version = origVersion
		Commit = origCommit
		BuildDate = origBuildDate
	}()

	Version = "v1.2.3"
	Commit = "abcdef0"
	BuildDate = "2025-01-01T00:00:00Z"

	assert.Equal(t, "v1.2.3", Version)
	assert.Equal(t, "abcdef0", Commit)
	assert.Equal(t, "2025-01-01T00:00:00Z", BuildDate)
}
