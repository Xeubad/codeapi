<template>
  <section class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-dark-700 dark:bg-dark-800" aria-labelledby="resource-analysis-title">
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div>
        <h2 id="resource-analysis-title" class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.ops.resourceAnalysis.title') }}</h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.description') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <span :class="['rounded-full px-3 py-1 text-xs font-medium', stateBadgeClass]">{{ stateLabel }}</span>
        <button
          type="button"
          class="rounded-lg bg-primary-600 px-3 py-1.5 text-xs font-medium text-white transition-colors hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="tuningLoading || !tuning"
          @click="openEdit"
        >
          {{ t('admin.ops.resourceAnalysis.manage') }}
        </button>
      </div>
    </div>

    <p v-if="tuningError" class="mt-3 rounded-lg bg-red-50 p-3 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">{{ tuningError }}</p>
    <p v-if="tuning?.managed_state === 'invalid'" class="mt-3 rounded-lg bg-red-50 p-3 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
      {{ t('admin.ops.resourceAnalysis.managedInvalid', { error: tuning.managed_error || 'unknown', path: tuning.file_path }) }}
    </p>
    <p v-else-if="tuning?.restart_required" class="mt-3 rounded-lg bg-amber-50 p-3 text-sm text-amber-700 dark:bg-amber-900/20 dark:text-amber-400">
      {{ t('admin.ops.resourceAnalysis.restartRequired') }}
    </p>

    <p v-if="!metrics" class="mt-4 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.noData') }}</p>
    <template v-else>
      <p class="mt-4 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.sampleTime', { time: sampleTime }) }}</p>
      <p v-if="stale" class="mt-2 text-sm text-amber-700 dark:text-amber-400">{{ t('admin.ops.resourceAnalysis.stale') }}</p>
      <div class="mt-4 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <div v-for="item in readings" :key="item.label" class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ item.label }}</div>
          <div class="mt-2 text-lg font-semibold tabular-nums text-gray-900 dark:text-white">{{ item.value }}</div>
          <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ item.detail }}</div>
        </div>
      </div>
      <!-- 自动优化建议：只建议不应用，需管理员确认后重启生效 -->
      <div v-if="adviceItems.length" class="mt-4 rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-900">
        <div class="flex flex-wrap items-center justify-between gap-2">
          <div class="text-xs font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400">
            {{ t('admin.ops.resourceAnalysis.adviceTitle') }}
            <span class="ml-2 font-normal normal-case text-gray-400">{{ t('admin.ops.resourceAnalysis.adviceSamples', { count: tuning?.advice?.sample_count ?? 0 }) }}</span>
          </div>
          <button
            v-if="hasApplyable"
            type="button"
            class="rounded-lg border border-primary-500 px-3 py-1 text-xs font-medium text-primary-600 transition-colors hover:bg-primary-50 disabled:opacity-50 dark:border-primary-400 dark:text-primary-400 dark:hover:bg-dark-800"
            @click="applyAdvice"
          >
            {{ t('admin.ops.resourceAnalysis.applyAdvice') }}
          </button>
        </div>
        <ul class="mt-3 space-y-2">
          <li v-for="it in adviceItems" :key="it.key" class="flex flex-wrap items-baseline justify-between gap-2 text-sm">
            <span class="font-medium text-gray-700 dark:text-gray-300">{{ fieldLabel(it.key) }}</span>
            <span class="tabular-nums text-gray-600 dark:text-gray-400">
              {{ it.current }} <span class="text-gray-400">→</span>
              <span :class="candidateClass(it)">{{ it.candidate }}</span>
              <span class="ml-2 text-xs" :class="candidateClass(it)">{{ actionLabel(it) }}</span>
            </span>
          </li>
        </ul>
        <p v-if="adviceWarnings.length" class="mt-3 text-xs text-amber-700 dark:text-amber-400">
          {{ adviceWarnings.join(' · ') }}
        </p>
        <p class="mt-1 text-xs text-gray-400 dark:text-gray-500">{{ t('admin.ops.resourceAnalysis.adviceHint') }}</p>
      </div>

      <ul class="mt-4 space-y-2 text-sm text-gray-600 dark:text-gray-300">
        <li v-if="metrics.db_ok === false" class="text-red-600 dark:text-red-400">{{ t('admin.ops.resourceAnalysis.dbUnavailable') }}</li>
        <li v-if="metrics.redis_ok === false" class="text-red-600 dark:text-red-400">{{ t('admin.ops.resourceAnalysis.redisUnavailable') }}</li>
        <li>{{ t('admin.ops.resourceAnalysis.redisScope') }}</li>
        <li>{{ t('admin.ops.resourceAnalysis.limits') }}</li>
      </ul>
    </template>

    <BaseDialog :show="editOpen" :title="t('admin.ops.resourceAnalysis.editTitle')" width="wide" @close="editOpen = false">
      <div class="space-y-4">
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('admin.ops.resourceAnalysis.editHint') }}</p>
        <div class="grid gap-4 sm:grid-cols-2">
          <div v-for="field in fields" :key="field.key" class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
            <label :for="`pool-${field.key}`" class="block text-sm font-medium text-gray-700 dark:text-gray-300">{{ field.label }}</label>
            <input
              :id="`pool-${field.key}`"
              v-model.number="form[field.key]"
              type="number"
              min="0"
              step="1"
              class="mt-2 w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm tabular-nums text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-800 dark:text-white"
            />
            <p class="mt-1.5 text-xs text-gray-500 dark:text-gray-400">{{ field.hint }}</p>
          </div>
        </div>
        <p v-if="formError" class="text-sm text-red-600 dark:text-red-400">{{ formError }}</p>
        <div class="flex flex-wrap items-center justify-between gap-3">
          <button
            v-if="tuning?.managed"
            type="button"
            class="text-sm font-medium text-red-600 hover:text-red-700 disabled:opacity-50 dark:text-red-400"
            :disabled="saving"
            @click="restoreBaseline"
          >
            {{ t('admin.ops.resourceAnalysis.restoreBaseline') }}
          </button>
          <span v-else></span>
          <div class="flex gap-2">
            <button
              type="button"
              class="rounded-lg px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 disabled:opacity-50 dark:text-gray-300 dark:hover:bg-dark-700"
              :disabled="saving"
              @click="editOpen = false"
            >
              {{ t('common.cancel') }}
            </button>
            <button
              type="button"
              class="rounded-lg bg-primary-600 px-4 py-2 text-sm font-medium text-white hover:bg-primary-700 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="saving || !formValid"
              @click="save"
            >
              {{ saving ? t('admin.ops.resourceAnalysis.saving') : t('admin.ops.resourceAnalysis.save') }}
            </button>
          </div>
        </div>
      </div>
    </BaseDialog>
    <TotpStepUpDialog :controller="stepUp" />
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import TotpStepUpDialog from '@/components/auth/TotpStepUpDialog.vue'
import { isStepUpCancelled, useStepUp } from '@/composables/useStepUp'
import { useAppStore } from '@/stores'
import {
  opsAPI,
  type OpsPoolValues,
  type OpsResourceTuningStatus,
  type OpsSystemMetricsSnapshot
} from '@/api/admin/ops'

