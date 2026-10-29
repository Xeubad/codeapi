package admin

import (
	"errors"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// UpstreamHubHandler 上游聚合管理：厂家台账 + 透传条目 + 分组引用开关。
type UpstreamHubHandler struct {
	svc *service.UpstreamHubService
}

func NewUpstreamHubHandler(svc *service.UpstreamHubService) *UpstreamHubHandler {
	return &UpstreamHubHandler{svc: svc}
}

type upstreamVendorRequest struct {
	Name              string   `json:"name"`
	ConcurrencyBudget int      `json:"concurrency_budget"`
	Notes             *string  `json:"notes"`
	BaseURLs          []string `json:"base_urls"`
}

type upstreamEntryRequest struct {
	Name        string   `json:"name"`
	Platform    string   `json:"platform"`
	BaseURL     string   `json:"base_url"`
	APIKey      string   `json:"api_key"`
	Concurrency *int     `json:"concurrency"`
	Priority    *int     `json:"priority"`
	Schedulable *bool    `json:"schedulable"`
	BaseURLs    []string `json:"base_urls"`
}

type upstreamEntryGroupEnabledRequest struct {
	Enabled bool `json:"enabled"`
}

type upstreamEntryGroupsRequest struct {
	GroupIDs []int64 `json:"group_ids"`
}

type upstreamImportEntryRequest struct {
	AccountID int64 `json:"account_id"`
}

// ImportEntry 把账号管理里已有的转发账号纳入指定厂家。
// POST /api/v1/admin/upstream-vendors/:id/entries/import
func (h *UpstreamHubHandler) ImportEntry(c *gin.Context) {
	vendorID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamImportEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if req.AccountID <= 0 {
		response.BadRequest(c, "account_id is required")
		return
	}
	entry, err := h.svc.ImportAccountToVendor(c.Request.Context(), req.AccountID, vendorID)
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, entry)
}

func (h *UpstreamHubHandler) ListVendors(c *gin.Context) {
	vendors, err := h.svc.ListVendors(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, vendors)
}

func (h *UpstreamHubHandler) ListAllEntries(c *gin.Context) {
	entries, err := h.svc.ListAllEntries(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, entries)
}

func (h *UpstreamHubHandler) GetVendor(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	detail, err := h.svc.GetVendorDetail(c.Request.Context(), id)
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, detail)
}

func (h *UpstreamHubHandler) CreateVendor(c *gin.Context) {
	var req upstreamVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	vendor, err := h.svc.CreateVendor(c.Request.Context(), service.CreateUpstreamVendorInput{
		Name:              req.Name,
		ConcurrencyBudget: req.ConcurrencyBudget,
		Notes:             req.Notes,
		BaseURLs:          req.BaseURLs,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	c.JSON(http.StatusCreated, vendor)
}

func (h *UpstreamHubHandler) UpdateVendor(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	vendor, err := h.svc.UpdateVendor(c.Request.Context(), id, service.CreateUpstreamVendorInput{
		Name:              req.Name,
		ConcurrencyBudget: req.ConcurrencyBudget,
		Notes:             req.Notes,
		BaseURLs:          req.BaseURLs,
	})
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, vendor)
}

func (h *UpstreamHubHandler) DeleteVendor(c *gin.Context) {
	id, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteVendor(c.Request.Context(), id); err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *UpstreamHubHandler) CreateEntry(c *gin.Context) {
	vendorID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	concurrency := 1
	if req.Concurrency != nil {
		concurrency = *req.Concurrency
	}
	priority := 50
	if req.Priority != nil {
		priority = *req.Priority
	}
	entry, err := h.svc.CreateEntry(c.Request.Context(), vendorID, service.CreateUpstreamEntryInput{
		Name:        req.Name,
		Platform:    req.Platform,
		BaseURL:     req.BaseURL,
		APIKey:      req.APIKey,
		Concurrency: concurrency,
		Priority:    priority,
		Schedulable: req.Schedulable,
		BaseURLs:    req.BaseURLs,
	})
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *UpstreamHubHandler) UpdateEntry(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	entry, err := h.svc.UpdateEntry(c.Request.Context(), accountID, service.UpdateUpstreamEntryInput{
		Name:        req.Name,
		BaseURL:     req.BaseURL,
		APIKey:      req.APIKey,
		Concurrency: req.Concurrency,
		Priority:    req.Priority,
		Schedulable: req.Schedulable,
		BaseURLs:    req.BaseURLs,
	})
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, entry)
}

type upstreamEntryDuplicateRequest struct {
	Name        string `json:"name"`
	Platform    string `json:"platform"`
	Concurrency *int   `json:"concurrency"`
	Priority    *int   `json:"priority"`
}

