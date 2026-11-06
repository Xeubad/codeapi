package service

import (
	"context"
	"errors"
	"sort"
	"time"
)

// ops_tuning_advisor.go — observation-driven resource tuning advisor.
//
// Design contract (docs/RESOURCE_TUNING_DEVELOPMENT_CN.md):
//   - The advisor never mutates pool settings. It produces candidate values
//     plus evidence; the administrator confirms and a restart applies.
//   - Sustained evidence is required: a single spike never generates advice.
//   - Missing, stale or contradictory data yields "insufficient" verdicts
//     instead of guesses — hardware size alone never produces a suggestion.

const (
	// opsAdviceMinSamples is the minimum recent samples needed before any
	// candidate is emitted.
	opsAdviceMinSamples = 5
	// opsAdviceWindowSamples is how many recent samples are examined.
	opsAdviceWindowSamples = 10
	// opsAdviceMaxSampleAge bounds how old the newest sample may be before the
	// whole window is flagged stale.
	opsAdviceMaxSampleAge = 10 * time.Minute

	// opsAdviceNearLimitPct: usage ≥ this fraction of the effective cap for
	// every sample in the window → propose increasing the cap.
	opsAdviceNearLimitPct = 0.85
	// opsAdviceFarBelowPct: usage < this fraction of the cap for every sample
	// → propose decreasing the cap.
	opsAdviceFarBelowPct = 0.20
	// opsAdviceWaitTriggerPct: a per-window waiting delta ≥ this fraction of
	// the cap counts as a wait event.
	opsAdviceWaitTriggerPct = 0.10
	// opsAdviceGrowStepPct / opsAdviceShrinkStepPct bound each suggestion to a
	// modest single step; repeated confirmations can iterate further.
	opsAdviceGrowStepPct   = 0.25
	opsAdviceShrinkStepPct = 0.50
	// opsAdviceDBIdleFloorPct: idle suggestions never drop below this fraction
	// of max-open so reconnect storms cannot be reintroduced by advice.
	opsAdviceDBIdleFloorPct = 0.10
	// opsAdviceMinRedisIdle keeps at least a small idle base by default.
	opsAdviceMinRedisIdle = 8
)

// OpsTuningAdviceAction is the recommendation for one tunable value.
type OpsTuningAdviceAction string

const (
	OpsAdviceKeep         OpsTuningAdviceAction = "keep"         // evidence says leave it
	OpsAdviceIncrease     OpsTuningAdviceAction = "increase"     // candidate is higher
	OpsAdviceDecrease     OpsTuningAdviceAction = "decrease"     // candidate is lower
	OpsAdviceUnchanged    OpsTuningAdviceAction = "unchanged"    // candidate equals current
	OpsAdviceInsufficient OpsTuningAdviceAction = "insufficient" // not enough trustworthy data
)

// OpsTuningPoolValues mirrors the four managed pool knobs. Kept as a service
// type so the advisor does not import config.
type OpsTuningPoolValues struct {
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	RedisPoolSize     int
	RedisMinIdleConns int
}

// OpsTuningAdviceItem is one parameter's recommendation with evidence.
type OpsTuningAdviceItem struct {
	Key       string                `json:"key"`
	Current   int                   `json:"current"`
	Candidate int                   `json:"candidate"`
	Action    OpsTuningAdviceAction `json:"action"`
	Reason    string                `json:"reason"`
	Evidence  map[string]any        `json:"evidence,omitempty"`
}

// OpsTuningAdvice is the full advisor output for the ops UI.
type OpsTuningAdvice struct {
	GeneratedAt  time.Time             `json:"generated_at"`
	SampleCount  int                   `json:"sample_count"`
	OldestSample string                `json:"oldest_sample,omitempty"`
	NewestSample string                `json:"newest_sample,omitempty"`
	Stalled      bool                  `json:"stalled"` // pending restart — advice shown but flagged
	Items        []OpsTuningAdviceItem `json:"items"`
	Warnings     []string              `json:"warnings,omitempty"`
}

// OpsResourceTuningAdvisor computes advice from a window of system metrics.
type OpsResourceTuningAdvisor struct {
	now func() time.Time // injectable for deterministic tests
}

// NewOpsResourceTuningAdvisor returns an advisor using the wall clock.
func NewOpsResourceTuningAdvisor() *OpsResourceTuningAdvisor {
	return &OpsResourceTuningAdvisor{now: time.Now}
}

// ListSystemMetrics returns recent global one-minute system metric snapshots,
// newest first. Thin pass-through so the handler stays out of the repository.
func (s *OpsService) ListSystemMetrics(ctx context.Context, limit int) ([]OpsSystemMetricsSnapshot, error) {
	if s == nil || s.opsRepo == nil {
		return nil, errors.New("ops repository unavailable")
	}
	return s.opsRepo.ListSystemMetrics(ctx, 1, limit)
}

