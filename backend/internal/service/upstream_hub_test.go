package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var errUpstreamTestCatalogUnavailable = errors.New("catalog unavailable")

type stubUpstreamVendorRepo struct {
	vendors []UpstreamVendor
	entries map[int64][]UpstreamEntry
	created *UpstreamVendor
	updated *UpstreamVendor
	deleted int64
}

func (s *stubUpstreamVendorRepo) List(context.Context) ([]UpstreamVendor, error) {
	return s.vendors, nil
}

func (s *stubUpstreamVendorRepo) GetByID(_ context.Context, id int64) (*UpstreamVendor, error) {
	for i := range s.vendors {
		if s.vendors[i].ID == id {
			vendor := s.vendors[i]
			return &vendor, nil
		}
	}
	return nil, ErrUpstreamVendorNotFound
}

func (s *stubUpstreamVendorRepo) Create(_ context.Context, name string, budget int, _ *string, baseURLs []string) (*UpstreamVendor, error) {
	s.created = &UpstreamVendor{ID: 7, Name: name, ConcurrencyBudget: budget, BaseURLs: baseURLs, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return s.created, nil
}

func (s *stubUpstreamVendorRepo) Update(_ context.Context, id int64, name string, budget int, _ *string, baseURLs []string) (*UpstreamVendor, error) {
	s.updated = &UpstreamVendor{ID: id, Name: name, ConcurrencyBudget: budget, BaseURLs: baseURLs}
	return s.updated, nil
}

func (s *stubUpstreamVendorRepo) SoftDelete(_ context.Context, id int64) error {
	s.deleted = id
	return nil
}

func (s *stubUpstreamVendorRepo) ListEntries(_ context.Context, vendorID int64) ([]UpstreamEntry, error) {
	return s.entries[vendorID], nil
}

func (s *stubUpstreamVendorRepo) ListGroups(context.Context) ([]UpstreamGroupOption, error) {
	return nil, nil
}

func (s *stubUpstreamVendorRepo) SetEntryGroupEnabled(context.Context, int64, int64, bool) error {
	return nil
}

func (s *stubUpstreamVendorRepo) SetEntryGroups(context.Context, int64, []int64) error {
	return nil
}

func TestUpstreamVendorIDFromExtra(t *testing.T) {
	tests := []struct {
		name  string
		extra map[string]any
		want  int64
		ok    bool
	}{
		{name: "float64 from jsonb", extra: map[string]any{UpstreamVendorIDExtraKey: float64(12)}, want: 12, ok: true},
		{name: "int64", extra: map[string]any{UpstreamVendorIDExtraKey: int64(3)}, want: 3, ok: true},
		{name: "string", extra: map[string]any{UpstreamVendorIDExtraKey: "9"}, want: 9, ok: true},
		{name: "missing", extra: map[string]any{}, ok: false},
		{name: "zero", extra: map[string]any{UpstreamVendorIDExtraKey: float64(0)}, ok: false},
		{name: "nil map", extra: nil, ok: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := upstreamVendorIDFromExtra(tc.extra)
			require.Equal(t, tc.ok, ok)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestNormalizeUpstreamBaseURL(t *testing.T) {
	got, err := normalizeUpstreamBaseURL("  https://api.example.com/v1/  ")
	require.NoError(t, err)
	require.Equal(t, "https://api.example.com/v1", got)

	for _, invalid := range []string{"", "ftp://example.com", "not a url", "//example.com"} {
		_, err := normalizeUpstreamBaseURL(invalid)
		require.Error(t, err, invalid)
	}
}

func TestEnrichUpstreamVendor(t *testing.T) {
	vendor := UpstreamVendor{ConcurrencyBudget: 50}
	entries := []UpstreamEntry{
		{Concurrency: 30, Groups: []UpstreamEntryGroupRef{{GroupID: 1}, {GroupID: 2}}},
		{Concurrency: 30, Groups: []UpstreamEntryGroupRef{{GroupID: 2}}},
	}
	enrichUpstreamVendor(&vendor, entries)
	require.Equal(t, 2, vendor.EntryCount)
	require.Equal(t, 60, vendor.AllocatedConcurrency)
	require.Equal(t, 2, vendor.ReferencedGroupCount)
	require.True(t, vendor.OverBudget)
}

func TestEnrichUpstreamVendorBudgetZeroDisablesLimit(t *testing.T) {
	vendor := UpstreamVendor{ConcurrencyBudget: 0}
	enrichUpstreamVendor(&vendor, []UpstreamEntry{{Concurrency: 100}})
	require.False(t, vendor.OverBudget)
}

func TestCreateVendorValidation(t *testing.T) {
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, nil, nil)

	_, err := svc.CreateVendor(context.Background(), CreateUpstreamVendorInput{Name: "   "})
	require.Error(t, err)

	_, err = svc.CreateVendor(context.Background(), CreateUpstreamVendorInput{Name: "A", ConcurrencyBudget: -1})
	require.Error(t, err)

	vendor, err := svc.CreateVendor(context.Background(), CreateUpstreamVendorInput{Name: " 中转商 A ", ConcurrencyBudget: 50})
	require.NoError(t, err)
	require.Equal(t, "中转商 A", vendor.Name)
	require.Equal(t, 50, vendor.ConcurrencyBudget)
}

func TestDeleteVendorRefusesWithEntries(t *testing.T) {
	repo := &stubUpstreamVendorRepo{
		vendors: []UpstreamVendor{{ID: 1, Name: "A"}},
		entries: map[int64][]UpstreamEntry{1: {{AccountID: 9}}},
	}
	svc := NewUpstreamHubService(repo, nil, nil)
	require.ErrorIs(t, svc.DeleteVendor(context.Background(), 1), ErrUpstreamVendorNotEmpty)

	emptyRepo := &stubUpstreamVendorRepo{
		vendors: []UpstreamVendor{{ID: 2, Name: "B"}},
		entries: map[int64][]UpstreamEntry{},
	}
	svc = NewUpstreamHubService(emptyRepo, nil, nil)
	require.NoError(t, svc.DeleteVendor(context.Background(), 2))
	require.Equal(t, int64(2), emptyRepo.deleted)
}

func TestGetVendorDetailWarnsOverBudget(t *testing.T) {
	repo := &stubUpstreamVendorRepo{
		vendors: []UpstreamVendor{{ID: 3, Name: "C", ConcurrencyBudget: 10}},
		entries: map[int64][]UpstreamEntry{3: {{AccountID: 1, Concurrency: 8}, {AccountID: 2, Concurrency: 8}}},
	}
	svc := NewUpstreamHubService(repo, nil, nil)
	detail, err := svc.GetVendorDetail(context.Background(), 3)
	require.NoError(t, err)
	require.True(t, detail.Vendor.OverBudget)
	require.Len(t, detail.Warnings, 1)
}

// stubUpstreamAdmin 只实现上游聚合用到的三个方法，其余方法由嵌入的 nil 接口满足。
type stubUpstreamAdmin struct {
	AdminService
	account *Account
	updated *UpdateAccountInput
	created *CreateAccountInput
}

func (s *stubUpstreamAdmin) GetAccount(context.Context, int64) (*Account, error) {
	if s.account == nil {
		return nil, ErrUpstreamEntryNotFound
	}
	return s.account, nil
}

func (s *stubUpstreamAdmin) UpdateAccount(_ context.Context, _ int64, input *UpdateAccountInput) (*Account, error) {
	s.updated = input
	return s.account, nil
}

func (s *stubUpstreamAdmin) GetProxy(context.Context, int64) (*Proxy, error) {
	return &Proxy{ID: 3, Name: "ip-mgmt-proxy"}, nil
}

func (s *stubUpstreamAdmin) CreateAccount(_ context.Context, input *CreateAccountInput) (*Account, error) {
	s.created = input
	return &Account{
		ID:          99,
		Name:        input.Name,
		Platform:    input.Platform,
		Credentials: input.Credentials,
		Extra:       input.Extra,
	}, nil
}

type stubUpstreamProber struct {
	models []string
	err    error
	called bool

	catalog      []string
	catalogErr   error
	catalogCalls int

	measureFail bool
	lastAccount *Account
}

func (s *stubUpstreamProber) FetchUpstreamSupportedModels(context.Context, *Account) ([]string, error) {
	s.called = true
	return s.models, s.err
}

func (s *stubUpstreamProber) SuggestUpstreamModelsFromCatalog(context.Context, *Account) ([]string, error) {
	s.catalogCalls++
	return s.catalog, s.catalogErr
}

func (s *stubUpstreamProber) MeasureBaseURLLatencies(_ context.Context, account *Account, urls []string) []UpstreamBaseURLLatency {
	s.lastAccount = account
	out := make([]UpstreamBaseURLLatency, 0, len(urls))
	for index, url := range urls {
		if s.measureFail {
			out = append(out, UpstreamBaseURLLatency{URL: url, Error: "unreachable"})
			continue
		}
		out = append(out, UpstreamBaseURLLatency{URL: url, OK: true, LatencyMS: int64((index + 1) * 10)})
	}
	return out
}

func TestUpstreamGroupAllowsPlatform(t *testing.T) {
	tests := []struct {
		name        string
		group       string
		entry       string
		wantAllowed bool
	}{
		{name: "same concrete platform", group: PlatformDeepseek, entry: PlatformDeepseek, wantAllowed: true},
		{name: "different concrete platform", group: PlatformAnthropic, entry: PlatformDeepseek, wantAllowed: false},
		{name: "composite accepts CN provider", group: PlatformComposite, entry: PlatformDeepseek, wantAllowed: true},
		{name: "composite accepts openai", group: PlatformComposite, entry: PlatformOpenAI, wantAllowed: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantAllowed, upstreamGroupAllowsPlatform(tc.group, tc.entry))
		})
	}
}

func TestProbeEntryModelsNormalizesAndReportsExisting(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:    5,
		Extra: map[string]any{UpstreamVendorIDExtraKey: int64(2)},
		Credentials: map[string]any{
			"api_key":       "sk-x",
			"model_mapping": map[string]any{"deepseek-chat": "deepseek-chat"},
		},
	}}
	prober := &stubUpstreamProber{models: []string{" deepseek-chat ", "deepseek-reasoner", "deepseek-chat", ""}}
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, admin, prober)

	probe, err := svc.ProbeEntryModels(context.Background(), 5)
	require.NoError(t, err)
	require.True(t, prober.called)
	require.Equal(t, []string{"deepseek-chat", "deepseek-reasoner"}, probe.Models)
	require.Equal(t, []string{"deepseek-chat"}, probe.Existing)
}

