package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// Resource tuning endpoints let administrators persist pool overrides that the
// next process start applies. They deliberately do not mutate the running
// pools: the Redis client and DB pool are built at startup, and an online
// resize would silently diverge from the persisted values. The UI surfaces
// "pending restart" state until the effective values match the managed file.

type opsResourceTuningLimits struct {
	DBMaxOpenConns int `json:"db_max_open_conns"`
	RedisPoolSize  int `json:"redis_pool_size"`
}

type opsResourceTuningStatus struct {
	Running          config.ManagedPoolValues  `json:"running"`
	Baseline         config.ManagedPoolValues  `json:"baseline"`
	Managed          *config.ManagedPoolValues `json:"managed"`
	ManagedUpdatedAt string                    `json:"managed_updated_at,omitempty"`
	// none: no managed file. pending: file saved, differs from running values.
	// active: file saved and running values match it. invalid: file unreadable
	// or rejected at startup/now — baseline is in effect; save or remove to fix.
	ManagedState    string                  `json:"managed_state"`
	ManagedError    string                  `json:"managed_error,omitempty"`
	RestartRequired bool                    `json:"restart_required"`
	FilePath        string                  `json:"file_path"`
	Limits          opsResourceTuningLimits `json:"limits"`
}

func (h *OpsHandler) resourceTuningStatus() (*opsResourceTuningStatus, error) {
	if h.cfg == nil {
		return nil, errors.New("config unavailable")
	}
	running := config.CurrentPoolValues(h.cfg)
	baseline := h.cfg.PoolBaseline
	// Hand-built configs (tests, odd boot paths) never populated the baseline.
	if baseline == (config.ManagedPoolValues{}) {
		baseline = running
	}

	st := &opsResourceTuningStatus{
		Running:  running,
		Baseline: baseline,
		FilePath: config.ManagedPoolPath(),
		Limits: opsResourceTuningLimits{
			DBMaxOpenConns: config.ManagedPoolMaxDBOpenConns,
			RedisPoolSize:  config.ManagedPoolMaxRedisPoolSize,
		},
	}

	f, err := config.ReadManagedPool()
	switch {
	case err != nil:
		st.ManagedState = "invalid"
		st.ManagedError = err.Error()
		// Effective values are the baseline — the file was (or will be) ignored.
		st.RestartRequired = running != baseline
	case f == nil:
		st.ManagedState = "none"
		// If a file was applied this boot but removed afterwards, the override is
		// still live until restart; deleting only takes effect on next start.
		st.RestartRequired = h.cfg.ManagedPoolApplied && running != baseline
	default:
		v := f.Values
		st.Managed = &v
		st.ManagedUpdatedAt = f.UpdatedAt
		if h.cfg.ManagedPoolApplied && running == v {
			st.ManagedState = "active"
		} else {
			st.ManagedState = "pending"
		}
		st.RestartRequired = running != v
	}
	if h.cfg.ManagedPoolIgnored != "" && st.ManagedState != "invalid" {
		st.ManagedState = "invalid"
		st.ManagedError = h.cfg.ManagedPoolIgnored
	}
	return st, nil
}

// GetResourceTuning reports effective/baseline/managed pool values.
// GET /api/v1/admin/ops/resource-tuning
func (h *OpsHandler) GetResourceTuning(c *gin.Context) {
	if h.opsService == nil || h.cfg == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
	}
	st, err := h.resourceTuningStatus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, st)
}

type updateResourceTuningRequest struct {
	DBMaxOpenConns    *int `json:"db_max_open_conns" binding:"required"`
	DBMaxIdleConns    *int `json:"db_max_idle_conns" binding:"required"`
	RedisPoolSize     *int `json:"redis_pool_size" binding:"required"`
	RedisMinIdleConns *int `json:"redis_min_idle_conns" binding:"required"`
}

// UpdateResourceTuning saves the managed pool override (applies on restart).
// PUT /api/v1/admin/ops/resource-tuning
func (h *OpsHandler) UpdateResourceTuning(c *gin.Context) {
	if h.cfg == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
	}
	var req updateResourceTuningRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	values := config.ManagedPoolValues{
		DBMaxOpenConns:    *req.DBMaxOpenConns,
		DBMaxIdleConns:    *req.DBMaxIdleConns,
		RedisPoolSize:     *req.RedisPoolSize,
		RedisMinIdleConns: *req.RedisMinIdleConns,
	}
	if err := values.Validate(); err != nil {
		response.BadRequest(c, fmt.Sprintf("invalid pool values: %v", err))
		return
	}
	if err := config.WriteManagedPool(values); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("save managed pool override: %v", err))
		return
	}
	st, err := h.resourceTuningStatus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, st)
}

// DeleteResourceTuning removes the managed file; baseline applies after restart.
// DELETE /api/v1/admin/ops/resource-tuning
func (h *OpsHandler) DeleteResourceTuning(c *gin.Context) {
	if h.cfg == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
	}
	if err := config.RemoveManagedPool(); err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("remove managed pool override: %v", err))
		return
	}
	st, err := h.resourceTuningStatus()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, st)
}
