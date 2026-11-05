// Platform-managed connection pool overrides.
//
// Administrators can change these four pool values from the ops UI. The values
// are stored in a managed JSON file inside the resolved data directory, take
// precedence over config.yaml and environment variables, and are applied on the
// next process start. Removing the file restores the baseline configuration.
//
// Recovery contract: if the file is missing, unreadable or invalid at startup,
// the baseline values are used and the service still starts; the ops status
// endpoint reports the failure so the administrator can fix or remove the file.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ManagedPoolFileName is the managed override file stored under the resolved data dir.
const ManagedPoolFileName = "managed-pool.json"

const managedPoolFileVersion = 1

// Bounds enforced both when saving through the admin API and when loading the
// managed file at startup. config.Validate() only checks positivity and the
// idle<=open relationship; these caps keep a bad edit from allocating absurd
// pools on a 1C1G-class host. They are generous enough to stay out of the way
// of legitimate deployments (a PostgreSQL server rarely supports >2000 clients).
const (
	ManagedPoolMaxDBOpenConns   = 2000
	ManagedPoolMaxRedisPoolSize = 8192
)

// ManagedPoolValues holds the four pool settings that can be overridden from
// the admin UI. Field names mirror the config keys they replace.
type ManagedPoolValues struct {
	DBMaxOpenConns    int `json:"db_max_open_conns"`
	DBMaxIdleConns    int `json:"db_max_idle_conns"`
	RedisPoolSize     int `json:"redis_pool_size"`
	RedisMinIdleConns int `json:"redis_min_idle_conns"`
}

// managedPoolFile is the on-disk envelope. Version allows future migrations.
type managedPoolFile struct {
	Version   int               `json:"version"`
	UpdatedAt string            `json:"updated_at"`
	Values    ManagedPoolValues `json:"values"`
}

// maxManagedPoolFileBytes bounds the managed file size; it is a four-integer
// JSON document, anything larger indicates corruption or tampering.
const maxManagedPoolFileBytes = 8 << 10

var managedPoolFileMu sync.Mutex

// ManagedPoolPath resolves the managed file location. It mirrors
// setup.GetDataDir ordering without importing the setup package (which imports
// this package): DATA_DIR env > /app/data > working directory. The writability
// probe used by setup is intentionally skipped here; write errors surface at
// save time instead of creating probe files on every read.
func ManagedPoolPath() string {
	if dir := stringOrEmpty(os.Getenv("DATA_DIR")); dir != "" {
		return filepath.Join(dir, ManagedPoolFileName)
	}
	const dockerDataDir = "/app/data"
	if info, err := os.Stat(dockerDataDir); err == nil && info.IsDir() {
		return filepath.Join(dockerDataDir, ManagedPoolFileName)
	}
	return ManagedPoolFileName
}

func stringOrEmpty(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' && s[i] != '\t' && s[i] != '\n' && s[i] != '\r' {
			return s
		}
	}
	return ""
}

// Validate enforces the same invariants as Config.Validate plus the managed
// bounds. Called on save and on load so a hand-edited file can never push the
// pools outside the accepted envelope.
func (v ManagedPoolValues) Validate() error {
	if v.DBMaxOpenConns <= 0 {
		return errors.New("db_max_open_conns must be positive")
	}
	if v.DBMaxOpenConns > ManagedPoolMaxDBOpenConns {
		return fmt.Errorf("db_max_open_conns must not exceed %d", ManagedPoolMaxDBOpenConns)
	}
	if v.DBMaxIdleConns < 0 {
		return errors.New("db_max_idle_conns must be non-negative")
	}
	if v.DBMaxIdleConns > v.DBMaxOpenConns {
		return errors.New("db_max_idle_conns cannot exceed db_max_open_conns")
	}
	if v.RedisPoolSize <= 0 {
		return errors.New("redis_pool_size must be positive")
	}
	if v.RedisPoolSize > ManagedPoolMaxRedisPoolSize {
		return fmt.Errorf("redis_pool_size must not exceed %d", ManagedPoolMaxRedisPoolSize)
	}
	if v.RedisMinIdleConns < 0 {
		return errors.New("redis_min_idle_conns must be non-negative")
	}
	if v.RedisMinIdleConns > v.RedisPoolSize {
		return errors.New("redis_min_idle_conns cannot exceed redis_pool_size")
	}
	return nil
}

