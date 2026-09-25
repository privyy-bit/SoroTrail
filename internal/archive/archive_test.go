package archive

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sorotrail/sorotrail/internal/store"
)

func TestOptionsDisabledByDefault(t *testing.T) {
	opts := Options{}
	assert.False(t, opts.Enabled(), "archive should be disabled when bucket is empty")
}

func TestOptionsEnabled(t *testing.T) {
	opts := Options{Bucket: "my-bucket"}
	assert.True(t, opts.Enabled(), "archive should be enabled when bucket is set")
}

func TestNewReturnsNilWhenDisabled(t *testing.T) {
	var st store.Store // nil interface
	a, err := New(st, Options{})
	assert.NoError(t, err)
	assert.Nil(t, a, "archiver should be nil when disabled")
}

func TestArchiverIsNilSafe(t *testing.T) {
	var a *Archiver
	// All methods should be safe to call on a nil archiver.
	archived, err := a.IsArchived(context.Background(), 100, 200)
	assert.NoError(t, err)
	assert.False(t, archived)

	uri, sha, err := a.ArchiveRange(context.Background(), 100, 200)
	assert.NoError(t, err)
	assert.Empty(t, uri)
	assert.Empty(t, sha)
}

func TestObjectPrefix(t *testing.T) {
	a := &Archiver{prefix: ""}
	assert.Equal(t, "events/schema=1/ledger_start=0240000", a.objectPrefix(240000))

	a2 := &Archiver{prefix: "my-prefix"}
	assert.Equal(t, "my-prefix/events/schema=1/ledger_start=0240000", a2.objectPrefix(240000))
}

func TestManifestHash(t *testing.T) {
	a := &Archiver{log: slog.Default()}
	m1 := Manifest{
		SchemaVersion: 1,
		Chunk: ChunkInfo{
			LedgerStart: 100,
			LedgerEnd:   200,
			RowCount:    50,
		},
		Producer: "sorotrail",
	}
	m2 := Manifest{
		SchemaVersion: 1,
		Chunk: ChunkInfo{
			LedgerStart: 100,
			LedgerEnd:   200,
			RowCount:    50,
		},
		Producer: "sorotrail",
	}
	// Same manifest should produce same hash.
	assert.Equal(t, a.manifestHash(m1), a.manifestHash(m2))

	// Different manifest should produce different hash.
	m3 := m1
	m3.Chunk.RowCount = 100
	assert.NotEqual(t, a.manifestHash(m1), a.manifestHash(m3))
}

func TestSchemaVersionConstant(t *testing.T) {
	assert.Equal(t, 1, SchemaVersion)
}

func TestArchiveErrorPaths(t *testing.T) {
	// Test context cancellation and error paths on Archiver operations
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var a *Archiver
	_, err := a.IsArchived(ctx, 100, 200)
	assert.NoError(t, err, "nil archiver is safe even with cancelled context")

	// Test with an uninitialized or failing store / bucket configuration
	opts := Options{
		Bucket:          "test-bucket",
		Endpoint:        "invalid-endpoint",
		UseSSL:          false,
		AccessKeyID:     "bad",
		SecretAccessKey: "worse",
	}
	// New with invalid minio client connection should error or handle appropriately
	var st store.Store
	arch, err := New(st, opts)
	// If client initialization fails or succeeds depending on lazy evaluation, test methods return wrapped errors.
	if err == nil && arch != nil {
		_, _, err = arch.ArchiveRange(ctx, 100, 200)
		assert.Error(t, err, "cancelled context or storage failure must be returned")
	}
}