func TestProbeEntryModelsRejectsNonHubAccount(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{ID: 6}}
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, admin, &stubUpstreamProber{})

	_, err := svc.ProbeEntryModels(context.Background(), 6)
	require.ErrorIs(t, err, ErrUpstreamEntryNotFound)
}

func TestProbeEntryModelsFallsBackToCatalog(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:          5,
		Platform:    PlatformDeepseek,
		Extra:       map[string]any{UpstreamVendorIDExtraKey: int64(2)},
		Credentials: map[string]any{"api_key": "sk-x", "base_url": "https://relay.example.com"},
	}}
	prober := &stubUpstreamProber{
		err:     &UpstreamModelSyncError{Kind: UpstreamModelSyncErrorUpstream, Message: "HTTP 404", StatusCode: 404},
		catalog: []string{"deepseek-reasoner", "deepseek-chat"},
	}
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, admin, prober)

	probe, err := svc.ProbeEntryModels(context.Background(), 5)
	require.NoError(t, err)
	require.Equal(t, 1, prober.catalogCalls)
	require.Equal(t, "models.dev", probe.Source)
	require.NotEmpty(t, probe.Warning)
	require.Equal(t, []string{"deepseek-chat", "deepseek-reasoner"}, probe.Models)
}

func TestProbeEntryModelsKeepsLiveErrorWhenCatalogUnavailable(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:          5,
		Platform:    PlatformDeepseek,
		Extra:       map[string]any{UpstreamVendorIDExtraKey: int64(2)},
		Credentials: map[string]any{"api_key": "sk-x"},
	}}
	liveErr := &UpstreamModelSyncError{Kind: UpstreamModelSyncErrorUpstream, Message: "HTTP 500", StatusCode: 500}
	prober := &stubUpstreamProber{err: liveErr, catalogErr: errUpstreamTestCatalogUnavailable}
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, admin, prober)

	_, err := svc.ProbeEntryModels(context.Background(), 5)
	require.ErrorIs(t, err, liveErr)
}

