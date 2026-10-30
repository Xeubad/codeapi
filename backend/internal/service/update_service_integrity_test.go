//go:build unit

package service

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsReleaseArchiveAssetMatchesOnlyPlatformArchive(t *testing.T) {
	t.Parallel()

	require.True(t, isReleaseArchiveAsset("sub2api_0.1.132_linux_amd64.tar.gz", "linux_amd64"))
	require.True(t, isReleaseArchiveAsset("sub2api_0.1.132_linux_amd64.tgz", "linux_amd64"))
	require.False(t, isReleaseArchiveAsset("sub2api_0.1.132_linux_amd64.tar.gz.sig", "linux_amd64"))
	require.False(t, isReleaseArchiveAsset("sub2api_0.1.132_linux_amd64.tar.gz.sha256", "linux_amd64"))
	require.False(t, isReleaseArchiveAsset("checksums.txt", "linux_amd64"))
	require.False(t, isReleaseArchiveAsset("sub2api_0.1.132_linux_arm64.tar.gz", "linux_amd64"))
	require.False(t, isReleaseArchiveAsset("sub2api_0.1.132_windows_amd64.zip", "linux_amd64"))
}

func TestExtractBinaryHandlesTgzArchives(t *testing.T) {
	t.Parallel()

	svc := &UpdateService{}
	archivePath := filepath.Join(t.TempDir(), "sub2api_0.1.132_linux_amd64.tgz")
	destPath := filepath.Join(t.TempDir(), "sub2api")

	f, err := os.Create(archivePath)
	require.NoError(t, err)
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	payload := []byte("binary-payload")
	require.NoError(t, tw.WriteHeader(&tar.Header{
		Name:     "sub2api",
		Mode:     0o755,
		Size:     int64(len(payload)),
		Typeflag: tar.TypeReg,
	}))
	_, err = tw.Write(payload)
	require.NoError(t, err)
	require.NoError(t, tw.Close())
	require.NoError(t, gw.Close())
	require.NoError(t, f.Close())

	require.NoError(t, svc.extractBinary(archivePath, destPath))
	got, err := os.ReadFile(destPath)
	require.NoError(t, err)
	require.Equal(t, payload, got)
}
