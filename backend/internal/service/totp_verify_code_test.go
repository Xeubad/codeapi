//go:build unit

package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type totpVerifyCacheStub struct {
	mu       sync.Mutex
	attempts int
	err      error
}

func (c *totpVerifyCacheStub) GetSetupSession(context.Context, int64) (*TotpSetupSession, error) {
	return nil, nil
}

func (c *totpVerifyCacheStub) SetSetupSession(context.Context, int64, *TotpSetupSession, time.Duration) error {
	return nil
}

func (c *totpVerifyCacheStub) DeleteSetupSession(context.Context, int64) error {
	return nil
}

func (c *totpVerifyCacheStub) GetLoginSession(context.Context, string) (*TotpLoginSession, error) {
	return nil, nil
}

func (c *totpVerifyCacheStub) SetLoginSession(context.Context, string, *TotpLoginSession, time.Duration) error {
	return nil
}

func (c *totpVerifyCacheStub) DeleteLoginSession(context.Context, string) error {
	return nil
}

func (c *totpVerifyCacheStub) IncrementVerifyAttempts(context.Context, int64) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err != nil {
		return 0, c.err
	}
	c.attempts++
	return c.attempts, nil
}

func (c *totpVerifyCacheStub) GetVerifyAttempts(context.Context, int64) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.attempts, nil
}

func (c *totpVerifyCacheStub) ClearVerifyAttempts(context.Context, int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.attempts = 0
	return nil
}

func (c *totpVerifyCacheStub) SetStepUpGrant(context.Context, int64, string, time.Duration) error {
	return nil
}

func (c *totpVerifyCacheStub) HasStepUpGrant(context.Context, int64, string) (bool, error) {
	return false, nil
}

type totpVerifyEncryptorStub struct{}

func (totpVerifyEncryptorStub) Encrypt(plaintext string) (string, error) {
	return plaintext, nil
}

func (totpVerifyEncryptorStub) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}

func newTotpVerifyService(cache TotpCache) (*TotpService, *totpVMUserRepoStub) {
	secret := "JBSWY3DPEHPK3PXP"
	userRepo := &totpVMUserRepoStub{user: &User{
		ID:                  7,
		Email:               "totp-verify@example.com",
		TotpEnabled:         true,
		TotpSecretEncrypted: &secret,
	}}
	return NewTotpService(userRepo, totpVerifyEncryptorStub{}, cache, nil, nil, nil), userRepo
}

func TestTotpVerifyCodeReservesAttemptsBeforeValidation(t *testing.T) {
	t.Parallel()

	cache := &totpVerifyCacheStub{}
	svc, _ := newTotpVerifyService(cache)

	const workers = 50
	var wg sync.WaitGroup
	var invalid, tooMany int
	var mu sync.Mutex
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := svc.VerifyCode(context.Background(), 7, "000000")
			mu.Lock()
			defer mu.Unlock()
			switch {
			case errors.Is(err, ErrTotpInvalidCode):
				invalid++
			case errors.Is(err, ErrTotpTooManyAttempts):
				tooMany++
			default:
				t.Errorf("unexpected verify result: %v", err)
			}
		}()
	}
	wg.Wait()

	require.Equal(t, maxTotpAttempts, invalid)
	require.Equal(t, workers-maxTotpAttempts, tooMany)
}

func TestTotpVerifyCodeFailsClosedWhenAttemptReservationFails(t *testing.T) {
	t.Parallel()

	cache := &totpVerifyCacheStub{err: errors.New("redis unavailable")}
	svc, _ := newTotpVerifyService(cache)

	err := svc.VerifyCode(context.Background(), 7, "000000")
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTotpInvalidCode)
	require.NotErrorIs(t, err, ErrTotpTooManyAttempts)
	attempts, attemptsErr := cache.GetVerifyAttempts(context.Background(), 7)
	require.NoError(t, attemptsErr)
	require.Zero(t, attempts)
}