const props = defineProps<{ metrics: OpsSystemMetricsSnapshot | null | undefined; refreshToken?: number }>()
const { t } = useI18n()
const appStore = useAppStore()
const stepUp = useStepUp()

const tuning = ref<OpsResourceTuningStatus | null>(null)
const tuningLoading = ref(false)
const tuningError = ref('')
const editOpen = ref(false)
const saving = ref(false)
const formError = ref('')

const form = ref<OpsPoolValues>({ db_max_open_conns: 0, db_max_idle_conns: 0, redis_pool_size: 0, redis_min_idle_conns: 0 })

async function refreshTuning() {
  tuningLoading.value = true
  try {
    tuning.value = await opsAPI.getResourceTuning()
    tuningError.value = ''
  } catch (err: any) {
    tuningError.value = err?.message || t('admin.ops.resourceAnalysis.loadFailed')
  } finally {
    tuningLoading.value = false
  }
}

onMounted(refreshTuning)
watch(() => props.refreshToken, () => { void refreshTuning() })

const stateLabel = computed(() => {
  switch (tuning.value?.managed_state) {
    case 'active': return t('admin.ops.resourceAnalysis.stateActive')
    case 'pending': return t('admin.ops.resourceAnalysis.statePending')
    case 'invalid': return t('admin.ops.resourceAnalysis.stateInvalid')
    default: return t('admin.ops.resourceAnalysis.observation')
  }
})
const stateBadgeClass = computed(() => {
  switch (tuning.value?.managed_state) {
    case 'active': return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
    case 'pending': return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400'
    case 'invalid': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
    default: return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'
  }
})

const fields = computed(() => {
  const lim = tuning.value?.limits
  return [
    { key: 'db_max_open_conns' as const, label: t('admin.ops.resourceAnalysis.fieldDbOpen'), hint: t('admin.ops.resourceAnalysis.hintDbOpen', { max: lim?.db_max_open_conns ?? 0 }) },
    { key: 'db_max_idle_conns' as const, label: t('admin.ops.resourceAnalysis.fieldDbIdle'), hint: t('admin.ops.resourceAnalysis.hintDbIdle') },
    { key: 'redis_pool_size' as const, label: t('admin.ops.resourceAnalysis.fieldRedisPool'), hint: t('admin.ops.resourceAnalysis.hintRedisPool', { max: lim?.redis_pool_size ?? 0 }) },
    { key: 'redis_min_idle_conns' as const, label: t('admin.ops.resourceAnalysis.fieldRedisIdle'), hint: t('admin.ops.resourceAnalysis.hintRedisIdle') },
  ]
})

// ---- 自动优化建议（只读，不自动应用） ----
const adviceItems = computed(() => tuning.value?.advice?.items ?? [])
const hasApplyable = computed(() => adviceItems.value.some(it => it.action === 'increase' || it.action === 'decrease'))
const adviceWarnings = computed(() => (tuning.value?.advice?.warnings ?? []).map(w => t(`admin.ops.resourceAnalysis.warn_${w}`, w)))