func TestAlternateUpstreamModelsURLs(t *testing.T) {
	tests := []struct {
		name    string
		primary string
		want    []string
	}{
		{
			name:    "versioned models path falls back to bare models",
			primary: "https://relay.example.com/v1/models",
			want:    []string{"https://relay.example.com/models"},
		},
		{
			name:    "bare models path falls back to versioned",
			primary: "https://relay.example.com/models",
			want:    []string{"https://relay.example.com/v1/models"},
		},
		{
			name:    "prefix path is preserved",
			primary: "https://relay.example.com/api/v1/models",
			want:    []string{"https://relay.example.com/api/models"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, alternateUpstreamModelsURLs(tc.primary))
		})
	}
}

func TestSetEntryModelsWritesMappingAndKeepsCredentials(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:   5,
		Name: "A-OpenAI",
		// 历史条目（type=upstream）在写模型时应顺带转为 apikey。
		Type:  AccountTypeUpstream,
		Extra: map[string]any{UpstreamVendorIDExtraKey: int64(2)},
		Credentials: map[string]any{
			"api_key":          "sk-x",
			"base_url":         "https://api.example.com/v1",
			"header_overrides": map[string]any{"x": "y"},
			"model_mapping":    map[string]any{"old": "old"},
		},
	}}
	repo := &stubUpstreamVendorRepo{entries: map[int64][]UpstreamEntry{2: {{AccountID: 5, VendorID: 2}}}}
	svc := NewUpstreamHubService(repo, admin, nil)

	entry, err := svc.SetEntryModels(context.Background(), 5, []string{"deepseek-reasoner", " deepseek-chat ", "deepseek-chat"})
	require.NoError(t, err)
	require.Equal(t, int64(5), entry.AccountID)
	require.NotNil(t, admin.updated)
	require.Equal(t, AccountTypeAPIKey, admin.updated.Type)
	require.Equal(t, "sk-x", admin.updated.Credentials["api_key"])
	require.Equal(t, "https://api.example.com/v1", admin.updated.Credentials["base_url"])
	require.Equal(t, map[string]any{"x": "y"}, admin.updated.Credentials["header_overrides"])
	require.Equal(t, map[string]any{
		"deepseek-chat":     "deepseek-chat",
		"deepseek-reasoner": "deepseek-reasoner",
	}, admin.updated.Credentials["model_mapping"])
}