// Advise evaluates the metrics window against the running pool values.
// pendingManaged is the saved-but-not-yet-active override (nil when none) and
// is only used to flag a stalled state — the advisor never writes it.
func (a *OpsResourceTuningAdvisor) Advise(
	samples []OpsSystemMetricsSnapshot,
	current OpsTuningPoolValues,
	pendingManaged *OpsTuningPoolValues,
) OpsTuningAdvice {
	out := OpsTuningAdvice{
		GeneratedAt: a.now().UTC(),
		Items:       make([]OpsTuningAdviceItem, 0, 4),
	}
	if pendingManaged != nil {
		out.Stalled = true
		out.Warnings = append(out.Warnings, "pending_restart")
	}

	usables := a.usableSamples(samples)
	out.SampleCount = len(usables)
	if len(usables) > 0 {
		out.NewestSample = usables[0].CreatedAt.UTC().Format(time.RFC3339)
		out.OldestSample = usables[len(usables)-1].CreatedAt.UTC().Format(time.RFC3339)
		if age := a.now().UTC().Sub(usables[0].CreatedAt.UTC()); age > opsAdviceMaxSampleAge {
			out.Warnings = append(out.Warnings, "stale_samples")
		}
	}

	if len(usables) < opsAdviceMinSamples {
		out.Warnings = append(out.Warnings, "insufficient_samples")
		out.Items = append(out.Items,
			a.insufficient("db_max_open_conns", current.DBMaxOpenConns, "need ≥5 recent samples"),
			a.insufficient("db_max_idle_conns", current.DBMaxIdleConns, "need ≥5 recent samples"),
			a.insufficient("redis_pool_size", current.RedisPoolSize, "need ≥5 recent samples"),
			a.insufficient("redis_min_idle_conns", current.RedisMinIdleConns, "need ≥5 recent samples"),
		)
		return out
	}

	if unhealthy := a.unhealthyReasons(usables); len(unhealthy) > 0 {
		out.Warnings = append(out.Warnings, unhealthy...)
		out.Items = append(out.Items,
			a.insufficient("db_max_open_conns", current.DBMaxOpenConns, "dependency unhealthy"),
			a.insufficient("db_max_idle_conns", current.DBMaxIdleConns, "dependency unhealthy"),
			a.insufficient("redis_pool_size", current.RedisPoolSize, "dependency unhealthy"),
			a.insufficient("redis_min_idle_conns", current.RedisMinIdleConns, "dependency unhealthy"),
		)
		return out
	}

	out.Items = append(out.Items,
		a.dbOpenAdvice(usables, current),
		a.dbIdleAdvice(usables, current),
		a.redisPoolAdvice(usables, current),
		a.redisIdleAdvice(usables, current),
	)
	return out
}