function fieldLabel(key: string): string {
  const f = fields.value.find(x => x.key === key)
  return f ? f.label : key
}

function actionLabel(it: { action: string }): string {
  return t(`admin.ops.resourceAnalysis.action_${it.action}`, it.action)
}

function candidateClass(it: { action: string }): string {
  switch (it.action) {
    case 'increase': return 'text-emerald-600 dark:text-emerald-400'
    case 'decrease': return 'text-amber-600 dark:text-amber-400'
    case 'insufficient': return 'text-gray-400'
    default: return 'text-gray-600 dark:text-gray-300'
  }
}

// 一键填充：把建议候选写进表单，仍需人工确认 + step-up 后保存
function applyAdvice() {
  if (!tuning.value?.advice) return
  const next = { ...form.value }
  for (const it of tuning.value.advice.items) {
    if (it.action === 'increase' || it.action === 'decrease') {
      next[it.key] = it.candidate
    }
  }
  form.value = next
  formError.value = ''
  editOpen.value = true
}

const formValid = computed(() => {
  const f = form.value
  const lim = tuning.value?.limits
  return Number.isInteger(f.db_max_open_conns) && f.db_max_open_conns >= 1 && (!lim || f.db_max_open_conns <= lim.db_max_open_conns)
    && Number.isInteger(f.db_max_idle_conns) && f.db_max_idle_conns >= 0 && f.db_max_idle_conns <= f.db_max_open_conns
    && Number.isInteger(f.redis_pool_size) && f.redis_pool_size >= 1 && (!lim || f.redis_pool_size <= lim.redis_pool_size)
    && Number.isInteger(f.redis_min_idle_conns) && f.redis_min_idle_conns >= 0 && f.redis_min_idle_conns <= f.redis_pool_size
})

function openEdit() {
  if (!tuning.value) return
  const src = tuning.value.managed ?? tuning.value.running
  form.value = { ...src }
  formError.value = ''
  editOpen.value = true
}

async function save() {
  if (!formValid.value) {
    formError.value = t('admin.ops.resourceAnalysis.formInvalid')
    return
  }
  saving.value = true
  try {
    tuning.value = await stepUp.run(() => opsAPI.updateResourceTuning({ ...form.value }))
    editOpen.value = false
    appStore.showSuccess(t('admin.ops.resourceAnalysis.saved'))
  } catch (err: any) {
    if (isStepUpCancelled(err)) return
    formError.value = err?.message || t('admin.ops.resourceAnalysis.saveFailed')
  } finally {
    saving.value = false
  }
}

async function restoreBaseline() {
  saving.value = true
  try {
    tuning.value = await stepUp.run(() => opsAPI.deleteResourceTuning())
    editOpen.value = false
    appStore.showSuccess(t('admin.ops.resourceAnalysis.restored'))
  } catch (err: any) {
    if (isStepUpCancelled(err)) return
    formError.value = err?.message || t('admin.ops.resourceAnalysis.saveFailed')
  } finally {
    saving.value = false
  }
}

const sampleTime = computed(() => {
  const date = props.metrics?.created_at ? new Date(props.metrics.created_at) : null
  return date && Number.isFinite(date.getTime()) ? date.toLocaleString() : t('admin.ops.noData')
})
const stale = computed(() => {
  const time = props.metrics?.created_at ? Date.parse(props.metrics.created_at) : NaN
  return !Number.isFinite(time) || Date.now() - time > 10 * 60 * 1000
})
const percent = (value: number | null | undefined) => typeof value === 'number' && Number.isFinite(value) ? `${value.toFixed(1)}%` : '—'
const pendingSuffix = (key: keyof OpsPoolValues) =>
  tuning.value && tuning.value.managed_state === 'pending' && tuning.value.managed && tuning.value.managed[key] !== tuning.value.running[key]
    ? ` → ${tuning.value.managed[key]}`
    : ''
const readings = computed(() => {
  const m = props.metrics
  const db = m?.db_conn_active != null && m.db_conn_idle != null ? m.db_conn_active + m.db_conn_idle : null
  const count = (value: number | null | undefined) => value == null ? '—' : String(value)
  return [
    { label: 'CPU', value: percent(m?.cpu_usage_percent), detail: t('admin.ops.resourceAnalysis.cpuScope') },
    { label: t('admin.ops.memory'), value: percent(m?.memory_usage_percent), detail: m?.memory_used_mb != null && m.memory_total_mb != null ? `${m.memory_used_mb} / ${m.memory_total_mb} MiB` : t('admin.ops.noData') },
    { label: t('admin.ops.db'), value: `${count(db)} / ${count(m?.db_max_open_conns)}${pendingSuffix('db_max_open_conns')}`, detail: t('admin.ops.resourceAnalysis.dbScope') },
    { label: 'Redis', value: `${count(m?.redis_conn_total)} / ${count(m?.redis_pool_size)}${pendingSuffix('redis_pool_size')}`, detail: t('admin.ops.resourceAnalysis.redisBase') },
  ]
})
</script>