func (h *UpstreamHubHandler) DuplicateEntry(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamEntryDuplicateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	concurrency := 1
	if req.Concurrency != nil {
		concurrency = *req.Concurrency
	}
	priority := 50
	if req.Priority != nil {
		priority = *req.Priority
	}
	entry, err := h.svc.DuplicateEntry(c.Request.Context(), accountID, service.DuplicateUpstreamEntryInput{
		Name:        req.Name,
		Platform:    req.Platform,
		Concurrency: concurrency,
		Priority:    priority,
	})
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *UpstreamHubHandler) DeleteEntry(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteEntry(c.Request.Context(), accountID); err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "deleted"})
}

func (h *UpstreamHubHandler) SetEntrySchedulable(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req struct {
		Schedulable bool `json:"schedulable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	entry, err := h.svc.SetEntrySchedulable(c.Request.Context(), accountID, req.Schedulable)
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, entry)
}

func (h *UpstreamHubHandler) SetEntryGroups(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamEntryGroupsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	groups, err := h.svc.SetEntryGroups(c.Request.Context(), accountID, req.GroupIDs)
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, groups)
}

func (h *UpstreamHubHandler) SetEntryGroupEnabled(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	groupID, ok := parseIDParam(c, "groupId")
	if !ok {
		return
	}
	var req upstreamEntryGroupEnabledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.svc.SetEntryGroupEnabled(c.Request.Context(), accountID, groupID, req.Enabled); err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, gin.H{"message": "updated"})
}

type upstreamEntryModelsRequest struct {
	Models []string `json:"models"`
}

// ProbeEntryModels 探测条目上游的模型列表（只读，不落库）。
// POST /api/v1/admin/upstream-entries/:id/probe-models
func (h *UpstreamHubHandler) ProbeEntryModels(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	probe, err := h.svc.ProbeEntryModels(c.Request.Context(), accountID)
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, probe)
}

// SetEntryModels 把勾选的模型写成条目的 model_mapping（提交列表即最终白名单）。
// PUT /api/v1/admin/upstream-entries/:id/models
func (h *UpstreamHubHandler) SetEntryModels(c *gin.Context) {
	accountID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var req upstreamEntryModelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	entry, err := h.svc.SetEntryModels(c.Request.Context(), accountID, req.Models)
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, entry)
}

type upstreamBaseURLLatencyRequest struct {
	AccountID int64    `json:"account_id"`
	Platform  string   `json:"platform"`
	APIKey    string   `json:"api_key"`
	BaseURLs  []string `json:"base_urls"`
	// ProxyID 是本次测速的出口（来自 IP 管理的已有代理），仅用于该请求。
	ProxyID *int64 `json:"proxy_id"`
}

// MeasureBaseURLs 并发测候选接入地址的延迟（只读），供保存前挑选最快者。
// POST /api/v1/admin/upstream-base-urls/latency
func (h *UpstreamHubHandler) MeasureBaseURLs(c *gin.Context) {
	var req upstreamBaseURLLatencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	results, err := h.svc.MeasureBaseURLs(c.Request.Context(), service.MeasureUpstreamBaseURLsInput{
		AccountID: req.AccountID,
		Platform:  req.Platform,
		APIKey:    req.APIKey,
		BaseURLs:  req.BaseURLs,
		ProxyID:   req.ProxyID,
	})
	if err != nil {
		writeUpstreamHubError(c, err)
		return
	}
	response.Success(c, results)
}

func (h *UpstreamHubHandler) ListGroups(c *gin.Context) {
	groups, err := h.svc.ListGroups(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, groups)
}

func writeUpstreamHubError(c *gin.Context, err error) {
	var syncErr *service.UpstreamModelSyncError
	if errors.As(err, &syncErr) {
		switch syncErr.Kind {
		case service.UpstreamModelSyncErrorConfiguration, service.UpstreamModelSyncErrorUnsupported:
			response.BadRequest(c, syncErr.SafeMessage())
		case service.UpstreamModelSyncErrorInternal:
			response.InternalError(c, syncErr.SafeMessage())
		default:
			response.Error(c, http.StatusBadGateway, syncErr.SafeMessage())
		}
		return
	}
	switch {
	case errors.Is(err, service.ErrUpstreamVendorNotFound),
		errors.Is(err, service.ErrUpstreamEntryNotFound):
		response.NotFound(c, err.Error())
	case errors.Is(err, service.ErrUpstreamAccountIneligible):
		response.BadRequest(c, err.Error())
	case errors.Is(err, service.ErrUpstreamAccountAlreadyOwned),
		errors.Is(err, service.ErrUpstreamVendorNotEmpty):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		response.ErrorFrom(c, err)
	}
}
