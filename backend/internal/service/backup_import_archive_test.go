package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func makeImportTestGzip(t *testing.T, content string) []byte {
	t.Helper()
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return b.Bytes()
}

func TestInspectBackupImportArchive(t *testing.T) {
	content := "-- PostgreSQL database dump\nSELECT 1;\n"
	archive := makeImportTestGzip(t, content)
	sum := sha256.Sum256(archive)
	digest := hex.EncodeToString(sum[:])
	limits := backupImportArchiveLimits{CompressedBytes: int64(len(archive)), UncompressedBytes: int64(len(content))}
	corrupt := bytes.Clone(archive)
	corrupt[len(corrupt)-8] ^= 1 // gzip trailer CRC32

	info, err := inspectBackupImportArchive(context.Background(), bytes.NewReader(archive), limits, digest)
	require.NoError(t, err)
	require.Equal(t, digest, info.SHA256)
	require.EqualValues(t, len(archive), info.CompressedBytes)
	require.EqualValues(t, len(content), info.UncompressedBytes)

	tests := []struct {
		name   string
		data   []byte
		limits backupImportArchiveLimits
		digest string
		want   error
	}{
		{"compressed limit", archive, backupImportArchiveLimits{CompressedBytes: int64(len(archive) - 1), UncompressedBytes: limits.UncompressedBytes}, "", errImportArchiveLimit},
		{"uncompressed limit", archive, backupImportArchiveLimits{CompressedBytes: limits.CompressedBytes, UncompressedBytes: int64(len(content) - 1)}, "", errImportArchiveLimit},
		{"wrong checksum", archive, limits, hex.EncodeToString(make([]byte, sha256.Size)), errImportArchiveDigest},
		{"bad checksum", archive, limits, "invalid", errImportArchiveDigest},
		{"not gzip", []byte("not gzip"), limits, "", errImportArchiveFormat},
		{"missing header", makeImportTestGzip(t, "hello"), limits, "", errImportArchiveFormat},
		{"truncated gzip", archive[:len(archive)-1], limits, "", errImportArchiveFormat},
		{"corrupt checksum", corrupt, limits, "", errImportArchiveFormat},
		{"trailing junk", append(bytes.Clone(archive), 'x'), backupImportArchiveLimits{CompressedBytes: int64(len(archive) + 1), UncompressedBytes: limits.UncompressedBytes}, "", errImportArchiveFormat},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := inspectBackupImportArchive(context.Background(), bytes.NewReader(tc.data), tc.limits, tc.digest)
			require.ErrorIs(t, err, tc.want)
		})
	}
}

func TestInspectBackupImportArchiveCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := inspectBackupImportArchive(ctx, bytes.NewReader(nil), backupImportArchiveLimits{CompressedBytes: 1024, UncompressedBytes: 1024}, "")
	require.True(t, errors.Is(err, context.Canceled))
}