func TestSetEntryModelsRejectsTooManyModels(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:    7,
		Extra: map[string]any{UpstreamVendorIDExtraKey: int64(3)},
	}}
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, admin, nil)

	models := make([]string, maxUpstreamEntryModels+1)
	for i := range models {
		models[i] = fmt.Sprintf("model-%d", i)
	}
	_, err := svc.SetEntryModels(context.Background(), 7, models)
	require.Error(t, err)
	require.Nil(t, admin.updated)
}

func TestStoredUpstreamEntryModelsSkipsEmptyMappings(t *testing.T) {
	credentials := map[string]any{
		"model_mapping": map[string]any{
			"a": "a",
			"b": "",
			"":  "c",
			"d": 1,
		},
	}
	require.Equal(t, []string{"a"}, storedUpstreamEntryModels(credentials))
	require.Empty(t, storedUpstreamEntryModels(map[string]any{}))
}

func TestNormalizeVendorBaseURLs(t *testing.T) {
	got, err := normalizeVendorBaseURLs([]string{
		" https://a.example.com/ ",
		"https://a.example.com",
		"http://b.example.com/v1",
		"   ",
	})
	require.NoError(t, err)
	require.Equal(t, []string{"https://a.example.com", "http://b.example.com/v1"}, got)

	_, err = normalizeVendorBaseURLs([]string{"ftp://a.example.com"})
	require.Error(t, err)

	many := make([]string, maxVendorBaseURLs+1)
	for i := range many {
		many[i] = fmt.Sprintf("https://host-%d.example.com", i)
	}
	_, err = normalizeVendorBaseURLs(many)
	require.Error(t, err)
}