// usableSamples returns the newest ≤ opsAdviceWindowSamples samples with a
// valid timestamp, newest first.
func (a *OpsResourceTuningAdvisor) usableSamples(in []OpsSystemMetricsSnapshot) []OpsSystemMetricsSnapshot {
	out := make([]OpsSystemMetricsSnapshot, 0, len(in))
	for _, m := range in {
		if !m.CreatedAt.IsZero() {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	if len(out) > opsAdviceWindowSamples {
		out = out[:opsAdviceWindowSamples]
	}
	return out
}

func (a *OpsResourceTuningAdvisor) unhealthyReasons(samples []OpsSystemMetricsSnapshot) []string {
	var reasons []string
	dbDown, redisDown := 0, 0
	for i := range samples {
		if samples[i].DBOK != nil && !*samples[i].DBOK {
			dbDown++
		}
		if samples[i].RedisOK != nil && !*samples[i].RedisOK {
			redisDown++
		}
	}
	// Only flag when the dependency was unhealthy in a majority of the window;
	// a transient blip should not suppress advice.
	if dbDown > len(samples)/2 {
		reasons = append(reasons, "db_unhealthy")
	}
	if redisDown > len(samples)/2 {
		reasons = append(reasons, "redis_unhealthy")
	}
	return reasons
}

func (a *OpsResourceTuningAdvisor) insufficient(key string, current int, why string) OpsTuningAdviceItem {
	return OpsTuningAdviceItem{
		Key: key, Current: current, Candidate: current,
		Action: OpsAdviceInsufficient, Reason: why,
	}
}

func dbPoolUsed(m *OpsSystemMetricsSnapshot) (int, bool) {
	if m.DBConnActive == nil || m.DBConnIdle == nil {
		return 0, false
	}
	return *m.DBConnActive + *m.DBConnIdle, true
}

func optInt(v *int) (int, bool) {
	if v == nil {
		return 0, false
	}
	return *v, true
}

func (a *OpsResourceTuningAdvisor) dbOpenAdvice(samples []OpsSystemMetricsSnapshot, cur OpsTuningPoolValues) OpsTuningAdviceItem {
	const key = "db_max_open_conns"
	capv := cur.DBMaxOpenConns
	if capv <= 0 {
		return a.insufficient(key, capv, "current value unreadable")
	}
	nearLimit, farBelow, waits, usable, maxUsed := 0, 0, 0, 0, 0
	for i := range samples {
		used, ok := dbPoolUsed(&samples[i])
		if !ok {
			continue
		}
		usable++
		if used > maxUsed {
			maxUsed = used
		}
		pct := float64(used) / float64(capv)
		if pct >= opsAdviceNearLimitPct {
			nearLimit++
		}
		if pct < opsAdviceFarBelowPct {
			farBelow++
		}
		if w, ok := optInt(samples[i].DBConnWaiting); ok && float64(w) >= float64(capv)*opsAdviceWaitTriggerPct {
			waits++
		}
	}
	ev := map[string]any{
		"samples": usable, "near_limit": nearLimit, "far_below": farBelow,
		"wait_events": waits, "max_used": maxUsed, "cap": capv,
	}
	if usable < opsAdviceMinSamples {
		return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: capv, Action: OpsAdviceInsufficient, Reason: "db connection stats missing", Evidence: ev}
	}
	if waits > 0 || nearLimit == usable {
		cand := adviceGrow(capv, opsAdviceGrowStepPct, maxUsed+2)
		return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: cand, Action: adviceActionFor(cand, capv), Reason: "sustained near-limit usage or wait events", Evidence: ev}
	}
	if farBelow == usable {
		floor := max(maxUsed+2, 4)
		cand := adviceShrink(capv, opsAdviceShrinkStepPct, floor)
		return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: cand, Action: adviceActionFor(cand, capv), Reason: "sustained low usage", Evidence: ev}
	}
	return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: capv, Action: OpsAdviceKeep, Reason: "usage within band", Evidence: ev}
}

func (a *OpsResourceTuningAdvisor) dbIdleAdvice(samples []OpsSystemMetricsSnapshot, cur OpsTuningPoolValues) OpsTuningAdviceItem {
	const key = "db_max_idle_conns"
	capOpen, curIdle := cur.DBMaxOpenConns, cur.DBMaxIdleConns
	if capOpen <= 0 || curIdle < 0 || curIdle > capOpen {
		return a.insufficient(key, curIdle, "invalid current values")
	}
	floor := max(int(float64(capOpen)*opsAdviceDBIdleFloorPct), 1)
	idleHeavy, usable, maxUsed, minUsed := 0, 0, 0, -1
	for i := range samples {
		used, ok := dbPoolUsed(&samples[i])
		idle, ok2 := optInt(samples[i].DBConnIdle)
		if !ok || !ok2 {
			continue
		}
		usable++
		if used > maxUsed {
			maxUsed = used
		}
		if minUsed < 0 || used < minUsed {
			minUsed = used
		}
		if curIdle > 0 && float64(idle) > float64(curIdle)*0.9 {
			idleHeavy++
		}
	}
	ev := map[string]any{"samples": usable, "idle_saturated": idleHeavy, "min_used": minUsed, "max_used": maxUsed, "floor": floor}
	if usable < opsAdviceMinSamples {
		return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: curIdle, Action: OpsAdviceInsufficient, Reason: "db connection stats missing", Evidence: ev}
	}
	// Shrink idle only when the observed working set never approaches the cap
	// and the idle reserve stays saturated — evidence it is never drained.
	persistentHeadroom := minUsed >= 0 && float64(minUsed) < float64(capOpen)*opsAdviceFarBelowPct
	if idleHeavy == usable && persistentHeadroom {
		cand := adviceShrink(curIdle, opsAdviceShrinkStepPct, floor)
		return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: cand, Action: adviceActionFor(cand, curIdle), Reason: "idle reserve never used", Evidence: ev}
	}
	// Grow idle toward the observed working set when it sits above the
	// current idle — comfort increase capped at the open limit.
	if maxUsed > curIdle && curIdle < capOpen {
		cand := min(maxUsed, capOpen)
		return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: cand, Action: adviceActionFor(cand, curIdle), Reason: "idle below observed working set", Evidence: ev}
	}
	return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: curIdle, Action: OpsAdviceKeep, Reason: "idle appropriate", Evidence: ev}
}

