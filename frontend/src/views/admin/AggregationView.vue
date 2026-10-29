<template>
  <AppLayout>
    <!-- 单元列表：与分组/渠道等列表页一致，使用 TablePageLayout + 工具栏 -->
    <TablePageLayout v-if="!selectedUnitId">
      <template #filters>
        <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-start">
          <div class="flex flex-1 flex-wrap items-center gap-3">
            <div class="relative w-full sm:w-64">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
              />
              <input
                v-model="unitSearch"
                type="text"
                :placeholder="t('admin.aggregation.searchUnits')"
                class="input pl-10"
              />
            </div>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.aggregation.unitsTotal', { count: filteredUnits.length }) }}
            </span>
          </div>
          <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 lg:w-auto">
            <button class="btn btn-secondary" :disabled="loading" @click="loadAll">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button class="btn btn-primary" @click="openCreateUnit">
              <Icon name="plus" size="md" class="mr-2" />
              {{ t('admin.aggregation.createUnit') }}
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('admin.aggregation.unitName') }}</th>
              <th>{{ t('admin.aggregation.entryCount') }}</th>
              <th>{{ t('admin.aggregation.routeCount') }}</th>
              <th>RPM</th>
              <th>{{ t('admin.aggregation.status') }}</th>
              <th class="text-right">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="unit in filteredUnits" :key="unit.id">
              <td>
                <div class="font-medium text-gray-900 dark:text-white">{{ unit.name }}</div>
                <div v-if="unit.description" class="max-w-[28rem] truncate text-xs text-gray-500">
                  {{ unit.description }}
                </div>
              </td>
              <td>{{ attachedCountFor(unit.id) }}</td>
              <td>{{ routesCountByGroup[unit.id] ?? '—' }}</td>
              <td>{{ unit.rpm_limit || '∞' }}</td>
              <td>
                <span class="badge" :class="unit.status === 'active' ? 'badge-primary' : 'badge-gray'">
                  {{ unit.status }}
                </span>
              </td>
              <td class="text-right">
                <button class="btn btn-sm btn-secondary" @click="openUnit(unit.id)">
                  {{ t('admin.aggregation.manage') }}
                </button>
              </td>
            </tr>
            <tr v-if="filteredUnits.length === 0">
              <td colspan="6">
                <EmptyState
                  :title="t('admin.aggregation.noUnits')"
                  :description="t('admin.aggregation.noUnitsHint')"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </template>
    </TablePageLayout>

    <!-- 单元详情：头部卡片 + 页签，与渠道监控等详情页一致 -->
    <div v-else class="w-full min-w-0 space-y-6 pb-8">
      <header
        class="page-header mb-0 rounded-3xl bg-white p-5 shadow-sm ring-1 ring-gray-900/5 dark:bg-dark-800 dark:ring-dark-700 sm:p-6"
      >
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h1 class="page-title flex items-center gap-2 text-xl font-black text-gray-900 dark:text-white">
              <button class="btn btn-ghost btn-sm" @click="backToList">← {{ t('admin.aggregation.back') }}</button>
              <span
                class="inline-flex h-8 w-8 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >
                <Icon name="swap" size="sm" />
              </span>
              {{ selectedUnit?.name }}
              <span class="badge badge-primary">{{ t('admin.aggregation.compositeBadge') }}</span>
            </h1>
            <p class="page-description mt-1.5 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.aggregation.unitDesc') }}
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-3">
            <button class="btn btn-secondary btn-sm" @click="goToGroups">
              {{ t('admin.aggregation.configureInGroups') }}
            </button>
            <button class="btn btn-secondary btn-sm" :disabled="loading" @click="loadAll">
              {{ t('admin.aggregation.refresh') }}
            </button>
          </div>
        </div>
        <div class="mt-4 border-t border-gray-100 pt-4 dark:border-dark-700">
          <div class="tabs inline-flex w-full max-w-xl flex-wrap sm:w-auto" role="tablist">
            <button
              type="button"
              role="tab"
              class="tab flex-1 sm:flex-none"
              :class="detailTab === 'supply' ? 'tab-active' : ''"
              :aria-selected="detailTab === 'supply'"
              @click="detailTab = 'supply'"
            >
              {{ t('admin.aggregation.tabSupply') }}
            </button>
            <button
              type="button"
              role="tab"
              class="tab flex-1 sm:flex-none"
              :class="detailTab === 'routes' ? 'tab-active' : ''"
              :aria-selected="detailTab === 'routes'"
              @click="detailTab = 'routes'"
            >
              {{ t('admin.aggregation.tabRoutes') }}
            </button>
            <button
              type="button"
              role="tab"
              class="tab flex-1 sm:flex-none"
              :class="detailTab === 'chain' ? 'tab-active' : ''"
              :aria-selected="detailTab === 'chain'"
              @click="detailTab = 'chain'"
            >
              {{ t('admin.aggregation.tabChain') }}
            </button>
          </div>
        </div>
      </header>

      <div class="grid gap-4 sm:grid-cols-3">
        <div class="card p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.statEntries') }}</p>
          <p class="mt-1 text-2xl font-semibold text-gray-900 dark:text-white">{{ attachedCount }}</p>
        </div>
        <div class="card p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.statRoutes') }}</p>
          <p class="mt-1 text-2xl font-semibold text-gray-900 dark:text-white">{{ routes.length }}</p>
        </div>
        <div class="card p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.statKeys') }}</p>
          <p class="mt-1 text-2xl font-semibold text-gray-900 dark:text-white">{{ stats?.active_api_keys || 0 }}</p>
        </div>
      </div>

      <!-- 供给汇入 -->
      <section v-if="detailTab === 'supply'" class="card p-5">
        <div class="mb-3">
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.aggregation.supply') }}</h2>
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.supplyHint') }}</p>
        </div>
        <div v-if="entries.length === 0" class="py-6 text-center text-sm text-gray-500">
          {{ t('admin.aggregation.noEntries') }}
        </div>
        <div v-else class="space-y-4">
          <div v-for="vendor in vendorGroups" :key="vendor.id">
            <div class="mb-1 flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
              <span class="font-medium text-gray-700 dark:text-gray-300">{{ vendor.name }}</span>
              <span>
                {{ t('admin.aggregation.vendorBudget', {
                  used: vendor.allocated_concurrency,
                  total: vendor.concurrency_budget || '∞'
                }) }}
              </span>
            </div>
            <div class="divide-y divide-gray-100 dark:divide-dark-700">
              <div
                v-for="entry in vendor.entries"
                :key="entry.account_id"
                class="flex items-center gap-3 py-2"
              >
                <input
                  type="checkbox"
                  class="cursor-pointer"
                  :checked="isAttached(entry)"
                  @change="toggleAttach(entry)"
                />
                <div class="min-w-0 flex-1">
                  <div class="truncate text-sm text-gray-900 dark:text-white">{{ entry.name }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">
                    {{ entry.platform }} · {{ t('admin.aggregation.modelsCount', { count: entry.models?.length || 0 }) }}
                  </div>
                </div>
                <span v-if="isAttached(entry) && !isEnabledHere(entry)" class="badge badge-gray">
                  {{ t('admin.aggregation.disabledHere') }}
                </span>
                <Toggle
                  v-if="isAttached(entry)"
                  :model-value="isEnabledHere(entry)"
                  @update:model-value="() => toggleEntryEnabled(entry)"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 归一模型 -->
      <section v-else-if="detailTab === 'routes'" class="card p-5">
        <div class="mb-3">
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.aggregation.modelRoutes') }}</h2>
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.modelRoutesHint') }}</p>
        </div>
        <div class="table-container">
          <table class="table">
            <thead>
              <tr>
                <th>{{ t('admin.aggregation.publicModel') }}</th>
                <th>{{ t('admin.aggregation.targetPlatform') }}</th>
                <th>{{ t('admin.aggregation.upstreamModel') }}</th>
                <th>{{ t('admin.aggregation.priority') }}</th>
                <th>{{ t('admin.aggregation.enabled') }}</th>
                <th class="text-right">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="route in routes" :key="route.id">
                <td class="font-mono text-xs">{{ route.public_model }}</td>
                <td><span class="badge badge-primary">{{ route.target_platform }}</span></td>
                <td class="font-mono text-xs">{{ route.upstream_model || route.public_model }}</td>
                <td>{{ route.priority }}</td>
                <td>
                  <Toggle :model-value="route.enabled" @update:model-value="() => toggleRouteEnabled(route)" />
                </td>
                <td class="text-right">
                  <button class="btn btn-sm btn-ghost text-red-600" @click="removeRoute(route)">
                    {{ t('admin.aggregation.delete') }}
                  </button>
                </td>
              </tr>
              <tr v-if="routes.length === 0">
                <td colspan="6" class="py-4 text-center text-xs text-gray-500">
                  {{ t('admin.aggregation.noRoutes') }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="mt-4 grid gap-3 sm:grid-cols-2">
          <input v-model="routeForm.public_model" class="input" :placeholder="t('admin.aggregation.publicModelPlaceholder')" />
          <Select v-model="routeForm.target_platform" :options="platformOptions" />
          <input v-model="routeForm.upstream_model" class="input" :placeholder="t('admin.aggregation.upstreamModelPlaceholder')" />
          <input v-model.number="routeForm.priority" class="input" type="number" min="1" :placeholder="t('admin.aggregation.priority')" />
        </div>
        <div class="mt-3 flex flex-wrap items-center justify-between gap-2">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.routeFormHint') }}</p>
          <button class="btn btn-primary btn-sm" :disabled="saving || !routeForm.public_model.trim()" @click="addRoute">
            {{ t('admin.aggregation.addRoute') }}
          </button>
        </div>
      </section>

      <!-- 链路视图 -->
      <section v-else class="card p-5">
        <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
          <div>
            <h2 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('admin.aggregation.chainView') }}</h2>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.chainHint') }}</p>
          </div>
          <div class="flex items-center gap-3 text-xs text-gray-500 dark:text-gray-400">
            <span class="inline-flex items-center gap-1">
              <span class="h-2 w-4 rounded bg-primary-500"></span>{{ t('admin.aggregation.legendOn') }}
            </span>
            <span class="inline-flex items-center gap-1">
              <span class="h-2 w-4 rounded border border-dashed border-gray-400"></span>{{ t('admin.aggregation.legendOff') }}
            </span>
          </div>
        </div>
        <div class="grid gap-3 md:grid-cols-4">
          <div class="space-y-2">
            <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.vendors') }}</h3>
            <div
              v-for="vendor in vendors"
              :key="vendor.id"
              class="rounded-lg border border-gray-200 px-2 py-1.5 text-sm dark:border-dark-600"
            >
              <div class="truncate font-medium text-gray-900 dark:text-white">{{ vendor.name }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ t('admin.aggregation.vendorBudget', { used: vendor.allocated_concurrency, total: vendor.concurrency_budget || '∞' }) }}
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.entries') }}</h3>
            <div
              v-for="entry in entries"
              :key="entry.account_id"
              class="rounded-lg border px-2 py-1.5 text-sm"
              :class="
                isAttached(entry) && isEnabledHere(entry)
                  ? 'border-primary-300 bg-primary-50 text-primary-800 dark:border-primary-700 dark:bg-primary-900/20 dark:text-primary-300'
                  : 'border-dashed border-gray-300 text-gray-500 dark:border-dark-600 dark:text-dark-300'
              "
            >
              <div class="truncate">{{ entry.name }}</div>
              <div class="text-xs opacity-80">
                {{ entry.platform }} ·
                {{ isAttached(entry) && !isEnabledHere(entry) ? t('admin.aggregation.disabledHere') : t('admin.aggregation.modelsCount', { count: entry.models?.length || 0 }) }}
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.group') }}</h3>
            <div
              v-if="selectedUnit"
              class="rounded-lg border border-primary-300 bg-primary-50 px-2 py-1.5 text-sm text-primary-800 dark:border-primary-700 dark:bg-primary-900/20 dark:text-primary-300"
            >
              <div class="truncate font-medium">{{ selectedUnit.name }}</div>
              <div class="text-xs opacity-80">
                {{ t('admin.aggregation.attachedCount', { count: attachedCount }) }} · RPM {{ selectedUnit.rpm_limit || '∞' }}
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <h3 class="text-xs font-semibold uppercase text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.customers') }}</h3>
            <div class="rounded-lg border border-gray-200 px-2 py-1.5 text-sm dark:border-dark-600">
              <div class="text-gray-900 dark:text-white">
                {{ t('admin.aggregation.activeKeys', { count: stats?.active_api_keys || 0 }) }}
              </div>
              <div class="text-xs text-gray-500 dark:text-gray-400">
                {{ t('admin.aggregation.totalKeys', { count: stats?.total_api_keys || 0 }) }}
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <BaseDialog :show="showCreateDialog" :title="t('admin.aggregation.createUnit')" @close="showCreateDialog = false">
      <div class="space-y-3">
        <div>
          <label class="input-label">{{ t('admin.aggregation.unitName') }}</label>
          <input v-model="createForm.name" class="input" type="text" :placeholder="t('admin.aggregation.unitNamePlaceholder')" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.aggregation.unitDescription') }}</label>
          <input
            v-model="createForm.description"
            class="input"
            type="text"
            :placeholder="t('admin.aggregation.unitDescriptionPlaceholder')"
          />
        </div>
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.aggregation.createUnitHint') }}</p>
      </div>
      <template #footer>
        <button class="btn btn-secondary" @click="showCreateDialog = false">{{ t('admin.aggregation.cancel') }}</button>
        <button class="btn btn-primary" :disabled="creating || !createForm.name.trim()" @click="createUnit">
          {{ t('admin.aggregation.createUnit') }}
        </button>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import upstreamHubAPI, { type UpstreamEntry, type UpstreamVendor } from '@/api/admin/upstreamHub'