// CurrentPoolValues reads the effective pool values from the loaded config.
func CurrentPoolValues(c *Config) ManagedPoolValues {
	return ManagedPoolValues{
		DBMaxOpenConns:    c.Database.MaxOpenConns,
		DBMaxIdleConns:    c.Database.MaxIdleConns,
		RedisPoolSize:     c.Redis.PoolSize,
		RedisMinIdleConns: c.Redis.MinIdleConns,
	}
}

// ReadManagedPool loads the managed file. It returns (nil, nil) when the file
// does not exist, and an error when the file exists but cannot be used — the
// caller decides whether that is fatal (never for reads; writes may overwrite).
func ReadManagedPool() (*managedPoolFile, error) {
	return readManagedPoolAt(ManagedPoolPath())
}

func readManagedPoolAt(path string) (*managedPoolFile, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat managed pool file: %w", err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("managed pool path %q is not a regular file", path)
	}
	if info.Size() > maxManagedPoolFileBytes {
		return nil, fmt.Errorf("managed pool file exceeds %d bytes", maxManagedPoolFileBytes)
	}
	//nolint:gosec // path is derived from DATA_DIR/fixed candidates, not request input
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read managed pool file: %w", err)
	}
	var f managedPoolFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse managed pool file: %w", err)
	}
	if f.Version != managedPoolFileVersion {
		return nil, fmt.Errorf("unsupported managed pool file version %d", f.Version)
	}
	if err := f.Values.Validate(); err != nil {
		return nil, fmt.Errorf("managed pool values invalid: %w", err)
	}
	return &f, nil
}

// WriteManagedPool validates and atomically persists the managed values.
func WriteManagedPool(values ManagedPoolValues) error {
	if err := values.Validate(); err != nil {
		return err
	}
	managedPoolFileMu.Lock()
	defer managedPoolFileMu.Unlock()

	path := ManagedPoolPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create data dir %q: %w", dir, err)
	}
	if info, err := os.Lstat(path); err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("managed pool path %q is not a regular file", path)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat managed pool file: %w", err)
	}

	body, err := json.MarshalIndent(managedPoolFile{
		Version:   managedPoolFileVersion,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Values:    values,
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode managed pool file: %w", err)
	}
	tmp, err := os.CreateTemp(dir, ".managed-pool-*")
	if err != nil {
		return fmt.Errorf("create temp managed pool file: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()

	if err := tmp.Chmod(0o600); err == nil {
		_, err = tmp.Write(body)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("write managed pool file: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("activate managed pool file: %w", err)
	}
	return nil
}

// RemoveManagedPool deletes the managed file; the next start falls back to the
// config.yaml/env baseline. Missing file is not an error.
func RemoveManagedPool() error {
	managedPoolFileMu.Lock()
	defer managedPoolFileMu.Unlock()
	path := ManagedPoolPath()
	if info, err := os.Lstat(path); err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("managed pool path %q is not a regular file", path)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("stat managed pool file: %w", err)
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove managed pool file: %w", err)
	}
	return nil
}

// applyManagedPoolOverride records the baseline (config.yaml + env) values and
// applies the managed file when present and valid. A bad file is logged and
// ignored so the service keeps starting with the baseline; the admin UI reports
// ManagedPoolIgnored via the status endpoint.
func (c *Config) applyManagedPoolOverride() {
	c.PoolBaseline = CurrentPoolValues(c)
	c.ManagedPoolApplied = false
	c.ManagedPoolIgnored = ""

	f, err := ReadManagedPool()
	if err != nil {
		c.ManagedPoolIgnored = err.Error()
		slog.Warn("managed pool override ignored, using baseline values", "error", err)
		return
	}
	if f == nil {
		return
	}
	c.Database.MaxOpenConns = f.Values.DBMaxOpenConns
	c.Database.MaxIdleConns = f.Values.DBMaxIdleConns
	c.Redis.PoolSize = f.Values.RedisPoolSize
	c.Redis.MinIdleConns = f.Values.RedisMinIdleConns
	c.ManagedPoolApplied = true
	slog.Info("managed pool override applied",
		"db_max_open_conns", f.Values.DBMaxOpenConns,
		"db_max_idle_conns", f.Values.DBMaxIdleConns,
		"redis_pool_size", f.Values.RedisPoolSize,
		"redis_min_idle_conns", f.Values.RedisMinIdleConns,
		"updated_at", f.UpdatedAt)
}