func (a *OpsResourceTuningAdvisor) redisPoolAdvice(samples []OpsSystemMetricsSnapshot, cur OpsTuningPoolValues) OpsTuningAdviceItem {
	const key = "redis_pool_size"
	capv := cur.RedisPoolSize
	if capv <= 0 {
		return a.insufficient(key, capv, "current value unreadable")
	}
	nearLimit, farBelow, usable, maxTotal := 0, 0, 0, 0
	for i := range samples {
		total, ok := optInt(samples[i].RedisConnTotal)
		if !ok {
			continue
		}
		usable++
		if total > maxTotal {
			maxTotal = total
		}
		pct := float64(total) / float64(capv)
		if pct >= opsAdviceNearLimitPct {
			nearLimit++
		}
		if pct < opsAdviceFarBelowPct {
			farBelow++
		}
	}
	ev := map[string]any{"samples": usable, "near_limit": nearLimit, "far_below": farBelow, "max_total": maxTotal, "cap": capv}
	if usable < opsAdviceMinSamples {
		return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: capv, Action: OpsAdviceInsufficient, Reason: "redis stats missing", Evidence: ev}
	}
	if nearLimit == usable {
		cand := adviceGrow(capv, opsAdviceGrowStepPct, maxTotal+2)
		return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: cand, Action: adviceActionFor(cand, capv), Reason: "pool near baseline every sample (soft cap)", Evidence: ev}
	}
	if farBelow == usable {
		floor := max(maxTotal+2, opsAdviceMinRedisIdle)
		cand := adviceShrink(capv, opsAdviceShrinkStepPct, floor)
		return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: cand, Action: adviceActionFor(cand, capv), Reason: "pool far below baseline every sample", Evidence: ev}
	}
	return OpsTuningAdviceItem{Key: key, Current: capv, Candidate: capv, Action: OpsAdviceKeep, Reason: "usage within band", Evidence: ev}
}

func (a *OpsResourceTuningAdvisor) redisIdleAdvice(samples []OpsSystemMetricsSnapshot, cur OpsTuningPoolValues) OpsTuningAdviceItem {
	const key = "redis_min_idle_conns"
	pool, curIdle := cur.RedisPoolSize, cur.RedisMinIdleConns
	if pool <= 0 || curIdle < 0 || curIdle > pool {
		return a.insufficient(key, curIdle, "invalid current values")
	}
	floor := min(opsAdviceMinRedisIdle, pool)
	idleHeavy, usable, maxTotal := 0, 0, 0
	for i := range samples {
		idle, ok := optInt(samples[i].RedisConnIdle)
		if !ok {
			continue
		}
		usable++
		if curIdle > 0 && float64(idle) > float64(curIdle)*0.9 {
			idleHeavy++
		}
		if t, ok2 := optInt(samples[i].RedisConnTotal); ok2 && t > maxTotal {
			maxTotal = t
		}
	}
	ev := map[string]any{"samples": usable, "idle_saturated": idleHeavy, "max_total": maxTotal, "floor": floor}
	if usable < opsAdviceMinSamples {
		return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: curIdle, Action: OpsAdviceInsufficient, Reason: "redis stats missing", Evidence: ev}
	}
	if idleHeavy == usable && float64(maxTotal) < float64(pool)*opsAdviceFarBelowPct {
		cand := adviceShrink(curIdle, opsAdviceShrinkStepPct, floor)
		return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: cand, Action: adviceActionFor(cand, curIdle), Reason: "idle reserve unused under sustained low load", Evidence: ev}
	}
	return OpsTuningAdviceItem{Key: key, Current: curIdle, Candidate: curIdle, Action: OpsAdviceKeep, Reason: "idle appropriate", Evidence: ev}
}

func adviceActionFor(candidate, current int) OpsTuningAdviceAction {
	switch {
	case candidate > current:
		return OpsAdviceIncrease
	case candidate < current:
		return OpsAdviceDecrease
	default:
		return OpsAdviceUnchanged
	}
}

// adviceGrow grows base by stepPct, at least to floor and never below base+1.
func adviceGrow(base int, stepPct float64, floor int) int {
	step := int(float64(base) * stepPct)
	if step < 1 {
		step = 1
	}
	cand := base + step
	if floor > cand {
		cand = floor
	}
	return cand
}

// adviceShrink shrinks base by stepPct toward floor (never below floor or 1).
func adviceShrink(base int, stepPct float64, floor int) int {
	if base <= floor {
		return base
	}
	step := int(float64(base) * stepPct)
	if step < 1 {
		step = 1
	}
	cand := base - step
	if cand < floor {
		cand = floor
	}
	if cand < 1 {
		cand = 1
	}
	return cand
}