import groupsAPI from '@/api/admin/groups'
import { CONCRETE_PLATFORM_OPTIONS } from '@/constants/platforms'
import type { AdminGroup, CompositeModelRoute } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Select from '@/components/common/Select.vue'
import Toggle from '@/components/common/Toggle.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

const loading = ref(false)
const saving = ref(false)
const creating = ref(false)
const groups = ref<AdminGroup[]>([])
const vendors = ref<UpstreamVendor[]>([])
const entries = ref<UpstreamEntry[]>([])
const routes = ref<CompositeModelRoute[]>([])
const routesCountByGroup = ref<Record<number, number>>({})
const stats = ref<{ total_api_keys: number; active_api_keys: number } | null>(null)
const selectedUnitId = ref<number | null>(null)
const detailTab = ref<'supply' | 'routes' | 'chain'>('supply')
const unitSearch = ref('')

const units = computed(() => groups.value.filter((group) => group.platform === 'composite'))
const filteredUnits = computed(() => {
  const keyword = unitSearch.value.trim().toLowerCase()
  if (!keyword) return units.value
  return units.value.filter(
    (unit) =>
      unit.name.toLowerCase().includes(keyword) ||
      (unit.description || '').toLowerCase().includes(keyword)
  )
})
const selectedUnit = computed(() => units.value.find((group) => group.id === selectedUnitId.value) || null)
const platformOptions = computed(() =>
  CONCRETE_PLATFORM_OPTIONS.map((platform) => ({ label: platform.label, value: platform.value }))
)
const vendorGroups = computed(() =>
  vendors.value.map((vendor) => ({
    ...vendor,
    entries: entries.value.filter((entry) => entry.vendor_id === vendor.id)
  }))
)
const attachedCount = computed(() => entries.value.filter((entry) => isAttached(entry)).length)

