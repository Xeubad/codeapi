package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// These limits must be supplied by the caller according to its storage quota.
// Archive inspection establishes integrity, not trust or permission to execute SQL.
type backupImportArchiveLimits struct {
	CompressedBytes   int64
	UncompressedBytes int64
}

type backupImportArchiveInfo struct {
	CompressedBytes   int64
	UncompressedBytes int64
	SHA256            string
}

var (
	errImportArchiveLimit  = errors.New("backup archive exceeds size limit")
	errImportArchiveFormat = errors.New("invalid backup archive")
	errImportArchiveDigest = errors.New("backup archive checksum mismatch")
)

const backupImportHeaderLimit = 64 * 1024

type backupImportCountingReader struct {
	ctx context.Context
	r   io.Reader
	n   int64
}

func (r *backupImportCountingReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := r.r.Read(p)
	r.n += int64(n)
	return n, err
}

type backupImportHeaderWriter struct {
	buf bytes.Buffer
}

func (w *backupImportHeaderWriter) Write(p []byte) (int, error) {
	if left := backupImportHeaderLimit - w.buf.Len(); left > 0 {
		if left > len(p) {
			left = len(p)
		}
		_, _ = w.buf.Write(p[:left])
	}
	return len(p), nil
}

// inspectBackupImportArchive consumes the complete stream, validates every gzip
// member and computes the digest over the original compressed bytes. A pg_dump
// header is only a format hint: arbitrary SQL and psql commands may follow it.
func inspectBackupImportArchive(ctx context.Context, src io.Reader, limits backupImportArchiveLimits, expectedSHA256 string) (backupImportArchiveInfo, error) {
	var info backupImportArchiveInfo
	if ctx == nil || src == nil || limits.CompressedBytes <= 0 || limits.UncompressedBytes <= 0 || limits.CompressedBytes == int64(^uint64(0)>>1) || limits.UncompressedBytes == int64(^uint64(0)>>1) {
		return info, fmt.Errorf("%w: invalid limits or input", errImportArchiveFormat)
	}
	var expected []byte
	if expectedSHA256 != "" {
		var err error
		expected, err = hex.DecodeString(expectedSHA256)
		if err != nil || len(expected) != sha256.Size {
			return info, errImportArchiveDigest
		}
	}

	hash := sha256.New()
	compressed := &backupImportCountingReader{ctx: ctx, r: io.TeeReader(io.LimitReader(src, limits.CompressedBytes+1), hash)}
	gz, err := gzip.NewReader(compressed)
	if err != nil {
		if ctx.Err() != nil {
			return info, ctx.Err()
		}
		if compressed.n > limits.CompressedBytes {
			return info, errImportArchiveLimit
		}
		return info, fmt.Errorf("%w: gzip header", errImportArchiveFormat)
	}
	defer func() { _ = gz.Close() }()

	var header backupImportHeaderWriter
	uncompressed, readErr := io.Copy(&header, io.LimitReader(gz, limits.UncompressedBytes+1))
	if compressed.n > limits.CompressedBytes || uncompressed > limits.UncompressedBytes {
		return info, errImportArchiveLimit
	}
	if readErr != nil {
		if ctx.Err() != nil {
			return info, ctx.Err()
		}
		return info, fmt.Errorf("%w: gzip stream", errImportArchiveFormat)
	}
	if !bytes.Contains(header.buf.Bytes(), []byte("-- PostgreSQL database dump")) {
		return info, fmt.Errorf("%w: missing pg_dump header", errImportArchiveFormat)
	}
	sum := hash.Sum(nil)
	if expected != nil && subtle.ConstantTimeCompare(expected, sum) != 1 {
		return info, errImportArchiveDigest
	}
	info.CompressedBytes = compressed.n
	info.UncompressedBytes = uncompressed
	info.SHA256 = hex.EncodeToString(sum)
	return info, nil
}
