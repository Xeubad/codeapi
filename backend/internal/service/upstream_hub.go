package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/domain"
)

// UpstreamVendorIDExtraKey 是 accounts.extra 中记录所属上游厂家的键。
// 存在该键的账号由「上游聚合」管理，其它账号由「账号管理」管理。
const UpstreamVendorIDExtraKey = "upstream_vendor_id"

// 上游聚合领域类型。
type UpstreamVendor struct {
	ID                   int64     `json:"id"`
	Name                 string    `json:"name"`
	ConcurrencyBudget    int       `json:"concurrency_budget"`
	Notes                string    `json:"notes"`
	EntryCount           int       `json:"entry_count"`
	AllocatedConcurrency int       `json:"allocated_concurrency"`
	ReferencedGroupCount int       `json:"referenced_group_count"`
	OverBudget           bool      `json:"over_budget"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type UpstreamEntryGroupRef struct {
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	Platform  string `json:"platform"`
	Priority  int    `json:"priority"`
	Enabled   bool   `json:"enabled"`
}

type UpstreamEntry struct {
	AccountID   int64                   `json:"account_id"`
	VendorID    int64                   `json:"vendor_id"`
	Name        string                  `json:"name"`
	Platform    string                  `json:"platform"`
	BaseURL     string                  `json:"base_url"`
	HasAPIKey   bool                    `json:"has_api_key"`
	Concurrency int                     `json:"concurrency"`
	Priority    int                     `json:"priority"`
	Schedulable bool                    `json:"schedulable"`
	Status      string                  `json:"status"`
	Groups      []UpstreamEntryGroupRef `json:"groups"`
	// Models 是 credentials.model_mapping 的键（已登记的对外模型名，排序后）。
	// 网关的 /v1/models 与 composite 归属解析都只认这份映射，为空即无模型可暴露。
	Models []string `json:"models"`
}

// UpstreamEntryModelProbe 是探测上游模型列表的结果（只读，不落库）。
type UpstreamEntryModelProbe struct {
	Models   []string `json:"models"`
	Existing []string `json:"existing"`
}

type UpstreamGroupOption struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

type UpstreamVendorDetail struct {
	Vendor   UpstreamVendor  `json:"vendor"`
	Entries  []UpstreamEntry `json:"entries"`
	Warnings []string        `json:"warnings"`
}

type CreateUpstreamVendorInput struct {
	Name              string
	ConcurrencyBudget int
	Notes             *string
}

type CreateUpstreamEntryInput struct {
	Name        string
	Platform    string
	BaseURL     string
	APIKey      string
	Concurrency int
	Priority    int
	Schedulable *bool
}

type UpdateUpstreamEntryInput struct {
	Name        string
	BaseURL     string
	APIKey      string // 空表示保持不变
	Concurrency *int
	Priority    *int
	Schedulable *bool
}

var (
	ErrUpstreamVendorNotFound = errors.New("upstream vendor not found")
	ErrUpstreamVendorNotEmpty = errors.New("upstream vendor still has entries")
	ErrUpstreamEntryNotFound  = errors.New("upstream entry not found")
	ErrUpstreamGroupInvalid   = errors.New("group platform does not match entry platform")
)

// UpstreamVendorRepository 上游厂家与条目（accounts.extra 归属）的持久化接口。
type UpstreamVendorRepository interface {
	List(ctx context.Context) ([]UpstreamVendor, error)
	GetByID(ctx context.Context, id int64) (*UpstreamVendor, error)
	Create(ctx context.Context, name string, budget int, notes *string) (*UpstreamVendor, error)
	Update(ctx context.Context, id int64, name string, budget int, notes *string) (*UpstreamVendor, error)
	SoftDelete(ctx context.Context, id int64) error
	ListEntries(ctx context.Context, vendorID int64) ([]UpstreamEntry, error)
	ListGroups(ctx context.Context) ([]UpstreamGroupOption, error)
	SetEntryGroupEnabled(ctx context.Context, accountID, groupID int64, enabled bool) error
	SetEntryGroups(ctx context.Context, accountID int64, groupIDs []int64) error
}

// upstreamModelProber 是「读取上游真实模型列表」的能力，由 *AccountTestService 实现。
// 抽成接口只为本包单测可替换，生产装配仍是 AccountTestService。
type upstreamModelProber interface {
	FetchUpstreamSupportedModels(ctx context.Context, account *Account) ([]string, error)
}

// UpstreamHubService 上游聚合：厂家台账 + 透传条目编排。
// 条目本身是 type=apikey 的 Account（自定义 base_url 直连上游），由 AdminService
// 负责既有校验与生命周期；归属只认 extra.upstream_vendor_id。
// 之所以不用 type=upstream：网关各转发路径按 account.type 分支取凭据，upstream
// 只有 antigravity relay 有实现，其余平台会直接报 unsupported account type。
type UpstreamHubService struct {
	repo  UpstreamVendorRepository
	admin AdminService
	probe upstreamModelProber
}

func NewUpstreamHubService(
	repo UpstreamVendorRepository,
	admin AdminService,
	probe upstreamModelProber,
) *UpstreamHubService {
	return &UpstreamHubService{repo: repo, admin: admin, probe: probe}
}

func (s *UpstreamHubService) ListVendors(ctx context.Context) ([]UpstreamVendor, error) {
	vendors, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for i := range vendors {
		entries, err := s.repo.ListEntries(ctx, vendors[i].ID)
		if err != nil {
			return nil, err
		}
		enrichUpstreamVendor(&vendors[i], entries)
	}
	return vendors, nil
}

// ListAllEntries 返回全部厂家的条目，供聚合调度页做跨厂家编排。
func (s *UpstreamHubService) ListAllEntries(ctx context.Context) ([]UpstreamEntry, error) {
	vendors, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	entries := make([]UpstreamEntry, 0)
	for _, vendor := range vendors {
		vendorEntries, err := s.repo.ListEntries(ctx, vendor.ID)
		if err != nil {
			return nil, err
		}
		entries = append(entries, vendorEntries...)
	}
	return entries, nil
}

func (s *UpstreamHubService) GetVendorDetail(ctx context.Context, id int64) (*UpstreamVendorDetail, error) {
	vendor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	entries, err := s.repo.ListEntries(ctx, id)
	if err != nil {
		return nil, err
	}
	enrichUpstreamVendor(vendor, entries)
	detail := &UpstreamVendorDetail{Vendor: *vendor, Entries: entries, Warnings: []string{}}
	if vendor.ConcurrencyBudget > 0 && vendor.AllocatedConcurrency > vendor.ConcurrencyBudget {
		detail.Warnings = append(detail.Warnings, fmt.Sprintf(
			"各条目并发合计 %d 已超过厂家预算 %d，请调低条目并发或上调预算",
			vendor.AllocatedConcurrency, vendor.ConcurrencyBudget))
	}
	return detail, nil
}

func (s *UpstreamHubService) CreateVendor(ctx context.Context, input CreateUpstreamVendorInput) (*UpstreamVendor, error) {
	name, err := normalizeVendorName(input.Name)
	if err != nil {
		return nil, err
	}
	if input.ConcurrencyBudget < 0 {
		return nil, fmt.Errorf("concurrency_budget must be >= 0")
	}
	return s.repo.Create(ctx, name, input.ConcurrencyBudget, trimNotes(input.Notes))
}

func (s *UpstreamHubService) UpdateVendor(ctx context.Context, id int64, input CreateUpstreamVendorInput) (*UpstreamVendor, error) {
	name, err := normalizeVendorName(input.Name)
	if err != nil {
		return nil, err
	}
	if input.ConcurrencyBudget < 0 {
		return nil, fmt.Errorf("concurrency_budget must be >= 0")
	}
	return s.repo.Update(ctx, id, name, input.ConcurrencyBudget, trimNotes(input.Notes))
}

func (s *UpstreamHubService) DeleteVendor(ctx context.Context, id int64) error {
	entries, err := s.repo.ListEntries(ctx, id)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return ErrUpstreamVendorNotEmpty
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *UpstreamHubService) CreateEntry(ctx context.Context, vendorID int64, input CreateUpstreamEntryInput) (*UpstreamEntry, error) {
	if _, err := s.repo.GetByID(ctx, vendorID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("entry name is required")
	}
	if !domain.IsConcretePlatform(strings.TrimSpace(input.Platform)) {
		return nil, fmt.Errorf("platform %q is not a concrete platform", input.Platform)
	}
	baseURL, err := normalizeUpstreamBaseURL(input.BaseURL)
	if err != nil {
		return nil, err
	}
	apiKey := strings.TrimSpace(input.APIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}
	if input.Concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be >= 1")
	}
	schedulable := true
	if input.Schedulable != nil {
		schedulable = *input.Schedulable
	}
	account, err := s.admin.CreateAccount(ctx, &CreateAccountInput{
		Name:     name,
		Platform: strings.TrimSpace(input.Platform),
		Type:     AccountTypeAPIKey,
		Credentials: map[string]any{
			"api_key":  apiKey,
			"base_url": baseURL,
		},
		Extra: map[string]any{
			UpstreamVendorIDExtraKey: vendorID,
		},
		Concurrency:          input.Concurrency,
		Priority:             input.Priority,
		SkipDefaultGroupBind: true,
	})
	if err != nil {
		return nil, err
	}
	if !schedulable {
		if _, err := s.admin.SetAccountSchedulable(ctx, account.ID, false); err != nil {
			return nil, err
		}
	}
	return s.getEntry(ctx, account.ID)
}

func (s *UpstreamHubService) UpdateEntry(ctx context.Context, accountID int64, input UpdateUpstreamEntryInput) (*UpstreamEntry, error) {
	current, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if _, ok := upstreamVendorIDFromExtra(current.Extra); !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("entry name is required")
	}
	baseURL, err := normalizeUpstreamBaseURL(input.BaseURL)
	if err != nil {
		return nil, err
	}
	// UpdateAccount 的合并语义是「非敏感键完全由 incoming 决定」
	// （MergePreservingSensitiveCreds），所以必须以现有 credentials 为底，
	// 否则探测写入的 model_mapping 会在任意一次编辑后被抹掉。
	credentials := make(map[string]any, len(current.Credentials)+2)
	for key, value := range current.Credentials {
		credentials[key] = value
	}
	credentials["api_key"] = current.GetCredential("api_key")
	credentials["base_url"] = baseURL
	if apiKey := strings.TrimSpace(input.APIKey); apiKey != "" {
		credentials["api_key"] = apiKey
	}
	update := &UpdateAccountInput{
		Name:                  name,
		Type:                  AccountTypeAPIKey,
		Credentials:           credentials,
		Concurrency:           input.Concurrency,
		Priority:              input.Priority,
		SkipMixedChannelCheck: true,
	}
	if _, err := s.admin.UpdateAccount(ctx, accountID, update); err != nil {
		return nil, err
	}
	if input.Schedulable != nil {
		if _, err := s.admin.SetAccountSchedulable(ctx, accountID, *input.Schedulable); err != nil {
			return nil, err
		}
	}
	return s.getEntry(ctx, accountID)
}

// DuplicateEntry 以现有条目为模板新建条目（同厂家）：服务端复制凭据，
// API Key 不出后端；用于同 key 跨类别或同类别多条目场景。
type DuplicateUpstreamEntryInput struct {
	Name        string
	Platform    string
	Concurrency int
	Priority    int
}

func (s *UpstreamHubService) DuplicateEntry(ctx context.Context, accountID int64, input DuplicateUpstreamEntryInput) (*UpstreamEntry, error) {
	current, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	vendorID, ok := upstreamVendorIDFromExtra(current.Extra)
	if !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("entry name is required")
	}
	platform := strings.TrimSpace(input.Platform)
	if platform == "" {
		platform = current.Platform
	}
	if !domain.IsConcretePlatform(platform) {
		return nil, fmt.Errorf("platform %q is not a concrete platform", platform)
	}
	if input.Concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be >= 1")
	}
	apiKey := current.GetCredential("api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("source entry has no api key")
	}
	// 复制条目同时带走已登记的模型白名单，否则副本会没有任何可调用模型。
	credentials := map[string]any{
		"api_key":  apiKey,
		"base_url": current.GetCredential("base_url"),
	}
	if mapping, ok := current.Credentials["model_mapping"]; ok && mapping != nil {
		credentials["model_mapping"] = mapping
	}
	account, err := s.admin.CreateAccount(ctx, &CreateAccountInput{
		Name:        name,
		Platform:    platform,
		Type:        AccountTypeAPIKey,
		Credentials: credentials,
		Extra: map[string]any{
			UpstreamVendorIDExtraKey: vendorID,
		},
		Concurrency:          input.Concurrency,
		Priority:             input.Priority,
		SkipDefaultGroupBind: true,
	})
	if err != nil {
		return nil, err
	}
	return s.getEntry(ctx, account.ID)
}

func (s *UpstreamHubService) DeleteEntry(ctx context.Context, accountID int64) error {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}
	if _, ok := upstreamVendorIDFromExtra(account.Extra); !ok {
		return ErrUpstreamEntryNotFound
	}
	return s.admin.DeleteAccount(ctx, accountID)
}

func (s *UpstreamHubService) SetEntrySchedulable(ctx context.Context, accountID int64, schedulable bool) (*UpstreamEntry, error) {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if _, ok := upstreamVendorIDFromExtra(account.Extra); !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	if _, err := s.admin.SetAccountSchedulable(ctx, accountID, schedulable); err != nil {
		return nil, err
	}
	return s.getEntry(ctx, accountID)
}

func (s *UpstreamHubService) SetEntryGroupEnabled(ctx context.Context, accountID, groupID int64, enabled bool) error {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}
	if _, ok := upstreamVendorIDFromExtra(account.Extra); !ok {
		return ErrUpstreamEntryNotFound
	}
	return s.repo.SetEntryGroupEnabled(ctx, accountID, groupID, enabled)
}

// upstreamGroupAllowsPlatform 判断条目能否挂进该分组：具体平台分组要求平台相等，
// composite 综合分组可接收任意具体平台条目（对齐平台自身 canCopyAccountsFromGroupPlatform）。
func upstreamGroupAllowsPlatform(groupPlatform, entryPlatform string) bool {
	if groupPlatform == PlatformComposite {
		return true
	}
	return groupPlatform == entryPlatform
}

func (s *UpstreamHubService) SetEntryGroups(ctx context.Context, accountID int64, groupIDs []int64) ([]UpstreamEntryGroupRef, error) {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if _, ok := upstreamVendorIDFromExtra(account.Extra); !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	groups, err := s.repo.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	platformByID := make(map[int64]string, len(groups))
	for _, g := range groups {
		platformByID[g.ID] = g.Platform
	}
	for _, groupID := range groupIDs {
		platform, ok := platformByID[groupID]
		if !ok {
			return nil, fmt.Errorf("group %d not found", groupID)
		}
		// 综合分组（composite）天生跨平台，可接收任意具体平台条目；
		// 具体平台分组仍要求与条目平台一致，口径对齐平台自身的
		// canCopyAccountsFromGroupPlatform（admin_group.go）。
		if !upstreamGroupAllowsPlatform(platform, account.Platform) {
			return nil, ErrUpstreamGroupInvalid
		}
	}
	if err := s.repo.SetEntryGroups(ctx, accountID, groupIDs); err != nil {
		return nil, err
	}
	entry, err := s.getEntry(ctx, accountID)
	if err != nil {
		return nil, err
	}
	return entry.Groups, nil
}

func (s *UpstreamHubService) ListGroups(ctx context.Context) ([]UpstreamGroupOption, error) {
	return s.repo.ListGroups(ctx)
}

// ProbeEntryModels 读取条目上游的模型列表（只读，不落库）。
// 复用 AccountTestService 的探测链，覆盖 openai / 国产多协议供应商 / anthropic /
// gemini / grok；失败原因由 UpstreamModelSyncError 分类，handler 据此映射状态码。
func (s *UpstreamHubService) ProbeEntryModels(ctx context.Context, accountID int64) (*UpstreamEntryModelProbe, error) {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if _, ok := upstreamVendorIDFromExtra(account.Extra); !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	if s.probe == nil {
		return nil, errors.New("upstream model prober is not configured")
	}
	models, err := s.probe.FetchUpstreamSupportedModels(ctx, account)
	if err != nil {
		return nil, err
	}
	return &UpstreamEntryModelProbe{
		Models:   normalizeUpstreamModelIDs(models),
		Existing: storedUpstreamEntryModels(account.Credentials),
	}, nil
}

// SetEntryModels 把勾选的模型写成条目的 model_mapping（from=to）。
// 提交列表即最终白名单；空列表表示清空白名单（模型列表回落到平台默认表）。
func (s *UpstreamHubService) SetEntryModels(ctx context.Context, accountID int64, models []string) (*UpstreamEntry, error) {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	if _, ok := upstreamVendorIDFromExtra(account.Extra); !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	selected := normalizeUpstreamModelIDs(models)
	if len(selected) > maxUpstreamEntryModels {
		return nil, fmt.Errorf("too many models: %d exceeds limit %d", len(selected), maxUpstreamEntryModels)
	}
	mapping := make(map[string]any, len(selected))
	for _, model := range selected {
		mapping[model] = model
	}
	credentials := make(map[string]any, len(account.Credentials)+1)
	for key, value := range account.Credentials {
		credentials[key] = value
	}
	credentials["model_mapping"] = mapping
	if _, err := s.admin.UpdateAccount(ctx, accountID, &UpdateAccountInput{
		Name:                  account.Name,
		Type:                  AccountTypeAPIKey,
		Credentials:           credentials,
		SkipMixedChannelCheck: true,
	}); err != nil {
		return nil, err
	}
	return s.getEntry(ctx, accountID)
}

// maxUpstreamEntryModels 限制单条目可登记的模型数，避免误提交超大白名单。
const maxUpstreamEntryModels = 500

func normalizeUpstreamModelIDs(models []string) []string {
	seen := make(map[string]struct{}, len(models))
	out := make([]string, 0, len(models))
	for _, model := range models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		if _, ok := seen[model]; ok {
			continue
		}
		seen[model] = struct{}{}
		out = append(out, model)
	}
	sort.Strings(out)
	return out
}

// storedUpstreamEntryModels 只读 credentials 里真实存储的 model_mapping 键，
// 不走 Account.GetModelMapping（后者对部分平台会补默认映射，不是管理员登记的口径）。
func storedUpstreamEntryModels(credentials map[string]any) []string {
	stored := stringMappingFromRaw(credentials["model_mapping"])
	keys := make([]string, 0, len(stored))
	for key, mapped := range stored {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(mapped) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (s *UpstreamHubService) getEntry(ctx context.Context, accountID int64) (*UpstreamEntry, error) {
	account, err := s.admin.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}
	vendorID, ok := upstreamVendorIDFromExtra(account.Extra)
	if !ok {
		return nil, ErrUpstreamEntryNotFound
	}
	entries, err := s.repo.ListEntries(ctx, vendorID)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.AccountID == accountID {
			return &entry, nil
		}
	}
	return nil, ErrUpstreamEntryNotFound
}

func enrichUpstreamVendor(vendor *UpstreamVendor, entries []UpstreamEntry) {
	vendor.EntryCount = len(entries)
	allocated := 0
	groupSeen := make(map[int64]struct{})
	for _, entry := range entries {
		allocated += entry.Concurrency
		for _, ref := range entry.Groups {
			groupSeen[ref.GroupID] = struct{}{}
		}
	}
	vendor.AllocatedConcurrency = allocated
	vendor.ReferencedGroupCount = len(groupSeen)
	vendor.OverBudget = vendor.ConcurrencyBudget > 0 && allocated > vendor.ConcurrencyBudget
}

func normalizeVendorName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", fmt.Errorf("vendor name is required")
	}
	if len([]rune(trimmed)) > 100 {
		return "", fmt.Errorf("vendor name is too long")
	}
	return trimmed, nil
}

func normalizeUpstreamBaseURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("base url is required")
	}
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("base url is invalid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("base url scheme must be http or https")
	}
	return strings.TrimRight(trimmed, "/"), nil
}

func trimNotes(notes *string) *string {
	if notes == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*notes)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// upstreamVendorIDFromExtra 解析 accounts.extra 中的厂家归属（JSONB 数字读回为 float64）。
func upstreamVendorIDFromExtra(extra map[string]any) (int64, bool) {
	if extra == nil {
		return 0, false
	}
	value, ok := extra[UpstreamVendorIDExtraKey]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case int64:
		return typed, typed > 0
	case int:
		return int64(typed), typed > 0
	case float64:
		id := int64(typed)
		return id, id > 0
	case string:
		var id int64
		if _, err := fmt.Sscanf(typed, "%d", &id); err == nil && id > 0 {
			return id, true
		}
	}
	return 0, false
}