function attachedCountFor(groupId: number): number {
  return entries.value.filter((entry) => (entry.groups || []).some((group) => group.group_id === groupId)).length
}

function groupRef(entry: UpstreamEntry) {
  return entry.groups?.find((group) => group.group_id === selectedUnitId.value)
}

function isAttached(entry: UpstreamEntry): boolean {
  return !!groupRef(entry)
}

function isEnabledHere(entry: UpstreamEntry): boolean {
  return groupRef(entry)?.enabled ?? false
}

function errorMessage(err: unknown): string {
  return extractApiErrorMessage(err, t('admin.aggregation.loadFailed'))
}

async function loadUnits() {
  groups.value = await groupsAPI.getAllIncludingInactive()
  const counts: Record<number, number> = {}
  await Promise.all(
    units.value.map(async (unit) => {
      try {
        counts[unit.id] = (await groupsAPI.listCompositeRoutes(unit.id)).length
      } catch {
        counts[unit.id] = 0
      }
    })
  )
  routesCountByGroup.value = counts
}

async function loadUnitData() {
  if (!selectedUnitId.value) {
    routes.value = []
    stats.value = null
    return
  }
  try {
    const [routeList, groupStats] = await Promise.all([
      groupsAPI.listCompositeRoutes(selectedUnitId.value),
      groupsAPI.getStats(selectedUnitId.value)
    ])
    routes.value = routeList
    stats.value = { total_api_keys: groupStats.total_api_keys, active_api_keys: groupStats.active_api_keys }
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function loadAll() {
  loading.value = true
  try {
    const [vendorList, entryList] = await Promise.all([upstreamHubAPI.listVendors(), upstreamHubAPI.listAllEntries()])
    vendors.value = vendorList
    entries.value = entryList
    await loadUnits()
    await loadUnitData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    loading.value = false
  }
}

async function refreshEntries() {
  entries.value = await upstreamHubAPI.listAllEntries()
  vendors.value = await upstreamHubAPI.listVendors()
}

function openUnit(unitId: number) {
  selectedUnitId.value = unitId
  detailTab.value = 'supply'
  void router.replace({ query: { ...route.query, unit: String(unitId) } })
  void loadUnitData()
}

function backToList() {
  selectedUnitId.value = null
  const query = { ...route.query }
  delete query.unit
  void router.replace({ query })
  routes.value = []
  stats.value = null
}

function goToGroups() {
  void router.push('/admin/groups')
}

// ---------- 新建单元 ----------
const showCreateDialog = ref(false)
const createForm = reactive({ name: '', description: '' })

function openCreateUnit() {
  createForm.name = ''
  createForm.description = ''
  showCreateDialog.value = true
}

async function createUnit() {
  if (!createForm.name.trim()) return
  creating.value = true
  try {
    const group = await groupsAPI.create({
      name: createForm.name.trim(),
      description: createForm.description.trim() || undefined,
      platform: 'composite',
      rate_multiplier: 1,
      subscription_type: 'standard'
    })
    showCreateDialog.value = false
    appStore.showSuccess(t('admin.aggregation.unitCreated'))
    await loadAll()
    openUnit(group.id)
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    creating.value = false
  }
}

// ---------- 供给 ----------
async function toggleAttach(entry: UpstreamEntry) {
  if (!selectedUnitId.value) return
  const otherGroups = (entry.groups || [])
    .map((group) => group.group_id)
    .filter((groupId) => groupId !== selectedUnitId.value)
  const nextGroupIds = isAttached(entry) ? otherGroups : [...otherGroups, selectedUnitId.value]
  try {
    await upstreamHubAPI.setEntryGroups(entry.account_id, nextGroupIds)
    await refreshEntries()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function toggleEntryEnabled(entry: UpstreamEntry) {
  const entryGroupRef = groupRef(entry)
  if (!entryGroupRef) return
  try {
    await upstreamHubAPI.setEntryGroupEnabled(entry.account_id, entryGroupRef.group_id, !entryGroupRef.enabled)
    await refreshEntries()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

// ---------- 归一模型 ----------
const routeForm = reactive({
  public_model: '',
  target_platform: platformOptions.value[0]?.value || 'openai',
  upstream_model: '',
  priority: 100
})

async function addRoute() {
  if (!selectedUnitId.value) return
  if (!routeForm.public_model.trim()) {
    appStore.showError(t('admin.aggregation.routeRequired'))
    return
  }
  saving.value = true
  try {
    await groupsAPI.createCompositeRoute(selectedUnitId.value, {
      public_model: routeForm.public_model.trim(),
      match_type: 'exact',
      target_platform: routeForm.target_platform,
      upstream_model: routeForm.upstream_model.trim() || undefined,
      endpoint: 'any',
      priority: Number(routeForm.priority) || 100,
      enabled: true
    })
    routeForm.public_model = ''
    routeForm.upstream_model = ''
    appStore.showSuccess(t('admin.aggregation.routeCreated'))
    await loadUnitData()
    await loadUnits()
  } catch (err) {
    appStore.showError(errorMessage(err))
  } finally {
    saving.value = false
  }
}

async function toggleRouteEnabled(route: CompositeModelRoute) {
  if (!selectedUnitId.value) return
  try {
    await groupsAPI.updateCompositeRoute(selectedUnitId.value, route.id, {
      public_model: route.public_model,
      match_type: route.match_type,
      target_platform: route.target_platform,
      upstream_model: route.upstream_model,
      endpoint: route.endpoint,
      priority: route.priority,
      enabled: !route.enabled,
      notes: route.notes
    })
    await loadUnitData()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

async function removeRoute(route: CompositeModelRoute) {
  if (!selectedUnitId.value) return
  try {
    await groupsAPI.deleteCompositeRoute(selectedUnitId.value, route.id)
    appStore.showSuccess(t('admin.aggregation.routeDeleted'))
    await loadUnitData()
    await loadUnits()
  } catch (err) {
    appStore.showError(errorMessage(err))
  }
}

onMounted(async () => {
  await loadAll()
  const unitParam = Number(route.query.unit)
  if (Number.isFinite(unitParam) && unitParam > 0 && units.value.some((unit) => unit.id === unitParam)) {
    selectedUnitId.value = unitParam
    await loadUnitData()
  }
})
</script>