func TestCreateVendorStoresBaseURLs(t *testing.T) {
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, nil, nil)
	vendor, err := svc.CreateVendor(context.Background(), CreateUpstreamVendorInput{
		Name:     "A",
		BaseURLs: []string{"https://a.example.com", "https://a.example.com/"},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"https://a.example.com"}, vendor.BaseURLs)
}

func TestMeasureBaseURLsForDraftRow(t *testing.T) {
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, nil, &stubUpstreamProber{})
	results, err := svc.MeasureBaseURLs(context.Background(), MeasureUpstreamBaseURLsInput{
		Platform: PlatformDeepseek,
		APIKey:   "sk-x",
		BaseURLs: []string{"https://a.example.com", "https://b.example.com"},
	})
	require.NoError(t, err)
	require.Len(t, results, 2)
	require.Equal(t, "https://a.example.com", results[0].URL)
	require.True(t, results[0].OK)
}

func TestMeasureBaseURLsWorksWithoutAPIKey(t *testing.T) {
	// 基础连通性检查不需要 key：草稿行只靠出口（直连或代理）即可测。
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, nil, &stubUpstreamProber{})
	results, err := svc.MeasureBaseURLs(context.Background(), MeasureUpstreamBaseURLsInput{
		BaseURLs: []string{"https://a.example.com"},
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.True(t, results[0].OK)
}

func TestCreateEntryUsesPrimaryVendorBaseURL(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:    99,
		Name:  "A",
		Extra: map[string]any{UpstreamVendorIDExtraKey: int64(2)},
	}}
	repo := &stubUpstreamVendorRepo{
		vendors: []UpstreamVendor{{ID: 2, Name: "A"}},
		entries: map[int64][]UpstreamEntry{2: {{AccountID: 99, VendorID: 2, Name: "A"}}},
	}
	svc := NewUpstreamHubService(repo, admin, &stubUpstreamProber{})

	_, err := svc.CreateEntry(context.Background(), 2, CreateUpstreamEntryInput{
		Name:        "A",
		Platform:    PlatformDeepseek,
		BaseURL:     "https://slow.example.com",
		APIKey:      "sk-x",
		Concurrency: 1,
		Priority:    50,
		BaseURLs:    []string{"https://fast.example.com", "https://slow.example.com"},
	})
	require.NoError(t, err)
	require.NotNil(t, admin.created)
	require.Equal(t, "https://fast.example.com", admin.created.Credentials["base_url"])
}

func TestCreateEntryUsesSingleVendorBaseURL(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:    99,
		Name:  "A",
		Extra: map[string]any{UpstreamVendorIDExtraKey: int64(2)},
	}}
	repo := &stubUpstreamVendorRepo{
		vendors: []UpstreamVendor{{ID: 2, Name: "A"}},
		entries: map[int64][]UpstreamEntry{2: {{AccountID: 99, VendorID: 2, Name: "A"}}},
	}
	svc := NewUpstreamHubService(repo, admin, &stubUpstreamProber{})

	_, err := svc.CreateEntry(context.Background(), 2, CreateUpstreamEntryInput{
		Name:        "A",
		Platform:    PlatformDeepseek,
		APIKey:      "sk-x",
		Concurrency: 1,
		Priority:    50,
		BaseURLs:    []string{"https://only.example.com"},
	})
	require.NoError(t, err)
	require.Equal(t, "https://only.example.com", admin.created.Credentials["base_url"])
}

func TestMeasureBaseURLsUsesTestOnlyProxy(t *testing.T) {
	// 代理仅用于本次测速，不落库；这里断言它确实被带进了测量账号。
	admin := &stubUpstreamAdmin{}
	prober := &stubUpstreamProber{}
	svc := NewUpstreamHubService(&stubUpstreamVendorRepo{}, admin, prober)
	proxyID := int64(3)

	_, err := svc.MeasureBaseURLs(context.Background(), MeasureUpstreamBaseURLsInput{
		BaseURLs: []string{"https://a.example.com"},
		ProxyID:  &proxyID,
	})
	require.NoError(t, err)
	require.NotNil(t, prober.lastAccount)
	require.NotNil(t, prober.lastAccount.Proxy)
	require.Equal(t, proxyID, prober.lastAccount.Proxy.ID)
}

func TestCreateEntryUsesPrimaryEvenWhenUnmeasured(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:    99,
		Name:  "A",
		Extra: map[string]any{UpstreamVendorIDExtraKey: int64(2)},
	}}
	repo := &stubUpstreamVendorRepo{
		vendors: []UpstreamVendor{{ID: 2, Name: "A"}},
		entries: map[int64][]UpstreamEntry{2: {{AccountID: 99, VendorID: 2, Name: "A"}}},
	}
	svc := NewUpstreamHubService(repo, admin, &stubUpstreamProber{measureFail: true})

	_, err := svc.CreateEntry(context.Background(), 2, CreateUpstreamEntryInput{
		Name:        "A",
		Platform:    PlatformDeepseek,
		APIKey:      "sk-x",
		Concurrency: 1,
		Priority:    50,
		BaseURLs:    []string{"https://a.example.com", "https://b.example.com"},
	})
	require.NoError(t, err)
	require.NotNil(t, admin.created)
	require.Equal(t, "https://a.example.com", admin.created.Credentials["base_url"])
}

func TestUpdateEntryUsesPrimaryVendorBaseURL(t *testing.T) {
	admin := &stubUpstreamAdmin{account: &Account{
		ID:          5,
		Name:        "A",
		Platform:    PlatformDeepseek,
		Extra:       map[string]any{UpstreamVendorIDExtraKey: int64(2)},
		Credentials: map[string]any{"api_key": "sk-x", "base_url": "https://old.example.com"},
	}}
	repo := &stubUpstreamVendorRepo{entries: map[int64][]UpstreamEntry{2: {{AccountID: 5, VendorID: 2, Name: "A"}}}}
	svc := NewUpstreamHubService(repo, admin, &stubUpstreamProber{})

	_, err := svc.UpdateEntry(context.Background(), 5, UpdateUpstreamEntryInput{
		Name:     "A",
		BaseURL:  "",
		BaseURLs: []string{"https://fast.example.com", "https://slow.example.com"},
	})
	require.NoError(t, err)
	require.NotNil(t, admin.updated)
	require.Equal(t, "https://fast.example.com", admin.updated.Credentials["base_url"])
}
